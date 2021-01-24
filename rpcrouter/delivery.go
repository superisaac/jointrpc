package rpcrouter

import (
	"time"

	uuid "github.com/google/uuid"	
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"		
)

func (self *Router) DeliverMessage(cmdMsg CmdMsg) *ConnT {
	msg := cmdMsg.MsgVec.Msg
	msg.Log().WithFields(log.Fields{"from": cmdMsg.MsgVec.FromConnId}).Debugf("Deliver message")
	if msg.IsRequest() {
		return self.deliverRequest(cmdMsg)
	} else if msg.IsNotify() {
		return self.deliverNotify(cmdMsg)
	} else if msg.IsResultOrError() {
		return self.deliverResultOrError(cmdMsg)
	}
	return nil
}

func (self *Router) deliverNotify(cmdMsg CmdMsg) *ConnT {
	msg := cmdMsg.MsgVec.Msg
	toConn, found := self.SelectConn(msg.MustMethod(), cmdMsg.MsgVec.ToConnId)
	if found {
		return self.SendTo(
			toConn.ConnId, cmdMsg.MsgVec)
	}
	return nil
}

func (self *Router) deliverRequest(cmdMsg CmdMsg) *ConnT {
	msg := cmdMsg.MsgVec.Msg
	fromConnId := cmdMsg.MsgVec.FromConnId
	toConn, found := self.SelectConn(msg.MustMethod(), cmdMsg.MsgVec.ToConnId)
	if found {
		timeout := cmdMsg.Timeout
		if timeout <= 0 {
			timeout = DefaultRequestTimeout
		}

		//fmt.Printf("timeout %d\n", timeout)
		expireTime := time.Now().Add(timeout)
		reqMsg, ok := msg.(*jsonrpc.RequestMessage)
		misc.Assert(ok, "bad msg type other than request")

		msgId := reqMsg.Id

		func() {
			// update pending Request
			self.lock("deliverRequest")
			defer self.unlock("deliverRequest")

			if _, ok := self.pendingRequests[msgId]; ok {
				msgId = uuid.New().String()
				reqMsg.Log().Infof("msg id already exist, change a new one %s", msgId)
				reqMsg = reqMsg.Clone(msgId)
			}
			self.pendingRequests[msgId] = PendingT{
				ReqMsg:     reqMsg,
				FromConnId: cmdMsg.MsgVec.FromConnId,
				ToConnId:   toConn.ConnId,			
				Expire:     expireTime,
			}
		}()
		go func() {
			time.Sleep(timeout)
			time.Sleep(1 * time.Second)
			//time.Sleep(int64(timeout.Seconds() + 1) * time.Second)
			self.TryClearPendingRequest(msgId)
		}()
		targetVec := cmdMsg.MsgVec
		targetVec.Msg = reqMsg
		return self.SendTo(toConn.ConnId, targetVec)
	} else {
		errMsg := jsonrpc.RPCErrorMessage(msg, 404, "method not found", false)
		errMsg.SetTraceId(msg.TraceId())
		errMsgVec := MsgVec{Msg: errMsg}
		return self.SendTo(fromConnId, errMsgVec)
	}
}

func (self *Router) deliverResultOrError(cmdMsg CmdMsg) *ConnT {
	msg := cmdMsg.MsgVec.Msg
	//if msgId, ok := msg.MustId().(string); ok {
	msgId := msg.MustId()
	if reqt, ok := self.pendingRequests[msgId]; ok {
		self.DeletePending(msgId)
		
		if cmdMsg.MsgVec.FromConnId != reqt.ToConnId {
			msg.Log().Warnf("msg conn %d not from the delivered conn %d", cmdMsg.MsgVec.FromConnId, reqt.ToConnId)
		}
		origReq := reqt.ReqMsg
		if msg.TraceId() != origReq.TraceId() {
			msg.Log().Warnf("result trace is different from request %s", origReq.TraceId())
		}
		if resMsg, ok := msg.(*jsonrpc.ResultMessage); ok {
			newRes := jsonrpc.NewResultMessage(origReq, resMsg.Result, nil)
			newVec := cmdMsg.MsgVec
			newVec.Msg = newRes
			return self.SendTo(reqt.FromConnId, newVec)
		} else if errMsg, ok := msg.(*jsonrpc.ErrorMessage); ok {
			newErr := jsonrpc.NewErrorMessage(origReq, errMsg.Error, nil)
			newVec := cmdMsg.MsgVec
			newVec.Msg = newErr
			return self.SendTo(reqt.FromConnId, newVec)
		} else {
			msg.Log().Fatalf("msg is neither result nor error")
		}
	} else {
		msg.Log().Warn("fail to find request from this result/error")
	}
	return nil
}
