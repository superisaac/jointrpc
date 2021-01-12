package rpcrouter

import (
	"fmt"
	uuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
)

func (self *Router) SingleCall(msgvec MsgVec) (jsonrpc.IMessage, string, error) {
	msg := msgvec.Msg
	if msg.IsRequest() {
		conn := self.Join()
		defer self.Leave(conn)

		self.ChMsg <- CmdMsg{
			MsgVec: MsgVec{
				Msg:        msg,
				TraceId:    msgvec.TraceId,
				FromConnId: conn.ConnId}}
		msgvec := <-conn.RecvChannel
		return msgvec.Msg, msgvec.TraceId, nil
	} else if msg.IsNotify() {
		self.ChMsg <- CmdMsg{
			MsgVec: MsgVec{
				Msg:        msg,
				TraceId:    msgvec.TraceId,
				FromConnId: ZeroCID},
		}
		return nil, "", nil
	} else {
		return nil, "", ErrRequestNotifyRequired
	}
}

func (self *Router) GatherCall(msgvec MsgVec, limit int) (resmsg jsonrpc.IMessage, traceId string, err error) {
	msg := msgvec.Msg
	if msg.IsRequest() {
		conn := self.Join()
		defer self.Leave(conn)

		reqmsg, _ := msg.(*jsonrpc.RequestMessage)

		servoIds := self.ListConns(reqmsg.Method, limit)

		var arr []interface{}

		for _, servoId := range servoIds {
			msgId := uuid.New().String()
			newmsg := reqmsg.Clone(msgId)
			self.ChMsg <- CmdMsg{
				MsgVec: MsgVec{
					Msg:          newmsg,
					TraceId:      msgvec.TraceId,
					FromConnId:   conn.ConnId,
					TargetConnId: servoId},
			}
		}
		log.Infof("send request %s to %d handlers", reqmsg.Method, len(servoIds))
		// wait for results
		var resTraceId string
		for i := 0; i < len(servoIds); i++ {
			resMsgvec := <-conn.RecvChannel

			misc.Assert(resMsgvec.Msg.IsResultOrError(), "recved neither result not error")
			resTraceId = resMsgvec.TraceId
			misc.Assert(resTraceId == msgvec.TraceId,
				fmt.Sprintf("traceid mismatch %s %s", msgvec.TraceId, resTraceId))
			arr = append(arr, resMsgvec.Msg.Interface())
		}
		resmsg := jsonrpc.NewResultMessage(reqmsg.Id, arr, nil)
		return resmsg, resTraceId, nil
	} else if msg.IsNotify() {
		conn := self.Join()
		defer self.Leave(conn)

		notifymsg, _ := msg.(*jsonrpc.NotifyMessage)
		servoIds := self.ListConns(notifymsg.Method, limit)

		for _, servoId := range servoIds {
			self.ChMsg <- CmdMsg{
				MsgVec: MsgVec{
					Msg:          notifymsg,
					TraceId:      msgvec.TraceId,
					FromConnId:   conn.ConnId,
					TargetConnId: servoId},
			}
		}
		log.Infof("send notify %s to %d handlers", notifymsg.Method, len(servoIds))
		return nil, "", nil
	} else {
		return nil, "", ErrRequestNotifyRequired
	}
}

func (self *Router) CallOrNotify(msgvec MsgVec, broadcast bool) (jsonrpc.IMessage, string, error) {
	if broadcast {
		return self.GatherCall(msgvec, 100)
	} else {
		return self.SingleCall(msgvec)
	}
}
