package mirror

import (
	client "github.com/superisaac/rpctube/client"
	misc "github.com/superisaac/rpctube/misc"
	tube "github.com/superisaac/rpctube/tube"
	handler "github.com/superisaac/rpctube/tube/handler"
)

type Edge struct {
	remoteClient *client.RPCClient
	// set of names
	dlgMethods  []tube.MethodInfo
	methodNames misc.StringSet
}

type CmdStateChange struct {
	ServerAddress string
	State         *tube.TubeState
}

type Mirror struct {
	handler.HandlerManager
	conn          *tube.ConnT
	router        *tube.Router
	serverEntries []client.ServerEntry
	edges         map[string]*Edge
	//methodEdges map[string]StringSet
	ChState   chan CmdStateChange
	methodSig string
}
