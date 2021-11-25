package server

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/superisaac/jointrpc/dispatch"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jsonrpc"
)

type StreamDispatcher struct {
	disp     *dispatch.Dispatcher
	authDisp *dispatch.Dispatcher
}

type IReceiver interface {
	Recv() (jsonrpc.IMessage, error)
}

// Websocket server
type WSServer struct {
	//router *rpcrouter.Router
	rootCtx context.Context
}

type WSAdaptor struct {
	ws   *websocket.Conn
	done chan error
}

// gRPC server
type JointRPC struct {
	intf.UnimplementedJointRPCServer
}

type GRPCAdaptor struct {
	stream intf.JointRPC_LiveServer
	done   chan error
}
