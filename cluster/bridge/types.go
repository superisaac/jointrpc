package bridge

import (
	client "github.com/superisaac/jointrpc/client"
	misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
)

type Edge struct {
	remoteClient *client.RPCClient
	// set of names
	methodNames     misc.StringSet
	delegateMethods []string
}

type CmdStateChange struct {
	serverUrl string
	state     *rpcrouter.ServerState
}

type Bridge struct {
	serverEntries []client.ServerEntry
	edges         map[string]*Edge
	ChState       chan CmdStateChange
	methodSig     string
}
