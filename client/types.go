package client

import (
	//intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

type Handler func(msg *jsonrpc.RPCMessage) (*jsonrpc.RPCMessage, error)

type RPCClient struct {
	MethodHandlers map[string]Handler
}
