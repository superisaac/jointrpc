package dispatch

import (
	"context"
	//"errors"
	"github.com/pkg/errors"
	"github.com/superisaac/jointrpc/rpcrouter"
	"github.com/superisaac/jsonz"
)

type RPCRequest struct {
	Context context.Context
	CmdMsg  rpcrouter.CmdMsg
	Data    interface{} // user defined data
}

var (
	Deferred = errors.New("dispatch deferred")
)

type HandlerFunc func(req *RPCRequest, params []interface{}) (interface{}, error)
type DefaultHandlerFunc func(req *RPCRequest, method string, params []interface{}) (interface{}, error)

type OnChangeFunc func()

type MethodHandler struct {
	function HandlerFunc
	//Options HandlerOption
	SchemaJson string
	Help       string
}

type HandlerOption func(*MethodHandler)

type ResultT struct {
	ResMsg    jsonz.Message
	ReqCmdMsg rpcrouter.CmdMsg
}

type Dispatcher struct {
	//ChResult       chan ResultT
	spawnExec      bool
	methodHandlers map[string]MethodHandler
	defaultHandler DefaultHandlerFunc
	changeHandlers []OnChangeFunc
}

// listen to tube state change
type StateHandlerFunc func(newState *rpcrouter.ServerState)
type StateListener struct {
	stateHandlers []StateHandlerFunc
}

type ISender interface {
	SendMessage(ctx context.Context, msg jsonz.Message) error
	SendCmdMsg(ctx context.Context, cmdMsg rpcrouter.CmdMsg) error
	Done() chan error
}
