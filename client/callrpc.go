package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	//simplejson "github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/msgutil"
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
	reqmsg := jsonrpc.NewRequestMessage(msgId, method, params)
	return self.CallMessage(rootCtx, reqmsg, opts...)
}

func (self *RPCClient) CallMessage(rootCtx context.Context, reqmsg jsonrpc.IMessage, opts ...CallOptionFunc) (jsonrpc.IMessage, error) {
	if self.IsHttp() {
		return self.CallHTTPMessage(rootCtx, reqmsg, opts...)
	} else {
		return self.CallGRPCMessage(rootCtx, reqmsg, opts...)
	}
}

func (self *RPCClient) CallHTTPMessage(rootCtx context.Context, reqmsg jsonrpc.IMessage, opts ...CallOptionFunc) (jsonrpc.IMessage, error) {
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
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	marshaled, err := jsonrpc.MessageBytes(reqmsg)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(marshaled)
	req, err := http.NewRequestWithContext(ctx, "POST", self.serverEntry.ServerUrl, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Trace-Id", reqmsg.TraceId())
	// TODO: handle broadcast
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("bad resp %d", resp.StatusCode))
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if reqmsg.IsRequest() {
		respMsg, err := jsonrpc.ParseBytes(respBody)
		if err != nil {
			return nil, err
		}
		respMsg.SetTraceId(resp.Header.Get("X-Trace-Id"))
		return respMsg, nil
	} else {
		return nil, nil
	}
}

func (self *RPCClient) CallGRPCMessage(rootCtx context.Context, reqmsg jsonrpc.IMessage, opts ...CallOptionFunc) (jsonrpc.IMessage, error) {

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
	envolope := msgutil.MessageToEnvolope(reqmsg)
	req := &intf.JSONRPCCallRequest{
		Auth:      self.ClientAuth(),
		Envolope:  envolope,
		Broadcast: opt.broadcast,
	}
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	res, err := self.grpcClient.Call(ctx, req)
	if err != nil {
		return nil, err
	}
	err = self.CheckStatus(res.Status, "Call")
	if err != nil {
		return nil, err
	}

	resmsg, err := jsonrpc.ParseBytes([]byte(res.Envolope.Body))

	if err != nil {
		return nil, err
	}
	if !resmsg.IsResultOrError() {
		log.Warnf("bad result or error message %+v", res.Envolope.Body)
		return nil, &jsonrpc.RPCError{10409, "msg is neither result nor error", false}
	}
	return resmsg, nil
}

func (self *RPCClient) SendNotify(rootCtx context.Context, method string, params []interface{}, opts ...CallOptionFunc) error {
	if self.IsHttp() {
		return self.SendHTTPNotify(rootCtx, method, params, opts...)
	} else {
		return self.SendGRPCNotify(rootCtx, method, params, opts...)
	}
}

func (self *RPCClient) SendHTTPNotify(rootCtx context.Context, method string, params []interface{}, opts ...CallOptionFunc) error {
	opt := &CallOption{}
	for _, optfunc := range opts {
		optfunc(opt)
	}
	notify := jsonrpc.NewNotifyMessage(method, params)

	if opt.traceId == "" {
		opt.traceId = misc.NewUuid()
	}

	notify.SetTraceId(opt.traceId)

	marshaled, err := jsonrpc.MessageBytes(notify)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(marshaled)
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", self.serverEntry.ServerUrl, reader)
	if err != nil {
		return err
	}
	req.Header.Add("X-Trace-Id", notify.TraceId())
	// TODO: handle broadcast
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("bad resp %d", resp.StatusCode))
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (self *RPCClient) SendGRPCNotify(rootCtx context.Context, method string, params []interface{}, opts ...CallOptionFunc) error {
	opt := &CallOption{}
	for _, optfunc := range opts {
		optfunc(opt)
	}
	notify := jsonrpc.NewNotifyMessage(method, params)

	if opt.traceId == "" {
		opt.traceId = misc.NewUuid()
	}

	notify.SetTraceId(opt.traceId)

	env := &intf.JSONRPCEnvolope{
		Body: jsonrpc.MessageString(notify),
	}
	//		TraceId: opt.traceId}
	req := &intf.JSONRPCNotifyRequest{
		Auth:      self.ClientAuth(),
		Envolope:  env,
		Broadcast: opt.broadcast,
	}
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	notify.Log().Debug("notify message created")
	res, err := self.grpcClient.Notify(ctx, req)
	if err != nil {
		return err
	}

	err = self.CheckStatus(res.Status, "Notify")
	if err != nil {
		notify.Log().Warn(err.Error())
		return err
	}

	notify.Log().Debugf("notify result %s", res.Text)
	return nil
}
