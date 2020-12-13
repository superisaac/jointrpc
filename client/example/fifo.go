package example

import (
	//intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
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
	}, client.WithSchema(``))

	rpcClient.On("fifo.list",
		func(req *client.RPCRequest, params []interface{}) (interface{}, error) {
			return fifo.Items, nil
		})

	rpcClient.On("fifo.add",
		func(req *client.RPCRequest, params []interface{}) (interface{}, error) {
			if len(fifo.Items) < 2 {
				return nil, &jsonrpc.RPCError{10408, "items size < 2", false}
			}
			a, err := jsonrpc.ValidateFloat64(fifo.Items[len(fifo.Items) - 2])
			if err != nil {
				return nil, err
			}
			b, err := jsonrpc.ValidateFloat64(fifo.Items[len(fifo.Items) - 1])
			if err != nil {
				return nil, err
			}
			fifo.Items = fifo.Items[0:len(fifo.Items) - 2]
			return a + b, nil

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
