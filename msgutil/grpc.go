package msgutil

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/jsonrpc"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"io"
)

func GRPCClientSend(stream intf.JointRPC_LiveClient, msg jsonrpc.IMessage) error {
	envo := MessageToEnvolope(msg)
	return stream.Send(envo)
}

func GRPCClientRecv(stream intf.JointRPC_LiveClient) (jsonrpc.IMessage, error) {
	envo, err := stream.Recv()
	if err != nil {
		return nil, errors.Wrap(err, "stream.Recv()")
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
		return nil, errors.Wrap(err, "stream.Recv()")
	}

	msg, err := MessageFromEnvolope(envo)
	return msg, err
}

func GRPCHandleCodes(err error, safeCodes ...codes.Code) error {
	if errors.Is(err, io.EOF) {
		log.Infof("cannot connect stream")
		return nil
	}
	code := grpc.Code(err)
	if code == codes.Unknown {
		cause := errors.Cause(err)
		if cause != nil {
			code = grpc.Code(cause)
		}
	}
	for _, safeCode := range safeCodes {
		if code == safeCode {
			log.Debugf("grpc connect code %d %s", code, code)
			return nil
		}
	}
	// if code == codes.Unavailable { // || code == codes.Canceled {
	// 	log.Debugf("connect closed %d retrying", code)
	// 	return nil
	// }
	log.Warnf("error on handle %+v", err)
	return err
}
