package example

import (
	client "github.com/superisaac/rpctube/client"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

func ExampleArray(serverAddress string) error {
	items := make([]interface{}, 0)

	rpcClient := client.NewRPCClient(serverAddress)

	rpcClient.On("array.push", func(req *client.RPCRequest, params []interface{}) (interface{}, error) {
		for _, elem := range params {
			items = append(items, elem)
		}
		return "ok", nil
	})

	rpcClient.On("array.at", func(req *client.RPCRequest, params []interface{}) (interface{}, error) {
		if len(params) != 1 {
			return nil, &jsonrpc.RPCError{400, "params count not eq 1", false}
		}

		n, err := jsonrpc.ValidateInt(params[0])
		if err != nil {
			return nil, err
		}

		if n < 0 || n >= int64(len(items)) {
			return nil, &jsonrpc.RPCError{10423, "index out of range", false}
		}
		return items[n], nil
	}, client.WithSchema(``))

	rpcClient.On("array.pophead", func(req *client.RPCRequest, params []interface{}) (interface{}, error) {
		if len(items) > 0 {
			elem := items[0]
			items = items[1:]
			return elem, nil
		} else {
			return nil, nil
		}
	}, client.WithSchema(``))

	rpcClient.On("array.poptail", func(req *client.RPCRequest, params []interface{}) (interface{}, error) {
		if len(items) > 0 {
			elem := items[len(items)-1]
			items = items[:len(items)-1]
			return elem, nil
		} else {
			return nil, nil
		}
	},
		client.WithSchema(``),
		client.WithHelp("pop the last element from the array"),
		client.WithConcurrent(true))

	rpcClient.On("array.list",
		func(req *client.RPCRequest, params []interface{}) (interface{}, error) {
			return items, nil
		},
		client.WithConcurrent(true),
		client.WithHelp("list the array elements"))

	rpcClient.On("array.add",
		func(req *client.RPCRequest, params []interface{}) (interface{}, error) {
			if len(items) < 2 {
				return nil, &jsonrpc.RPCError{10408, "items size < 2", false}
			}
			a, err := jsonrpc.ValidateFloat(items[len(items)-2])
			if err != nil {
				return nil, err
			}
			b, err := jsonrpc.ValidateFloat(items[len(items)-1])
			if err != nil {
				return nil, err
			}
			items = items[0 : len(items)-2]
			v := a + b
			items = append(items, v)
			return v, nil

		},
		client.WithHelp("pop two integers from the array, add them and push the result to array"),
		client.WithConcurrent(true))

	rpcClient.OnDefault(func(req *client.RPCRequest, method string, params []interface{}) (interface{}, error) {
		return "I don't know how to respond", nil
	})

	err := rpcClient.Connect()
	if err != nil {
		return err
	}
	return rpcClient.RunHandlers()
} // end of ExampleArray
