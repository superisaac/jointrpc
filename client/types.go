package client

import (
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"time"
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

type WireCallback func(jsonrpc.IMessage)

type WireCallT struct {
	Expire   time.Time
	Callback WireCallback
}

type ConnectionLostCallback func()
type ConnectedCallback func()

type RPCClient struct {
	workerStream  intf.JointRPC_WorkerClient
	stateStream   intf.JointRPC_SubscribeStateClient
	serverEntry   ServerEntry
	serverUrl     *url.URL
	connected     bool
	grpcClient    intf.JointRPCClient
	sendUpChannel chan *intf.JointRPCUpPacket

	WorkerRetryTimes int
	onConnected      ConnectedCallback
	onConnectionLost ConnectionLostCallback

	wirePendingRequests map[interface{}]WireCallT
}

type MethodUpdateReceiver chan []*intf.MethodInfo
