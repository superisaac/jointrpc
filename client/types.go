package client

import (
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	handler "github.com/superisaac/rpctube/tube/handler"
)

type RPCRequest struct {
	Message *jsonrpc.RPCMessage
	// TODO: add more fields
}

type RPCClient struct {
	handler.HandlerManager
	serverAddress string
	certFile      string
	tubeClient    intf.JSONRPCTubeClient
	sendUpChannel chan *intf.JSONRPCUpPacket
}
