package bridge

import (
	client "github.com/superisaac/rpctube/client"
	misc "github.com/superisaac/rpctube/misc"
	tube "github.com/superisaac/rpctube/tube"
)

type Edge struct {
	remoteClient *client.RPCClient
	// set of names
	methodNames misc.StringSet
}

type CmdStateChange struct {
	serverAddress string
	state         *tube.TubeState
}

type Bridge struct {
	serverEntries []client.ServerEntry
	edges         map[string]*Edge
	ChState       chan CmdStateChange
	methodSig     string
}
