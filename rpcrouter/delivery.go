package rpcrouter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	"time"
)

func (self *Router) relayMessage(cmdMsg CmdMsg) {
	msg := cmdMsg.MsgVec.Msg
	misc.Assert(msg.IsRequestOrNotify(), "router only support request and notify")
	toConn, found := self.SelectConn(msg.MustMethod(), cmdMsg.MsgVec.ToConnId)
	if found {
		toConn.ChRouteMsg <- cmdMsg
	} else if msg.IsRequest() {
		reqMsg, _ := msg.(*jsonrpc.RequestMessage)
		errMsg := jsonrpc.ErrMethodNotFound.WithData(fmt.Sprintf("request method %s", reqMsg.Method)).ToMessage(reqMsg)
		errMsg.SetTraceId(reqMsg.TraceId())
		errMsgVec := MsgVec{Msg: errMsg}
		cmdMsg.ChRes <- errMsgVec
	}
}

func (self *Router) deliverMessage(cmdMsg CmdMsg) {
	msgvec := cmdMsg.MsgVec
	msg := cmdMsg.MsgVec.Msg
	msg.Log().WithFields(log.Fields{"from": msgvec.FromConnId}).Debugf("deliver message")
	if msg.IsRequest() {
		self.deliverRequest(cmdMsg) //msgvec, cmdMsg.Timeout, cmdMsg.ChRes)
	} else if msg.IsNotify() {
		self.deliverNotify(cmdMsg)
	} else if msg.IsResultOrError() {
		self.deliverResultOrError(cmdMsg)
	}
}

func (self *Router) deliverNotify(cmdMsg CmdMsg) {
	msgvec := cmdMsg.MsgVec
	notifyMsg, ok := msgvec.Msg.(*jsonrpc.NotifyMessage)
	misc.Assert(ok, "bad msg type other than notify")
	notifyMsg.Log().Debugf("deliver notify")
	toConn, found := self.SelectConn(notifyMsg.Method, msgvec.ToConnId)
	if found {
		notifyMsg.Log().Debugf("selected conn %d", toConn.ConnId)
		if self.factory.Config.ValidateSchema() {
			if v, err := toConn.ValidateNotifyMsg(notifyMsg); !v && err != nil {
				notifyMsg.Log().Errorf("notify not validated, %s", err.Error())
				return
			}
		}
		self.SendTo(toConn.ConnId, msgvec)
	}
}

func (self *Router) deliverRequest(cmdMsg CmdMsg) {
	msgvec := cmdMsg.MsgVec
	timeout := cmdMsg.Timeout
	chRes := cmdMsg.ChRes
	reqMsg, ok := msgvec.Msg.(*jsonrpc.RequestMessage)
	misc.Assert(ok, "bad msg type other than request")

	misc.Assert(chRes != nil, "chRes is nil")

	reqMsg.Log().Debugf("deliver request")
	msgId := reqMsg.Id
	toConn, found := self.SelectConn(reqMsg.Method, msgvec.ToConnId)
	if found {
		reqMsg.Log().Debugf("selected conn %d", toConn.ConnId)
		if self.factory.Config.ValidateSchema() {
			if v, errmsg := toConn.ValidateRequestMsg(reqMsg); !v && errmsg != nil {
				errVec := MsgVec{
					Msg:        errmsg,
					FromConnId: toConn.ConnId,
				}

				chRes <- errVec
				return
			}
		}
		if timeout <= 0 {
			timeout = DefaultRequestTimeout
		}
		//fmt.Printf("timeout %d\n", timeout)
		expireTime := time.Now().Add(timeout)
		origMsgId := msgId
		msgId = misc.NewUuid()
		reqMsg = reqMsg.Clone(msgId)
		self.addPending(msgId, PendingT{
			ReqMsg:     reqMsg,
			OrigMsgId:  origMsgId,
			FromConnId: msgvec.FromConnId,
			ToConnId:   toConn.ConnId,
			Expire:     expireTime,
			ChRes:      chRes,
		})

		targetVec := msgvec
		targetVec.Msg = reqMsg
		self.SendTo(toConn.ConnId, targetVec)
	} else {
		errMsg := jsonrpc.ErrMethodNotFound.WithData(fmt.Sprintf("request method %s", reqMsg.Method)).ToMessage(reqMsg)
		errMsg.SetTraceId(reqMsg.TraceId())
		errMsgVec := MsgVec{Msg: errMsg}
		chRes <- errMsgVec
	}
}

func (self *Router) deliverResultOrError(cmdMsg CmdMsg) {
	msgvec := cmdMsg.MsgVec
	msg := msgvec.Msg
	//if msgId, ok := msg.MustId().(string); ok {
	msg.Log().Infof("deliver result or error")
	msgId := msg.MustId()
	if reqt, ok := self.getAndDeletePendings(msgId); ok {
		if msgvec.FromConnId != reqt.ToConnId {
			msg.Log().Warnf("msg conn %d not from the delivered conn %d", msgvec.FromConnId, reqt.ToConnId)
		}
		origReq := reqt.ReqMsg.Clone(reqt.OrigMsgId)
		if msg.TraceId() != origReq.TraceId() {
			msg.Log().Warnf("result trace is different from request %s", origReq.TraceId())
		}
		if resMsg, ok := msg.(*jsonrpc.ResultMessage); ok {
			if self.factory.Config.ValidateSchema() {
				// validate result message
				if vConn, ok := self.GetConn(reqt.ToConnId); ok {
					if v, errmsg := vConn.ValidateResultMsg(resMsg, origReq); !v && errmsg != nil {
						errVec := MsgVec{
							Msg:        errmsg,
							FromConnId: msgvec.FromConnId,
						}
						reqt.ChRes <- errVec
						return
					}
				}
			}

			newRes := jsonrpc.NewResultMessage(origReq, resMsg.Result)
			newVec := msgvec
			newVec.Msg = newRes
			reqt.ChRes <- newVec
			return
		} else if errMsg, ok := msg.(*jsonrpc.ErrorMessage); ok {
			newErr := jsonrpc.NewErrorMessage(origReq, errMsg.Error)
			newVec := msgvec
			newVec.Msg = newErr
			reqt.ChRes <- newVec
			return
		} else {
			msg.Log().Fatalf("msg is neither result nor error")
		}
	} else {
		msg.Log().Warn("fail to find request from this result/error")
	}
}
