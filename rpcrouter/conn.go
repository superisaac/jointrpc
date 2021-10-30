package rpcrouter

import (
	//"fmt"
	"context"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	"github.com/superisaac/jointrpc/misc"
	"time"
)

func NewConn() *ConnT {
	connId := NextCID()
	ch := make(MsgChannel, 10000)
	methods := make(map[string]MethodInfo)
	pendings := make(map[interface{}]ConnPending)
	conn := &ConnT{ConnId: connId,
		RecvChannel:  ch,
		ServeMethods: methods,
		ChRouteMsg:   make(chan CmdMsg, 10000),
		pendings:     pendings,
	}
	return conn
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
	msg := cmdMsg.MsgVec.Msg
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
	reqMsg, _ := cmdMsg.MsgVec.Msg.(*jsonrpc.RequestMessage)
	if self.router.factory.Config.ValidateSchema() {
		if v, errmsg := self.ValidateRequestMsg(reqMsg); !v && errmsg != nil {
			errVec := MsgVec{
				Msg:        errmsg,
				FromConnId: self.ConnId,
			}
			cmdMsg.ChRes <- errVec
			return nil
		}
	}

	if cmdMsg.Timeout <= 0 {
		cmdMsg.Timeout = DefaultRequestTimeout
	}

	expireTime := time.Now().Add(cmdMsg.Timeout)
	newMsgId := misc.NewUuid()
	reqMsg = reqMsg.Clone(newMsgId)
	self.pendings[newMsgId] = ConnPending{
		cmdMsg: cmdMsg,
		Expire: expireTime,
	}
	self.RecvChannel <- MsgVec{Msg: reqMsg}
	return nil
}

func (self *ConnT) handleNotify(ctx context.Context, cmdMsg CmdMsg) error {
	notifyMsg, _ := cmdMsg.MsgVec.Msg.(*jsonrpc.NotifyMessage)
	if self.router.factory.Config.ValidateSchema() {
		if v, err := self.ValidateNotifyMsg(notifyMsg); !v && err != nil {
			notifyMsg.Log().Errorf("notify not valid, %s", err.Error())
			return nil
		}
	}
	self.RecvChannel <- MsgVec{Msg: notifyMsg}
	return nil
}

func (self *ConnT) handleResultOrError(ctx context.Context, cmdMsg CmdMsg) error {
	msg := cmdMsg.MsgVec.Msg
	msgId := msg.MustId()
	if pending, ok := self.pendings[msgId]; ok {
		origReq, ok := pending.cmdMsg.MsgVec.Msg.(*jsonrpc.RequestMessage)
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
					errVec := MsgVec{
						Msg:        errmsg,
						FromConnId: self.ConnId,
					}
					pending.cmdMsg.ChRes <- errVec
					return nil
				}
			}
			newRes := jsonrpc.NewResultMessage(origReq, resMsg.Result)
			newVec := cmdMsg.MsgVec
			newVec.Msg = newRes
			pending.cmdMsg.ChRes <- newVec
			return nil
		} else if errMsg, ok := msg.(*jsonrpc.ErrorMessage); ok {
			newErr := jsonrpc.NewErrorMessage(origReq, errMsg.Error)
			newVec := cmdMsg.MsgVec
			newVec.Msg = newErr
			pending.cmdMsg.ChRes <- newVec
			return nil
		}
	} else {
		cmdMsg.MsgVec.Msg.Log().Warnf("cannot find pending request")
	}
	return nil
}
