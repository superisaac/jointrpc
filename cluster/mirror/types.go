package mirror

import (
	client "github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/joint"
	handler "github.com/superisaac/jointrpc/joint/handler"
	misc "github.com/superisaac/jointrpc/misc"
)

type Edge struct {
	remoteClient *client.RPCClient
	// set of names
	dlgMethods  []joint.MethodInfo
	methodNames misc.StringSet
}

type CmdStateChange struct {
	ServerAddress string
	State         *joint.TubeState
}

type Mirror struct {
	handler.HandlerManager
	conn          *joint.ConnT
	router        *joint.Router
	serverEntries []client.ServerEntry
	edges         map[string]*Edge
	//methodEdges map[string]StringSet
	ChState   chan CmdStateChange
	methodSig string
}
