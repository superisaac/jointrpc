package neighbor

import (
	client "github.com/superisaac/jointrpc/client"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	dispatch "github.com/superisaac/jointrpc/dispatch"
)

type Edge struct {
	remoteClient *client.RPCClient
	// set of names
	dlgMethods  []rpcrouter.MethodInfo
	methodNames misc.StringSet
}

type CmdStateChange struct {
	ServerUrl string
	State     *rpcrouter.ServerState
}

type NeighborService struct {
	dispatch.Dispatcher
	conn          *rpcrouter.ConnT
	router        *rpcrouter.Router
	serverEntries []client.ServerEntry
	edges         map[string]*Edge
	//methodEdges map[string]StringSet
	ChState   chan CmdStateChange
	methodSig string
}
