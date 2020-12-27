package example

import (
	"context"
	client "github.com/superisaac/rpctube/client"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	handler "github.com/superisaac/rpctube/tube/handler"
)

func ExampleArray(serverAddress string, certFile string) error {
	items := make([]interface{}, 0)

	rpcClient := client.NewRPCClient(serverAddress, certFile)

	rpcClient.On("array.push", func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
		for _, elem := range params {
			items = append(items, elem)
		}
		return "ok", nil
	})

	rpcClient.On("array.at", func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
		if len(params) != 1 {
			return nil, &jsonrpc.RPCError{400, "params count not eq 1", false}
		}

		n, err := jsonrpc.ValidateInt(params[0], "parameter 1")
		if err != nil {
			//return nil, err
			panic(err)
		}

		if n < 0 || n >= int64(len(items)) {
			return nil, &jsonrpc.RPCError{10423, "parameter 1 index out of range", false}
		}
		return items[n], nil
	}, handler.WithSchema(`{"type": "method", "params": [{"type": "number"}]}`))

	rpcClient.On("array.size", func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
		return len(items), nil
	}, handler.WithHelp("get the size of array"))

	rpcClient.On("array.pophead", func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
		if len(items) > 0 {
			elem := items[0]
			items = items[1:]
			return elem, nil
		} else {
			return nil, nil
		}
	}, handler.WithSchema(``))

	rpcClient.On("array.poptail", func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
		if len(items) > 0 {
			elem := items[len(items)-1]
			items = items[:len(items)-1]
			return elem, nil
		} else {
			return nil, nil
		}
	},
		handler.WithSchema(``),
		handler.WithHelp("pop the last element from the array"),
		handler.WithConcurrent(true))

	rpcClient.On("array.list",
		func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
			return items, nil
		},
		handler.WithConcurrent(true),
		handler.WithHelp("list the array elements"))

	rpcClient.On("array.add",
		func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
			if len(items) < 2 {
				return nil, &jsonrpc.RPCError{10408, "array size < 2", false}
			}
			a, err := jsonrpc.ValidateFloat(items[len(items)-2], "array[-2]")
			if err != nil {
				return nil, err
			}
			b, err := jsonrpc.ValidateFloat(items[len(items)-1], "array[-1]")
			if err != nil {
				return nil, err
			}
			items = items[0 : len(items)-2]
			v := a + b
			items = append(items, v)
			return v, nil

		},
		handler.WithHelp("pop two integers from the array, add them and push the result to array"),
		handler.WithConcurrent(true))

	rpcClient.OnDefault(func(req *handler.RPCRequest, method string, params []interface{}) (interface{}, error) {
		return "I don't know how to respond", nil
	})

	err := rpcClient.Connect()
	if err != nil {
		return err
	}
	return rpcClient.Handle(context.Background())
} // end of ExampleArray
