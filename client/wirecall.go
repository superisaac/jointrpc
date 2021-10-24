package client

import (
	"errors"
	"time"
	//"fmt"
	"context"
	//"github.com/superisaac/jointrpc/encoding"
	//intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/misc"

	log "github.com/sirupsen/logrus"
)

func (self *RPCClient) CallInStream(rootCtx context.Context, reqmsg jsonrpc.IMessage, callback WireCallback, opts ...CallOptionFunc) error {
	//misc.Assert(self.workerStream != nil, "worker steam is empty")
	if self.workerStream == nil {
		log.Warnf("worker stream is empty")
		return errors.New("worker stream is empty")
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
	expire := time.Now().Add(time.Second * 30)
	reqId := reqmsg.MustId()
	wc := WireCallT{
		Expire:   expire,
		Callback: callback,
	}
	// TODO: assert wire pending requests
	self.wirePendingRequests[reqId] = wc
	self.chSendUp <- reqmsg
	return nil
}

func (self *RPCClient) handleWireResult(res jsonrpc.IMessage) {
	reqId := res.MustId()
	if wc, ok := self.wirePendingRequests[reqId]; ok {
		delete(self.wirePendingRequests, reqId)
		// callback
		wc.Callback(res)
	} else {
		res.Log().Warnf("res not found in wirePendingRequests")
	}
}
