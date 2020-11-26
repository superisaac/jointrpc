package client

import (
	"context"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	server "github.com/superisaac/rpctube/server"
	simplejson "github.com/bitly/go-simplejson"	
)

func CallRPC(c intf.JSONRPCTubeClient, method string, params []interface{}) (*jsonrpc.RPCMessage, error) {
	paramsJson := simplejson.New()
	paramsJson.SetPath(nil, params)
	paramsStr, err := jsonrpc.MarshalJson(paramsJson)
	if err != nil {
		return nil, err
	}
	req := &intf.JSONRPCRequest{
		Id: 1,
		Method: method,
		Params: paramsStr}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := c.Call(ctx, req)
	if err != nil {
		return nil, err
	}

	msg, err := server.ResultToMessage(res)
	if err != nil {
		return nil, err
	}
	return msg, err
}
