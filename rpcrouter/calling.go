package rpcrouter

import (
	//"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
)

func (self *Router) SingleCall(msg jsonrpc.IMessage) (jsonrpc.IMessage, error) {
	if msg.IsRequest() {
		conn := self.Join()
		defer self.Leave(conn)

		self.ChMsg <- CmdMsg{
			MsgVec: MsgVec{Msg: msg,
				FromConnId: conn.ConnId}}
		msgvec := <-conn.RecvChannel
		return msgvec.Msg, nil
	} else if msg.IsNotify() {
		self.ChMsg <- CmdMsg{
			MsgVec: MsgVec{Msg: msg, FromConnId: 0},
		}
		return nil, nil
	} else {
		return nil, ErrRequestNotifyRequired
	}
}

func (self *Router) GatherCall(msg jsonrpc.IMessage, limit int) (jsonrpc.IMessage, error) {
	if msg.IsRequest() {
		conn := self.Join()
		defer self.Leave(conn)

		reqmsg, _ := msg.(*jsonrpc.RequestMessage)

		servos := self.ListConns(reqmsg.Method, limit)

		var arr []interface{}

		for _, servoConn := range servos {
			msgId := uuid.New().String()
			newmsg := reqmsg.Clone(msgId)
			self.ChMsg <- CmdMsg{
				MsgVec: MsgVec{
					Msg:          newmsg,
					FromConnId:   conn.ConnId,
					TargetConnId: servoConn.ConnId},
			}
		}
		log.Infof("send request %s to %d handlers", reqmsg.Method,
			len(servos))
		// wait for results
		for i := 0; i < len(servos); i++ {
			msgvec := <-conn.RecvChannel
			misc.Assert(msgvec.Msg.IsResultOrError(), "recved neither result not error")
			arr = append(arr, msgvec.Msg.Interface())
		}
		resmsg := jsonrpc.NewResultMessage(reqmsg.Id, arr, nil)
		return resmsg, nil
	} else if msg.IsNotify() {
		conn := self.Join()
		defer self.Leave(conn)

		notifymsg, _ := msg.(*jsonrpc.NotifyMessage)
		servos := self.ListConns(notifymsg.Method, limit)

		for _, servoConn := range servos {
			self.ChMsg <- CmdMsg{
				MsgVec: MsgVec{
					Msg:          notifymsg,
					FromConnId:   conn.ConnId,
					TargetConnId: servoConn.ConnId},
			}
		}
		log.Infof("send notify %s to %d handlers", notifymsg.Method,
			len(servos))
		return nil, nil
	} else {
		return nil, ErrRequestNotifyRequired
	}
}

func (self *Router) CallOrNotify(msg jsonrpc.IMessage, broadcast bool) (jsonrpc.IMessage, error) {
	if broadcast {
		return self.GatherCall(msg, 100)
	} else {
		return self.SingleCall(msg)
	}
}
