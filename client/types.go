package client

import (
	handler "github.com/superisaac/rpctube/handler"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

type RPCRequest struct {
	Message *jsonrpc.RPCMessage
	// TODO: add more fields
}

type RPCClient struct {
	handler.HandlerManager
	ServerAddress string
	TubeClient    intf.JSONRPCTubeClient
	sendUpChannel chan *intf.JSONRPCUpPacket
}
