package client

import (
	intf "github.com/superisaac/rpctube/intf/tube"
	//jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

type Fifo struct {
	Items []interface{}
}

func ExampleFIFO(c intf.JSONRPCTubeClient) error {
	fifo := &Fifo{Items: make([]interface{}, 0)}

	client := NewRPCClient()

	client.Handle("fifo.put", func(req *RPCRequest, params []interface{}) (interface{}, error) {
		for _, elem := range params {
			fifo.Items = append(fifo.Items, elem)
		}
		return "ok", nil
	})

	client.Handle("fifo.get", func(req *RPCRequest, params []interface{}) (interface{}, error) {
		if len(fifo.Items) > 0 {
			elem := fifo.Items[0]
			fifo.Items = fifo.Items[1:len(fifo.Items)]
			return elem, nil
		} else {
			return nil, nil
		}
	})

	client.Handle("fifo.list", func(req *RPCRequest, params []interface{}) (interface{}, error) {
		return fifo.Items, nil
	})

	client.HandleMethods(c)

	return nil
} // end of Example FIFO
