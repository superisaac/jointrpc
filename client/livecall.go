package client

import (
	"github.com/pkg/errors"
	"time"
	//"fmt"
	"context"
	//"github.com/superisaac/jointrpc/msgutil"
	//intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/misc"

	log "github.com/sirupsen/logrus"
)

func (self *RPCClient) LiveCall(rootCtx context.Context, reqmsg jsonrpc.IMessage, callback LiveCallback, opts ...CallOptionFunc) error {
	if !self.connected {
		log.Warnf("live stream is not connected")
		return errors.New("live stream is not connected")
	}

	opt := &CallOption{}

	for _, optfunc := range opts {
		optfunc(opt)
	}

	if opt.traceId != "" {
		reqmsg.SetTraceId(opt.traceId)
	}
	if reqmsg.TraceId() == "" {
		reqmsg.SetTraceId(misc.NewUuid())
	}
	reqmsg.Log().Debug("request message created")

	// save request in pending map
	expire := time.Now().Add(time.Second * 10)
	reqId := reqmsg.MustId()
	wc := &LivecallT{
		Expire:   expire,
		Request:  reqmsg,
		Callback: callback,
	}
	// TODO: assert live pending requests
	self.pendingLivecalls.Store(reqId, wc)
	self.chSendUp <- reqmsg
	return nil
}

func (self *RPCClient) cleanTimeoutLivecalls() {
	now := time.Now()
	var arr []interface{}
	self.pendingLivecalls.Range(func(k, v interface{}) bool {
		wc, _ := v.(*LivecallT)
		if now.After(wc.Expire) {
			arr = append(arr, k)
		}
		return true
	})
	for _, k := range arr {
		if v, ok := self.pendingLivecalls.LoadAndDelete(k); ok {
			wc, _ := v.(*LivecallT)
			errmsg := jsonrpc.ErrTimeout.ToMessage(wc.Request)
			wc.Callback(errmsg)
		}
	}
}

func (self *RPCClient) handleLiveResult(res jsonrpc.IMessage) {
	reqId := res.MustId()
	if r, ok := self.pendingLivecalls.LoadAndDelete(reqId); ok {
		wc, _ := r.(*LivecallT)
		wc.Callback(res)
	} else {
		res.Log().Warnf("res not found in pendingLivecalls")
	}
}
