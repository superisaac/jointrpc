package client

import (
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

type RPCRequest struct {
	Message *jsonrpc.RPCMessage
	// TODO: add more fields
}

type HandlerFunc func(req *RPCRequest, params []interface{}) (interface{}, error)
type DefaultHandlerFunc func(req *RPCRequest, method string, params []interface{}) (interface{}, error)

type MethodHandler struct {
	function HandlerFunc
	//Options HandlerOption
	schema     string
	concurrent bool
}

type RPCClient struct {
	ServerAddress     string
	TubeClient        intf.JSONRPCTubeClient
	methodHandlers    map[string]MethodHandler
	defaultHandler    DefaultHandlerFunc
	defaultConcurrent bool
	sendUpChannel     chan *intf.JSONRPCUpPacket
}
