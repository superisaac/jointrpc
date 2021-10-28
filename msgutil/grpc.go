package msgutil

import (
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/jsonrpc"
)

func GRPCClientSend(stream intf.JointRPC_LiveClient, msg jsonrpc.IMessage) error {
	envo := MessageToEnvolope(msg)
	return stream.Send(envo)
}

func GRPCClientRecv(stream intf.JointRPC_LiveClient) (jsonrpc.IMessage, error) {
	envo, err := stream.Recv()
	if err != nil {
		return nil, err
	}

	msg, err := MessageFromEnvolope(envo)
	return msg, err
}

func GRPCServerSend(stream intf.JointRPC_LiveServer, msg jsonrpc.IMessage) error {
	envo := MessageToEnvolope(msg)
	return stream.Send(envo)
}

func GRPCServerRecv(stream intf.JointRPC_LiveServer) (jsonrpc.IMessage, error) {
	envo, err := stream.Recv()
	if err != nil {
		return nil, err
	}

	msg, err := MessageFromEnvolope(envo)
	return msg, err
}
