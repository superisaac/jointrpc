package rpcrouter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	"time"
)

func (self *Router) deliverMessage(cmdMsg CmdMsg) *ConnT {
	msgvec := cmdMsg.MsgVec
	msg := cmdMsg.MsgVec.Msg
	msg.Log().WithFields(log.Fields{"from": msgvec.FromConnId}).Debugf("deliver message")
	if msg.IsRequest() {
		return self.deliverRequest(msgvec, cmdMsg.Timeout, cmdMsg.ChRes)
	} else if msg.IsNotify() {
		return self.deliverNotify(msgvec)
	} else if msg.IsResultOrError() {
		return self.deliverResultOrError(msgvec)
	}
	return nil
}

func (self *Router) deliverNotify(msgvec MsgVec) *ConnT {
	notifyMsg, ok := msgvec.Msg.(*jsonrpc.NotifyMessage)
	misc.Assert(ok, "bad msg type other than notify")
	notifyMsg.Log().Debugf("deliver notify")
	toConn, found := self.SelectConn(notifyMsg.Method, msgvec.ToConnId)
	if found {
		notifyMsg.Log().Debugf("selected conn %d", toConn.ConnId)
		if self.factory.Config.ValidateSchema() {
			if v, err := toConn.ValidateNotifyMsg(notifyMsg); !v && err != nil {
				notifyMsg.Log().Errorf("notify not validated, %s", err.Error())
				return nil
			}
		}

		return self.SendTo(
			toConn.ConnId, msgvec)
	}
	return nil
}

func (self *Router) deliverRequest(msgvec MsgVec, timeout time.Duration, chRes MsgChannel) *ConnT {
	reqMsg, ok := msgvec.Msg.(*jsonrpc.RequestMessage)
	misc.Assert(ok, "bad msg type other than request")

	reqMsg.Log().Debugf("deliver request")
	msgId := reqMsg.Id
	fromConnId := msgvec.FromConnId
	toConn, found := self.SelectConn(reqMsg.Method, msgvec.ToConnId)
	if found {
		reqMsg.Log().Debugf("selected conn %d", toConn.ConnId)
		if self.factory.Config.ValidateSchema() {
			if v, errmsg := toConn.ValidateRequestMsg(reqMsg); !v && errmsg != nil {

				errVec := MsgVec{
					Msg:        errmsg,
					FromConnId: toConn.ConnId,
				}
				return self.SendTo(fromConnId, errVec)
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

		// go func() {

		// 	time.Sleep(timeout)
		// 	time.Sleep(1 * time.Second)
		// 	self.TryClearPendingRequest(msgId)
		// }()
		targetVec := msgvec
		targetVec.Msg = reqMsg
		return self.SendTo(toConn.ConnId, targetVec)
	} else {
		errMsg := jsonrpc.ErrMethodNotFound.WithData(fmt.Sprintf("request method %s", reqMsg.Method)).ToMessage(reqMsg)
		errMsg.SetTraceId(reqMsg.TraceId())
		errMsgVec := MsgVec{Msg: errMsg}
		return self.SendTo(fromConnId, errMsgVec)
	}
}

func (self *Router) deliverResultOrError(msgvec MsgVec) *ConnT {
	msg := msgvec.Msg
	//if msgId, ok := msg.MustId().(string); ok {
	msg.Log().Debugf("deliver result or error")
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
						return self.SendTo(reqt.FromConnId, errVec)
					}
				}
			}

			newRes := jsonrpc.NewResultMessage(origReq, resMsg.Result)
			newVec := msgvec
			newVec.Msg = newRes
			if reqt.ChRes != nil {
				reqt.ChRes <- newVec
				return nil
			} else {
				return self.SendTo(reqt.FromConnId, newVec)
			}
		} else if errMsg, ok := msg.(*jsonrpc.ErrorMessage); ok {
			newErr := jsonrpc.NewErrorMessage(origReq, errMsg.Error)
			newVec := msgvec
			newVec.Msg = newErr
			if reqt.ChRes != nil {
				reqt.ChRes <- newVec
				return nil
			} else {
				return self.SendTo(reqt.FromConnId, newVec)
			}
		} else {
			msg.Log().Fatalf("msg is neither result nor error")
		}
	} else {
		msg.Log().Warn("fail to find request from this result/error")
	}
	return nil
}
