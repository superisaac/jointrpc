package rpcrouter

import (
	//"fmt"
	log "github.com/sirupsen/logrus"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jsonz"
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

func (self *Router) SingleCall(msg jsonz.Message, ns string, callOption *CallOption) (jsonz.Message, error) {
	if msg.IsRequest() {
		chRes := make(MsgChannel, 10)

		self.PostMessage(CmdMsg{
			Msg:       msg,
			Namespace: ns,
			Timeout:   callOption.timeout,
			ChRes:     chRes, //conn.MsgOutput(),
		})
		resvec := <-chRes
		misc.AssertEqual(resvec.Msg.TraceId(), msg.TraceId(), "")
		return resvec.Msg, nil
	} else if msg.IsNotify() {
		self.PostMessage(CmdMsg{
			Msg:       msg,
			Namespace: self.Name(),
		})
		return nil, nil
	} else {
		return nil, ErrRequestNotifyRequired
	}
}

func (self *Router) GatherCall(msg jsonz.Message, ns string, limit int, callOption *CallOption) (resmsg jsonz.Message, err error) {
	if msg.IsRequest() {
		reqmsg, _ := msg.(*jsonz.RequestMessage)

		servoIds := self.ListConns(reqmsg.Method, limit)

		var arr []interface{}

		chRes := make(MsgChannel, len(servoIds))

		for _, servoId := range servoIds {
			newId := jsonz.NewUuid()
			newmsg := reqmsg.Clone(newId)
			self.PostMessage(CmdMsg{
				Msg:       newmsg,
				Namespace: ns,
				Timeout:   callOption.timeout,
				ChRes:     chRes,
				ConnId:    servoId,
			})
		}
		log.Infof("send request %s to %d handlers", reqmsg.Method, len(servoIds))
		// wait for results
		for i := 0; i < len(servoIds); i++ {
			resMsgvec := <-chRes

			misc.Assert(resMsgvec.Msg.IsResultOrError(), "recved neither result not error")
			misc.AssertEqual(resMsgvec.Msg.TraceId(), reqmsg.TraceId(), "")
			arr = append(arr, jsonz.MessageInterface(resMsgvec.Msg))
		}
		resmsg := jsonz.NewResultMessage(reqmsg, arr)
		return resmsg, nil
	} else if msg.IsNotify() {
		notifymsg, _ := msg.(*jsonz.NotifyMessage)
		servoIds := self.ListConns(notifymsg.Method, limit)

		for _, servoId := range servoIds {
			self.PostMessage(CmdMsg{
				Msg:       notifymsg,
				Namespace: ns,
				ConnId:    servoId,
			})
		}
		log.Infof("send notify %s to %d handlers", notifymsg.Method, len(servoIds))
		return nil, nil
	} else {
		return nil, ErrRequestNotifyRequired
	}
}

func (self *Router) CallOrNotify(msg jsonz.Message, ns string, opts ...CallOptionFunc) (jsonz.Message, error) {
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
