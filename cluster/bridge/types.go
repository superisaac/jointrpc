package bridge

import (
	client "github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/joint"
	misc "github.com/superisaac/jointrpc/misc"
)

type Edge struct {
	remoteClient *client.RPCClient
	// set of names
	methodNames misc.StringSet
}

type CmdStateChange struct {
	serverAddress string
	state         *joint.TubeState
}

type Bridge struct {
	serverEntries []client.ServerEntry
	edges         map[string]*Edge
	ChState       chan CmdStateChange
	methodSig     string
}
