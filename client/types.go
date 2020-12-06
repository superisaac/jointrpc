package client

import (
	//intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

type RPCRequest struct {
	Message *jsonrpc.RPCMessage
}

type Handler func(req *RPCRequest, params []interface{}) (interface{}, error)

type RPCClient struct {
	MethodHandlers map[string]Handler
}
