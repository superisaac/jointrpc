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
	stateListener *dispatch.StateListener
	// set of names
	dlgMethods  []rpcrouter.MethodInfo
	methodNames misc.StringSet
}

type NamedServerEntry struct {
	namespace string
	entry     client.ServerEntry
}

type CmdStateChange struct {
	ServerUrl string
	State     *rpcrouter.ServerState
}

type NeighborPort struct {
	dispatcher    *dispatch.Dispatcher
	namespace     string
	edges         map[string]*Edge
	ChState       chan CmdStateChange
	methodSig     string
	serverEntries []client.ServerEntry
	conn          *rpcrouter.ConnT
}

type NeighborService struct {

	//router        *rpcrouter.Router
	//namedServerEntries []NamedServerEntry

	ports map[string]*NeighborPort
	//methodEdges map[string]StringSet

}
