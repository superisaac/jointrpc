package client

import (
	"context"
	simplejson "github.com/bitly/go-simplejson"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	server "github.com/superisaac/rpctube/server"
	"log"
)

func CallRPC(c intf.JSONRPCTubeClient, method string, params []interface{}) (*jsonrpc.RPCMessage, error) {
	log.Printf("log methods %s, params %v", method, params)
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
	res, err := c.Call(ctx, req)
	log.Printf("res is %v", res)
	if err != nil {
		return nil, err
	}

	msg, err := server.ResultToMessage(res)
	if err != nil {
		return nil, err
	}
	return msg, err
}

func ListMethods(c intf.JSONRPCTubeClient) ([]string, error) {
	req := &intf.ListMethodsRequest{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := c.ListMethods(ctx, req)
	if err != nil {
		return []string{}, err
	}

	return res.Methods, nil
}
