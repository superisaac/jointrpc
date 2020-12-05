package client

import (
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

type Fifo struct {
	Items []interface{}
}

func ExampleFIFO(c intf.JSONRPCTubeClient) error {
	fifo := &Fifo{Items: make([]interface{}, 0)}

	client := NewRPCClient()

	client.Handle("fifo.put", func(msg *jsonrpc.RPCMessage) (*jsonrpc.RPCMessage, error) {
		for _, elem := range msg.Params.MustArray() {
			fifo.Items = append(fifo.Items, elem)
		}
		resmsg := jsonrpc.NewResultMessage(msg.Id, "ok")
		return resmsg, nil
	})

	client.Handle("fifo.get", func(msg *jsonrpc.RPCMessage) (*jsonrpc.RPCMessage, error) {
		if len(fifo.Items) > 0 {
			elem := fifo.Items[0]
			fifo.Items = fifo.Items[1:len(fifo.Items)]
			resmsg := jsonrpc.NewResultMessage(msg.Id, elem)
			return resmsg, nil
		} else {
			resmsg := jsonrpc.NewResultMessage(msg.Id, nil)
			return resmsg, nil
		}
	})

	client.Handle("fifo.list", func(msg *jsonrpc.RPCMessage) (*jsonrpc.RPCMessage, error) {
		resmsg := jsonrpc.NewResultMessage(msg.Id, fifo.Items)
		return resmsg, nil
	})

	client.HandleMethods(c)

	return nil
} // end of Example FIFO
