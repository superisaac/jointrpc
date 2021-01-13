package client

import (
	"context"
	uuid "github.com/google/uuid"
	//simplejson "github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	encoding "github.com/superisaac/jointrpc/encoding"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
)

type CallOption struct {
	broadcast bool
	traceId   string
}

type CallOptionFunc func(opt *CallOption)

func WithTraceId(traceId string) CallOptionFunc {
	return func(opt *CallOption) {
		opt.traceId = traceId
	}
}

func WithBroadcast(b bool) CallOptionFunc {
	return func(opt *CallOption) {
		opt.broadcast = b
	}
}

func (self *RPCClient) CallRPC(rootCtx context.Context, method string, params []interface{}, opts ...CallOptionFunc) (jsonrpc.IMessage, error) {
	msgId := 1

	reqmsg := jsonrpc.NewRequestMessage(msgId, method, params, nil)
	return self.CallMessage(rootCtx, reqmsg, opts...)
}

func (self *RPCClient) CallMessage(rootCtx context.Context, reqmsg jsonrpc.IMessage, opts ...CallOptionFunc) (jsonrpc.IMessage, error) {

	opt := &CallOption{}

	for _, optfunc := range opts {
		optfunc(opt)
	}

	if opt.traceId != "" {
		reqmsg.SetTraceId(opt.traceId)
	}
	if reqmsg.TraceId() == "" {
		reqmsg.SetTraceId(uuid.New().String())
	}
	reqmsg.Log().Debug("request message created")
	envolope := encoding.MessageToEnvolope(reqmsg)
	req := &intf.JSONRPCCallRequest{Envolope: envolope, Broadcast: opt.broadcast}

	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	res, err := self.tubeClient.Call(ctx, req)
	if err != nil {
		return nil, err
	}

	resmsg, err := jsonrpc.ParseBytes([]byte(res.Envolope.Body))

	if err != nil {
		return nil, err
	}
	resmsg.SetTraceId(res.Envolope.TraceId)
	if !resmsg.IsResultOrError() {
		log.Warnf("bad result or error message %+v", res.Envolope.Body)
		return nil, &jsonrpc.RPCError{10409, "msg is neither result nor error", false}
	}
	return resmsg, nil
}

func (self *RPCClient) ListDelegates(rootCtx context.Context) ([]string, error) {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	req := &intf.ListDelegatesRequest{}
	res, err := self.tubeClient.ListDelegates(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Delegates, err
}

func (self *RPCClient) SendNotify(rootCtx context.Context, method string, params []interface{}, opts ...CallOptionFunc) error {
	opt := &CallOption{}
	for _, optfunc := range opts {
		optfunc(opt)
	}
	notify := jsonrpc.NewNotifyMessage(method, params, nil)

	if opt.traceId == "" {
		opt.traceId = uuid.New().String()
	}

	notify.SetTraceId(opt.traceId)

	env := &intf.JSONRPCEnvolope{
		Body:    notify.MustString(),
		TraceId: opt.traceId}
	req := &intf.JSONRPCNotifyRequest{
		Envolope:  env,
		Broadcast: opt.broadcast,
	}
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	notify.Log().Debug("notify message created")
	resp, err := self.tubeClient.Notify(ctx, req)
	if err != nil {
		return err
	}
	notify.Log().Debugf("notify result %s", resp.Text)
	return nil
}
