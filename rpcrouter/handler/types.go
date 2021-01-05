package handler

import (
	//"context"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/rpcrouter"
)

type RPCRequest struct {
	MsgVec rpcrouter.MsgVec
}

type HandlerFunc func(req *RPCRequest, params []interface{}) (interface{}, error)
type DefaultHandlerFunc func(req *RPCRequest, method string, params []interface{}) (interface{}, error)

// listen to tube state change
type StateHandlerFunc func(newState *rpcrouter.TubeState)

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
