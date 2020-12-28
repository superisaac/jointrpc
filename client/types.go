package client

import (
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	handler "github.com/superisaac/rpctube/tube/handler"
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
	//serverAddress string
	//certFile      string
	serverEntry   ServerEntry
	tubeClient    intf.JSONRPCTubeClient
	sendUpChannel chan *intf.JSONRPCUpPacket
}

type MethodUpdateReceiver chan []*intf.MethodInfo
