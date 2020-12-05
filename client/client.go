package client

import (
	"context"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	server "github.com/superisaac/rpctube/server"
	"log"
)

func NewRPCClient() *RPCClient {
	methodHandlers := make(map[string](Handler))
	return &RPCClient{MethodHandlers: methodHandlers}
}

func (self *RPCClient) handleRequestMsg(msg *jsonrpc.RPCMessage) (*jsonrpc.RPCMessage, error) {
	handler, ok := self.MethodHandlers[msg.Method]
	if ok {
		return handler(msg)
	} else {
		errmsg := jsonrpc.NewErrorMessage(msg.Id, 404, "no such message", false)
		return errmsg, nil
	}
}

func (self *RPCClient) registerMethods(stream intf.JSONRPCTube_HandleClient) {
	methods := make([]string, 0, len(self.MethodHandlers))
	for method := range self.MethodHandlers {
		methods = append(methods, method)
	}
	reg := &intf.RegisterMethodsRequest{Methods: methods}
	payload := &intf.JSONRPCUpPacket_RegisterMethods{RegisterMethods: reg}
	up_pac := &intf.JSONRPCUpPacket{Payload: payload}
	stream.Send(up_pac)
}

func (self *RPCClient) Handle(method string, handler Handler) {
	//h, ok := self.MethodHandlers[method]
	self.MethodHandlers[method] = handler
}

func (self *RPCClient) HandleMethods(c intf.JSONRPCTubeClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.Handle(ctx)
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
