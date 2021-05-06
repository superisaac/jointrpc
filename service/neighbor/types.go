package neighbor

import (
	client "github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/dispatch"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
)

type Edge struct {
	remoteClient *client.RPCClient
	//disp         *dispatch.Dispatcher
	stateDisp *dispatch.StateDispatcher
	// set of names
	dlgMethods  []rpcrouter.MethodInfo
	methodNames misc.StringSet
}

type CmdStateChange struct {
	ServerUrl string
	State     *rpcrouter.ServerState
}

type NeighborService struct {
	dispatcher    *dispatch.Dispatcher
	conn          *rpcrouter.ConnT
	router        *rpcrouter.Router
	serverEntries []client.ServerEntry
	edges         map[string]*Edge
	//methodEdges map[string]StringSet
	ChState   chan CmdStateChange
	methodSig string
}
