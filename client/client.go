package client

import (
	"context"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	server "github.com/superisaac/rpctube/server"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"io"
	"log"
	"time"
	//transport "google.golang.org/grpc/internal/transport"
)

func NewRPCClient(serverAddress string) *RPCClient {
	methodHandlers := make(map[string](Handler))
	return &RPCClient{ServerAddress: serverAddress, methodHandlers: methodHandlers}

}

func (self *RPCClient) Connect() error {
	conn, err := grpc.Dial(self.ServerAddress,
		grpc.WithInsecure())
	if err != nil {
		return err
	}
	self.TubeClient = intf.NewJSONRPCTubeClient(conn)
	return nil
}

func (self *RPCClient) handleRequestMsg(msg *jsonrpc.RPCMessage) (*jsonrpc.RPCMessage, error) {
	handler, ok := self.methodHandlers[msg.Method]
	if ok {
		req := &RPCRequest{Message: msg}

		params := msg.Params.MustArray()

		res, err := handler(req, params)
		if err != nil {
			return nil, err
		} else if msg.IsRequest() {
			resmsg := jsonrpc.NewResultMessage(msg.Id, res)
			return resmsg, nil
		} else {
			return nil, nil
		}

	} else {
		errmsg := jsonrpc.NewErrorMessage(msg.Id, 404, "no such message", false)
		return errmsg, nil
	}
}

func (self *RPCClient) registerMethods(stream intf.JSONRPCTube_HandleClient) {
	methods := make([]string, 0, len(self.methodHandlers))
	for method := range self.methodHandlers {
		methods = append(methods, method)
	}
	reg := &intf.RegisterMethodsRequest{Methods: methods}
	payload := &intf.JSONRPCUpPacket_RegisterMethods{RegisterMethods: reg}
	uppac := &intf.JSONRPCUpPacket{Payload: payload}
	stream.Send(uppac)
}

func (self *RPCClient) On(method string, handler Handler) {
	//h, ok := self.methodHandlers[method]
	self.methodHandlers[method] = handler
}

func (self *RPCClient) HandleRPC() error {
	for {
		err := self.handleRPC()
		log.Printf("sss %+v", err)
		//if err == transport.ErrConnClosing {

		//if err != nil && err.Error() == "transport is closing" {
		//if false {
		if err != nil {
			if grpc.Code(err) == codes.Unavailable {
				log.Printf("connect closed retrying")
			} else {
				return err
			}
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (self *RPCClient) handleRPC() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := self.TubeClient.Handle(ctx, grpc_retry.WithMax(500))
	if err != nil {
		return err
	}

	// register methods first
	self.registerMethods(stream)
	for {
		downpac, err := stream.Recv()
		if err == io.EOF {
			log.Printf("eor close")
			return nil
		}
		if err != nil {
			log.Printf("down pack error %+v %d", err, grpc.Code(err))
			return err
		}

		// On Ping
		ping := downpac.GetPing()
		if ping != nil {
			// Send Pong
			pong := &intf.PONG{Text: ping.Text}
			payload := &intf.JSONRPCUpPacket_Pong{Pong: pong}
			uppac := &intf.JSONRPCUpPacket{Payload: payload}

			stream.Send(uppac)
			continue
		}

		// Handle JSONRPC Request
		req := downpac.GetRequest()
		if req != nil {
			msg, err := server.RequestToMessage(req)
			if err != nil {
				return err
			}
			resmsg, err := self.handleRequestMsg(msg)
			if err != nil {
				return err
			}
			if r := recover(); r != nil {
				log.Fatalf("error on handling request msg %+v", r)
				return nil
			}
			if resmsg != nil {
				rst, err := server.MessageToResult(resmsg)
				if err != nil {
					return err
				}
				payload := &intf.JSONRPCUpPacket_Result{Result: rst}
				uppac := &intf.JSONRPCUpPacket{Payload: payload}

				stream.Send(uppac)
			}
			continue
		}
	}
	return nil
}
