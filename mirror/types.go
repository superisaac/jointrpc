package mirror

import (
	client "github.com/superisaac/jointrpc/client"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	handler "github.com/superisaac/jointrpc/rpcrouter/handler"
)

type Edge struct {
	remoteClient *client.RPCClient
	// set of names
	dlgMethods  []rpcrouter.MethodInfo
	methodNames misc.StringSet
}

type CmdStateChange struct {
	ServerUrl string
	State     *rpcrouter.TubeState
}

type Mirror struct {
	handler.HandlerManager
	conn          *rpcrouter.ConnT
	router        *rpcrouter.Router
	serverEntries []client.ServerEntry
	edges         map[string]*Edge
	//methodEdges map[string]StringSet
	ChState   chan CmdStateChange
	methodSig string
}
