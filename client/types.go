package client

import (
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	handler "github.com/superisaac/jointrpc/rpcrouter/handler"
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

type RPCClient struct {
	handler.HandlerManager

	serverEntry   ServerEntry
	tubeClient    intf.JointRPCClient
	sendUpChannel chan *intf.JointRPCUpPacket
}

type MethodUpdateReceiver chan []*intf.MethodInfo
