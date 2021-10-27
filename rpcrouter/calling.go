package rpcrouter

import (
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	"time"
)

type CallOption struct {
	broadcast bool
	timeout   time.Duration
}

type CallOptionFunc func(opt *CallOption)

func WithBroadcast(b bool) CallOptionFunc {
	return func(opt *CallOption) {
		opt.broadcast = b
	}
}

func WithTimeout(timeout time.Duration) CallOptionFunc {
	return func(opt *CallOption) {
		opt.timeout = timeout
	}
}

func (self *Router) SingleCall(msg jsonrpc.IMessage, ns string, callOption *CallOption) (jsonrpc.IMessage, error) {
	if msg.IsRequest() {
		conn := self.Join()
		defer self.Leave(conn)

		msgvec := MsgVec{
			Msg:        msg,
			Namespace:  ns,
			FromConnId: conn.ConnId,
		}
		self.DeliverRequest(msgvec, callOption.timeout)
		resvec := <-conn.RecvChannel
		misc.AssertEqual(resvec.Msg.TraceId(), msg.TraceId(), "")
		return resvec.Msg, nil
	} else if msg.IsNotify() {
		self.DeliverNotify(MsgVec{
			Msg:       msg,
			Namespace: self.Name()})
		return nil, nil
	} else {
		return nil, ErrRequestNotifyRequired
	}
}

func (self *Router) GatherCall(msg jsonrpc.IMessage, ns string, limit int, callOption *CallOption) (resmsg jsonrpc.IMessage, err error) {
	if msg.IsRequest() {
		conn := self.Join()
		defer self.Leave(conn)

		reqmsg, _ := msg.(*jsonrpc.RequestMessage)

		servoIds := self.ListConns(reqmsg.Method, limit)

		var arr []interface{}

		for _, servoId := range servoIds {
			newId := misc.NewUuid()
			newmsg := reqmsg.Clone(newId)
			msgvec := MsgVec{
				Msg:        newmsg,
				Namespace:  ns,
				FromConnId: conn.ConnId,
				ToConnId:   servoId}

			self.DeliverRequest(msgvec, callOption.timeout)
		}
		log.Infof("send request %s to %d handlers", reqmsg.Method, len(servoIds))
		// wait for results
		for i := 0; i < len(servoIds); i++ {
			resMsgvec := <-conn.RecvChannel

			misc.Assert(resMsgvec.Msg.IsResultOrError(), "recved neither result not error")
			misc.AssertEqual(resMsgvec.Msg.TraceId(), reqmsg.TraceId(), "")
			arr = append(arr, jsonrpc.MessageInterface(resMsgvec.Msg))
		}
		resmsg := jsonrpc.NewResultMessage(reqmsg, arr)
		return resmsg, nil
	} else if msg.IsNotify() {
		conn := self.Join()
		defer self.Leave(conn)

		notifymsg, _ := msg.(*jsonrpc.NotifyMessage)
		servoIds := self.ListConns(notifymsg.Method, limit)

		for _, servoId := range servoIds {
			msgvec := MsgVec{
				Msg:        notifymsg,
				Namespace:  ns,
				FromConnId: conn.ConnId,
				ToConnId:   servoId}
			self.DeliverNotify(msgvec)
		}
		log.Infof("send notify %s to %d handlers", notifymsg.Method, len(servoIds))
		return nil, nil
	} else {
		return nil, ErrRequestNotifyRequired
	}
}

func (self *Router) CallOrNotify(msg jsonrpc.IMessage, ns string, opts ...CallOptionFunc) (jsonrpc.IMessage, error) {
	callOption := &CallOption{timeout: DefaultRequestTimeout}
	for _, optfunc := range opts {
		optfunc(callOption)
	}
	if callOption.broadcast {
		return self.GatherCall(msg, ns, 100, callOption)
	} else {
		return self.SingleCall(msg, ns, callOption)
	}
}
