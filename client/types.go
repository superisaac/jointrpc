package client

import (
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	handler "github.com/superisaac/jointrpc/joint/handler"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
)

type ServerEntry struct {
	Address  string
	CertFile string
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
