package client

import (
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

type RPCRequest struct {
	Message *jsonrpc.RPCMessage
	// TODO: add more fields
}

type MethodHandler func(req *RPCRequest, params []interface{}) (interface{}, error)
type MsgHandler func(req *RPCRequest, method string, params []interface{}) (interface{}, error)

type RPCClient struct {
	ServerAddress  string
	TubeClient     intf.JSONRPCTubeClient
	methodHandlers map[string]MethodHandler
	defaultHandler MsgHandler
}
