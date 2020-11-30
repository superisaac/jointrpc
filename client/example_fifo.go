package client

import (
	"context"
	"log"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	server "github.com/superisaac/rpctube/server"
	//simplejson "github.com/bitly/go-simplejson"	
)

type Fifo struct {
	elements []interface{}
}

func (self *Fifo) HandleRequestMsg(msg *jsonrpc.RPCMessage) (*jsonrpc.RPCMessage, error) {
	//
	log.Printf("handle request msg %s %d", msg.Method, msg.Id)
	if msg.Method == "greeting" {
		resmsg := jsonrpc.NewResultMessage(msg.Id, "echo")
		return resmsg, nil
	} else if msg.Method == "fifo.put" {
		for _, elem := range msg.Params.MustArray() {
			self.elements = append(self.elements, elem)
		}
		resmsg := jsonrpc.NewResultMessage(msg.Id, "ok")
		return resmsg, nil
	} else if msg.Method == "fifo.get" {
		if len(self.elements) > 0 {
			elem := self.elements[0]
			self.elements = self.elements[1:len(self.elements)]
			resmsg := jsonrpc.NewResultMessage(msg.Id, elem)
			return resmsg, nil
		} else {
			resmsg := jsonrpc.NewResultMessage(msg.Id, nil)
			return resmsg, nil
		}
	} else if msg.Method == "fifo.list" {
		resmsg := jsonrpc.NewResultMessage(msg.Id, self.elements)
		return resmsg, nil
	} else {
		errmsg := jsonrpc.NewErrorMessage(msg.Id, 404, "no such message", false)
		return errmsg, nil
	}
	return nil, nil
}

func ExampleFIFO(c intf.JSONRPCTubeClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	stream, err := c.Handle(ctx)
	if err != nil {
		return err
	}

	fifo := &Fifo{elements: make([]interface{}, 0)}

	// register methods first
	methods := []string{"greeting", "fifo.put", "fifo.get", "fifo.list"}
	reg := &intf.RegisterMethodsRequest{Methods: methods}
	payload := &intf.JSONRPCUpPacket_RegisterMethods{RegisterMethods: reg}
	up_pac := &intf.JSONRPCUpPacket{Payload: payload}
	stream.Send(up_pac)
	
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
			resmsg, err := fifo.HandleRequestMsg(msg)
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
} // end of Example FIFO
