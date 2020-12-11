package client

import (
	"context"
	simplejson "github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	server "github.com/superisaac/rpctube/server"
)

func (self *RPCClient) CallRPC(method string, params []interface{}) (*jsonrpc.RPCMessage, error) {
	log.Infof("log methods %s, params %v", method, params)
	paramsJson := simplejson.New()
	paramsJson.SetPath(nil, params)
	paramsStr, err := jsonrpc.MarshalJson(paramsJson)
	if err != nil {
		return nil, err
	}
	req := &intf.JSONRPCRequest{
		Id:     "1",
		Method: method,
		Params: paramsStr}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := self.TubeClient.Call(ctx, req)
	log.Infof("res is %v", res)
	if err != nil {
		return nil, err
	}

	msg, err := server.ResultToMessage(res)
	if err != nil {
		return nil, err
	}
	return msg, err
}

func (self *RPCClient) ListMethods() ([]string, error) {
	req := &intf.ListMethodsRequest{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := self.TubeClient.ListMethods(ctx, req)
	if err != nil {
		return []string{}, err
	}

	return res.Methods, nil
}
