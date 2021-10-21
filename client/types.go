package client

import (
	"github.com/superisaac/jointrpc/dispatch"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/jsonrpc"
	"net/url"
	"time"
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
	workerStream intf.JointRPC_WorkerClient
	serverEntry  ServerEntry
	serverUrl    *url.URL
	connected    bool
	grpcClient   intf.JointRPCClient
	chSendUp     chan jsonrpc.IMessage

	WorkerRetryTimes int
	onConnected      ConnectedCallback
	onConnectionLost ConnectionLostCallback

	wirePendingRequests map[interface{}]WireCallT

	chResult chan dispatch.ResultT
}

type MethodUpdateReceiver chan []*intf.MethodInfo
