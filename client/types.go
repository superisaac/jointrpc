package client

import (
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	//dispatch "github.com/superisaac/jointrpc/dispatch"
	"net/url"
)

type ServerEntry struct {
	ServerUrl string // raw url, scheme must be h2 or h2c
	CertFile  string
}

type ServerFlag struct {
	pAddress  *string
	pCertFile *string
}

type RPCRequest struct {
	Message jsonrpc.IMessage
	// TODO: add more fields
}

type RPCStatusError struct {
	Method string
	Code   int
	Reason string
}

type ConnectionLostCallback func()
type ConnectedCallback func()

type RPCClient struct {
	//disp *dispatch.Dispatcher

	//connPublicId  string
	workerStream  intf.JointRPC_WorkerClient
	serverEntry   ServerEntry
	serverUrl     *url.URL
	connected     bool
	rpcClient     intf.JointRPCClient
	sendUpChannel chan *intf.JointRPCUpPacket

	onConnected      ConnectedCallback
	onConnectionLost ConnectionLostCallback
}

type MethodUpdateReceiver chan []*intf.MethodInfo
