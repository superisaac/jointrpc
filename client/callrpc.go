package client

import (
	"context"
	//simplejson "github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	//server "github.com/superisaac/jointrpc/server"
)

func (self *RPCClient) CallRPC(rootCtx context.Context, method string, params []interface{}) (jsonrpc.IMessage, error) {
	msgId := 1

	msg := jsonrpc.NewRequestMessage(msgId, method, params, nil)
	return self.CallMessage(rootCtx, msg)
}

func (self *RPCClient) CallMessage(rootCtx context.Context, msg jsonrpc.IMessage) (jsonrpc.IMessage, error) {
	envolope := &intf.JSONRPCEnvolope{Body: msg.MustString()}
	req := &intf.JSONRPCCallRequest{Envolope: envolope}

	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	res, err := self.tubeClient.Call(ctx, req)
	log.Debugf("res is %v", res)
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

func (self *RPCClient) SendNotify(rootCtx context.Context, method string, params []interface{}, broadcast bool) error {
	notify := jsonrpc.NewNotifyMessage(method, params, nil)

	env := &intf.JSONRPCEnvolope{Body: notify.MustString()}
	req := &intf.JSONRPCNotifyRequest{
		Envolope:  env,
		Broadcast: broadcast,
	}
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()
	res, err := self.tubeClient.Notify(ctx, req)
	if err != nil {
		return err
	}
	log.Debugf("send notify result %s", res.Text)
	return nil
}
