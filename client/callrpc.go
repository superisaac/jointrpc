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

	var msgId string = "1"
	req := &intf.JSONRPCRequest{
		Id:     msgId,
		Method: method,
		Params: paramsStr}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := self.tubeClient.Call(ctx, req)
	log.Debugf("res is %v", res)
	if err != nil {
		return nil, err
	}

	msg, err := server.ResultToMessage(res)
	if err != nil {
		return nil, err
	}
	return msg, err
}

func (self *RPCClient) SendNotify(method string, params []interface{}, broadcast bool) error {
	log.Infof("log methods %s, params %v", method, params)
	paramsJson := simplejson.New()
	paramsJson.SetPath(nil, params)
	paramsStr, err := jsonrpc.MarshalJson(paramsJson)
	if err != nil {
		return err
	}

	req := &intf.JSONRPCNotifyRequest{
		Method:    method,
		Params:    paramsStr,
		Broadcast: broadcast,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := self.tubeClient.Notify(ctx, req)

	if err != nil {
		return err
	}
	log.Debugf("not ify res is %v", res)
	return nil
}

func (self *RPCClient) ListMethods() ([]*intf.MethodInfo, error) {
	req := &intf.ListMethodsRequest{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := self.tubeClient.ListMethods(ctx, req)
	if err != nil {
		return [](*intf.MethodInfo){}, err
	}

	return res.MethodInfos, nil
}
