package rpcrouter

import (
	//"fmt"
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/superisaac/jointrpc/misc"
	jsonrpc "github.com/superisaac/jsonrpc"
	schema "github.com/superisaac/jsonrpc/schema"
	"time"
)

func NewConn() *ConnT {
	connId := NextCID()
	methods := make(map[string]MethodInfo)
	pendings := make(map[interface{}]ConnPending)
	conn := &ConnT{ConnId: connId,
		ServeMethods: methods,
		msgOutput:    make(MsgChannel, 5000),
		msgInput:     make(MsgChannel, 5000),
		pendings:     pendings,
	}
	return conn
}

func (self *ConnT) Destruct() {
	self.Log().Debugf("conn destruct")
	self.Namespace = ""
	self.router = nil
	self.msgInput = nil
	self.msgOutput = nil
	for _, pending := range self.pendings {
		self.returnTimeout(pending)
	}
	self.pendings = nil
}

func (self ConnT) MsgOutput() MsgChannel {
	return self.msgOutput
}

func (self ConnT) MsgInput() MsgChannel {
	return self.msgInput
}

func (self ConnT) Joined() bool {
	return self.Namespace != ""
}

func (self ConnT) GetMethods() []string {
	var keys []string
	for k := range self.ServeMethods {
		keys = append(keys, k)
	}
	return keys
}

func (self ConnT) ValidateRequestMsg(reqMsg *jsonrpc.RequestMessage) (bool, jsonrpc.IMessage) {
	if info, ok := self.ServeMethods[reqMsg.Method]; ok && info.Schema() != nil {
		s := info.Schema()
		validator := schema.NewSchemaValidator()
		errPos := validator.Validate(s, jsonrpc.MessageInterface(reqMsg))
		if errPos != nil {
			errmsg := errPos.ToMessage(reqMsg)
			return false, errmsg
		}
	}
	return true, nil
}

func (self ConnT) ValidateNotifyMsg(notifyMsg *jsonrpc.NotifyMessage) (bool, error) {
	if info, ok := self.ServeMethods[notifyMsg.Method]; ok && info.Schema() != nil {
		s := info.Schema()
		validator := schema.NewSchemaValidator()
		errPos := validator.Validate(s, jsonrpc.MessageInterface(notifyMsg))
		if errPos != nil {
			return false, errPos
		}
	}
	return true, nil
}

func (self ConnT) ValidateResultMsg(resMsg *jsonrpc.ResultMessage, reqMsg *jsonrpc.RequestMessage) (bool, jsonrpc.IMessage) {
	if info, ok := self.ServeMethods[reqMsg.Method]; ok && info.Schema() != nil {
		s := info.Schema()
		validator := schema.NewSchemaValidator()
		errPos := validator.Validate(s, jsonrpc.MessageInterface(resMsg))
		if errPos != nil {
			errmsg := errPos.ToMessage(reqMsg)
			return false, errmsg
		}
	}
	return true, nil
}

func (self *ConnT) Touch() {
	self.lastPing = time.Now()
}
func (self *ConnT) SetWatchState(w bool) {
	self.watchState = w
}

func (self *ConnT) StateChannel() chan *ServerState {
	if self.stateChannel == nil {
		self.stateChannel = make(chan *ServerState, misc.DefaultChanSize())
	}
	return self.stateChannel
}

func (self ConnT) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"namespace":  self.Namespace,
		"connid":     self.ConnId,
		"remoteaddr": self.PeerAddr,
	})
}

func (self *ConnT) HandleRouteMessage(ctx context.Context, cmdMsg CmdMsg) error {
	msg := cmdMsg.Msg
	if msg.IsRequest() {
		// Down message
		return self.handleRequest(ctx, cmdMsg)
	} else if msg.IsNotify() {
		// Down Message
		return self.handleNotify(ctx, cmdMsg)
	} else {
		// Up message
		misc.Assert(msg.IsResultOrError(), "must be result or error")
		return self.handleResultOrError(ctx, cmdMsg)
	}
}

func (self *ConnT) handleRequest(ctx context.Context, cmdMsg CmdMsg) error {
	reqMsg, _ := cmdMsg.Msg.(*jsonrpc.RequestMessage)
	if self.router.factory.Config.ValidateSchema() {
		if v, errmsg := self.ValidateRequestMsg(reqMsg); !v && errmsg != nil {
			cmdMsg.ChRes <- cmdMsg.Res(errmsg)
			return nil
		}
	}

	if cmdMsg.Timeout <= 0 {
		cmdMsg.Timeout = DefaultRequestTimeout
	}

	expireTime := time.Now().Add(cmdMsg.Timeout)
	newMsgId := misc.NewUuid()
	newReqMsg := reqMsg.Clone(newMsgId)
	self.pendings[newMsgId] = ConnPending{
		cmdMsg: cmdMsg,
		Expire: expireTime,
	}

	newCmdMsg := CmdMsg{Msg: newReqMsg, Namespace: cmdMsg.Namespace}
	self.MsgOutput() <- newCmdMsg
	return nil
}

func (self *ConnT) handleNotify(ctx context.Context, cmdMsg CmdMsg) error {
	notifyMsg, _ := cmdMsg.Msg.(*jsonrpc.NotifyMessage)
	if self.router.factory.Config.ValidateSchema() {
		if v, err := self.ValidateNotifyMsg(notifyMsg); !v && err != nil {
			notifyMsg.Log().Errorf("notify not valid, %s", err.Error())
			return nil
		}
	}

	self.MsgOutput() <- cmdMsg
	return nil
}

func (self *ConnT) handleResultOrError(ctx context.Context, cmdMsg CmdMsg) error {
	msg := cmdMsg.Msg
	msgId := msg.MustId()
	if pending, ok := self.pendings[msgId]; ok {
		origReq, ok := pending.cmdMsg.Msg.(*jsonrpc.RequestMessage)
		misc.Assert(ok, "original is not request")
		// delete pendings
		delete(self.pendings, msgId)
		// check the expiration
		if time.Now().After(pending.Expire) {
			origReq.Log().Infof("request expired")
			return nil
		}

		if msg.TraceId() != origReq.TraceId() {
			msg.Log().Warnf("result traceid is different from original request %s", origReq.TraceId())
		}

		if resMsg, ok := msg.(*jsonrpc.ResultMessage); ok {
			if self.router.factory.Config.ValidateSchema() {
				if v, errmsg := self.ValidateResultMsg(resMsg, origReq); !v && errmsg != nil {
					pending.cmdMsg.ChRes <- cmdMsg.Res(errmsg)
					return nil
				}
			}
			newRes := jsonrpc.NewResultMessage(origReq, resMsg.Result)
			pending.cmdMsg.ChRes <- cmdMsg.Res(newRes)
			return nil
		} else if errMsg, ok := msg.(*jsonrpc.ErrorMessage); ok {
			newErr := jsonrpc.NewErrorMessage(origReq, errMsg.Error)
			pending.cmdMsg.ChRes <- cmdMsg.Res(newErr)
			return nil
		}
	} else {
		cmdMsg.Msg.Log().Warnf("cannot find pending request")
	}
	return nil
}

func (self *ConnT) returnTimeout(pending ConnPending) {
	defer func() {
		if r := recover(); r != nil {
			pending.cmdMsg.Msg.Log().Warnf("recovered send on timeout %+v", r)
			if err, ok := r.(error); ok && err.Error() == "send on closed channel" {
				pending.cmdMsg.Msg.Log().Warnf("channel already closed %+v", pending.cmdMsg)
			}
		}
	}()

	errMsg := jsonrpc.ErrTimeout.ToMessage(pending.cmdMsg.Msg)
	errCmdMsg := pending.cmdMsg
	errCmdMsg.Msg = errMsg
	pending.cmdMsg.ChRes <- errCmdMsg
}

func (self *ConnT) ClearPendings() {
	now := time.Now()
	newPendings := make(map[interface{}]ConnPending)

	for reqMsgId, pending := range self.pendings {
		if now.After(pending.Expire) {
			self.returnTimeout(pending)
		} else {
			newPendings[reqMsgId] = pending
		}
	}
	self.pendings = newPendings
}
