package bridge

import (
	client "github.com/superisaac/rpctube/client"
	tube "github.com/superisaac/rpctube/tube"
	handler "github.com/superisaac/rpctube/tube/handler"
)

type StringSet map[string]bool

type Edge struct {
	remoteClient *client.RPCClient
	// set of names
	dlgMethods  []tube.MethodInfo
	methodNames StringSet
}

type CmdStateChange struct {
	ServerAddress string
	State         *tube.TubeState
}

type Bridge struct {
	handler.HandlerManager
	conn          *tube.ConnT
	serverEntries []client.ServerEntry
	edges         map[string]*Edge
	//methodEdges map[string]StringSet
	ChState   chan CmdStateChange
	methodSig string
}
