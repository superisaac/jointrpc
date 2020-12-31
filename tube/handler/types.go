package handler

import (
	//"github.com/gorilla/websocket"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	tube "github.com/superisaac/rpctube/tube"
)

type RPCRequest struct {
	MsgVec tube.MsgVec
	//Message jsonrpc.IMessage
	//FromConnId tube.CID
	// TODO: add more fields
}

type HandlerFunc func(req *RPCRequest, params []interface{}) (interface{}, error)
type DefaultHandlerFunc func(req *RPCRequest, method string, params []interface{}) (interface{}, error)

// listen to tube state change
type StateHandlerFunc func(newState *tube.TubeState)

type OnChangeFunc func()

type MethodHandler struct {
	function HandlerFunc
	//Options HandlerOption
	SchemaJson string
	Help       string
	Concurrent bool
}

type HandlerManager struct {
	ChResultMsg    chan jsonrpc.IMessage
	MethodHandlers map[string]MethodHandler
	StateHandler   StateHandlerFunc

	defaultHandler    DefaultHandlerFunc
	defaultConcurrent bool
	onChange          OnChangeFunc
}
