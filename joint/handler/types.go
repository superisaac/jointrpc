package handler

import (
	//"context"
	"github.com/superisaac/jointrpc/joint"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
)

type RPCRequest struct {
	MsgVec joint.MsgVec
}

type HandlerFunc func(req *RPCRequest, params []interface{}) (interface{}, error)
type DefaultHandlerFunc func(req *RPCRequest, method string, params []interface{}) (interface{}, error)

// listen to tube state change
type StateHandlerFunc func(newState *joint.TubeState)

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
