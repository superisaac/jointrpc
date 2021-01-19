package rpcrouter

import (
	uuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
)

func (self *Router) SingleCall(msgvec MsgVec) (jsonrpc.IMessage, error) {
	msg := msgvec.Msg
	if msg.IsRequest() {
		conn := self.Join(false)
		defer self.Leave(conn)

		self.ChMsg <- CmdMsg{
			MsgVec: MsgVec{
				Msg:        msg,
				FromConnId: conn.ConnId}}
		resvec := <-conn.RecvChannel
		misc.AssertEqual(resvec.Msg.TraceId(), msg.TraceId(), "")
		return resvec.Msg, nil
	} else if msg.IsNotify() {
		self.ChMsg <- CmdMsg{
			MsgVec: MsgVec{
				Msg:        msg,
				FromConnId: ZeroCID},
		}
		return nil, nil
	} else {
		return nil, ErrRequestNotifyRequired
	}
}

func (self *Router) GatherCall(msgvec MsgVec, limit int) (resmsg jsonrpc.IMessage, err error) {
	msg := msgvec.Msg
	if msg.IsRequest() {
		conn := self.Join(false)
		defer self.Leave(conn)

		reqmsg, _ := msg.(*jsonrpc.RequestMessage)

		servoIds := self.ListConns(reqmsg.Method, limit)

		var arr []interface{}

		for _, servoId := range servoIds {
			newId := uuid.New().String()
			newmsg := reqmsg.Clone(newId)
			self.ChMsg <- CmdMsg{
				MsgVec: MsgVec{
					Msg:          newmsg,
					FromConnId:   conn.ConnId,
					TargetConnId: servoId},
			}
		}
		log.Infof("send request %s to %d handlers", reqmsg.Method, len(servoIds))
		// wait for results
		for i := 0; i < len(servoIds); i++ {
			resMsgvec := <-conn.RecvChannel

			misc.Assert(resMsgvec.Msg.IsResultOrError(), "recved neither result not error")
			misc.AssertEqual(resMsgvec.Msg.TraceId(), reqmsg.TraceId(), "")
			arr = append(arr, resMsgvec.Msg.Interface())
		}
		resmsg := jsonrpc.NewResultMessage(reqmsg, arr, nil)
		return resmsg, nil
	} else if msg.IsNotify() {
		conn := self.Join(false)
		defer self.Leave(conn)

		notifymsg, _ := msg.(*jsonrpc.NotifyMessage)
		servoIds := self.ListConns(notifymsg.Method, limit)

		for _, servoId := range servoIds {
			self.ChMsg <- CmdMsg{
				MsgVec: MsgVec{
					Msg:          notifymsg,
					FromConnId:   conn.ConnId,
					TargetConnId: servoId},
			}
		}
		log.Infof("send notify %s to %d handlers", notifymsg.Method, len(servoIds))
		return nil, nil
	} else {
		return nil, ErrRequestNotifyRequired
	}
}

func (self *Router) CallOrNotify(msgvec MsgVec, broadcast bool) (jsonrpc.IMessage, error) {
	if broadcast {
		return self.GatherCall(msgvec, 100)
	} else {
		return self.SingleCall(msgvec)
	}
}
