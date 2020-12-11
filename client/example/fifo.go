package example

import (
	//intf "github.com/superisaac/rpctube/intf/tube"
	//jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	client "github.com/superisaac/rpctube/client"
)

type Fifo struct {
	Items []interface{}
}

func ExampleFIFO(serverAddress string) error {
	fifo := &Fifo{Items: make([]interface{}, 0)}

	rpcClient := client.NewRPCClient(serverAddress)

	rpcClient.On("fifo.put", func(req *client.RPCRequest, params []interface{}) (interface{}, error) {
		for _, elem := range params {
			fifo.Items = append(fifo.Items, elem)
		}
		return "ok", nil
	})

	rpcClient.On("fifo.get", func(req *client.RPCRequest, params []interface{}) (interface{}, error) {
		if len(fifo.Items) > 0 {
			elem := fifo.Items[0]
			fifo.Items = fifo.Items[1:len(fifo.Items)]
			return elem, nil
		} else {
			return nil, nil
		}
	})

	rpcClient.On("fifo.list", func(req *client.RPCRequest, params []interface{}) (interface{}, error) {
		return fifo.Items, nil
	})

	rpcClient.OnDefault(func(req *client.RPCRequest, method string, params []interface{}) (interface{}, error) {
		return "I don't know how to respond", nil
	})

	err := rpcClient.Connect()
	if err != nil {
		return err
	}
	rpcClient.HandleRPC()

	return nil
} // end of Example FIFO
