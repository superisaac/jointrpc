package client

import (
	"github.com/superisaac/jointrpc/dispatch"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/jsonrpc"
	"net/url"
	"sync"
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
type AuthorizedCallback func()

type RPCClient struct {
	serverEntry ServerEntry
	serverUrl   *url.URL
	connected   bool

	// grpc transport
	grpcClient intf.JointRPCClient

	chSendUp chan jsonrpc.IMessage

	LiveRetryTimes int
	retry          int

	onConnected      ConnectedCallback
	onConnectionLost ConnectionLostCallback
	onAuthorized     AuthorizedCallback

	//wirePendingRequests map[interface{}]WireCallT
	wirePendingRequests sync.Map

	chResult chan dispatch.ResultT
}

type MethodUpdateReceiver chan []*intf.MethodInfo
