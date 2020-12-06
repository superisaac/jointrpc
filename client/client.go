package client

import (
	"context"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	server "github.com/superisaac/rpctube/server"
	grpc "google.golang.org/grpc"
	"log"
)

func NewRPCClient(serverAddress string) *RPCClient {
	methodHandlers := make(map[string](Handler))
	return &RPCClient{ServerAddress: serverAddress, methodHandlers: methodHandlers}

}

func (self *RPCClient) Connect() error {
	conn, err := grpc.Dial(self.ServerAddress, grpc.WithInsecure())
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
			return nil, nil
		} else {
			resmsg := jsonrpc.NewResultMessage(msg.Id, res)
			return resmsg, nil
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
	up_pac := &intf.JSONRPCUpPacket{Payload: payload}
	stream.Send(up_pac)
}

func (self *RPCClient) Handle(method string, handler Handler) {
	//h, ok := self.methodHandlers[method]
	self.methodHandlers[method] = handler
}

func (self *RPCClient) HandleMethods() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := self.TubeClient.Handle(ctx)
	if err != nil {
		return err
	}

	// register methods first
	self.registerMethods(stream)
	for {
		down_pac, err := stream.Recv()
		if err != nil {
			log.Printf("down pack error%v", err)
			return err
		}

		// On Ping
		ping := down_pac.GetPing()
		if ping != nil {
			// Send Pong
			pong := &intf.PONG{Text: ping.Text}
			payload := &intf.JSONRPCUpPacket_Pong{Pong: pong}
			up_pac := &intf.JSONRPCUpPacket{Payload: payload}

			stream.Send(up_pac)
			continue
		}

		// Handle JSONRPC Request
		req := down_pac.GetRequest()
		if req != nil {
			msg, err := server.RequestToMessage(req)
			if err != nil {
				return err
			}
			resmsg, err := self.handleRequestMsg(msg)
			if err != nil {
				return err
			}
			if resmsg != nil {
				rst, err := server.MessageToResult(resmsg)
				if err != nil {
					return err
				}
				payload := &intf.JSONRPCUpPacket_Result{Result: rst}
				up_pac := &intf.JSONRPCUpPacket{Payload: payload}

				stream.Send(up_pac)
			}
			continue
		}
	}
	return nil
}
