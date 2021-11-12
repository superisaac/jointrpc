package example

import (
	"context"
	"fmt"
	client "github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/dispatch"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
)

func ExampleArray(serverEntry client.ServerEntry) error {
	items := make([]interface{}, 0)

	disp := dispatch.NewDispatcher()

	rpcClient := client.NewRPCClient(serverEntry)

	// hooked methods
	disp.On("array.push",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			for _, elem := range params {
				items = append(items, elem)
			}
			return "ok", nil
		})

	disp.OnTyped("array.at",
		func(req *dispatch.RPCRequest, n int) (interface{}, error) {
			if n < 0 || n >= len(items) {
				return nil, &jsonrpc.RPCError{10423, "parameter 1 index out of range", false}
			}
			return items[n], nil
		},
		dispatch.WithSchema(`{"type": "method", "params": [{"type": "integer"}]}`),
		dispatch.WithHelp("return the element at index"))

	disp.OnTyped("array.size",
		func(req *dispatch.RPCRequest) (int, error) {
			return len(items), nil
		}, dispatch.WithHelp("get the size of array"))

	disp.OnTyped("array.pophead",
		func(req *dispatch.RPCRequest) (interface{}, error) {
			if len(items) > 0 {
				elem := items[0]
				items = items[1:]
				return elem, nil
			} else {
				return nil, nil
			}
		}, dispatch.WithSchema(``))

	disp.OnTyped("array.poptail",
		func(req *dispatch.RPCRequest) (interface{}, error) {
			if len(items) > 0 {
				elem := items[len(items)-1]
				items = items[:len(items)-1]
				return elem, nil
			} else {
				return nil, nil
			}
		},
		dispatch.WithSchema(``),
		dispatch.WithHelp("pop the last element from the array"))

	disp.OnTyped("array.list",
		func(req *dispatch.RPCRequest) ([]interface{}, error) {
			return items, nil
		},
		dispatch.WithHelp("list the array elements"))

	disp.OnTyped("array.add",
		func(req *dispatch.RPCRequest) (float64, error) {
			if len(items) < 2 {
				return 0, &jsonrpc.RPCError{Code: 10408, Message: "array size < 2"}
			}
			a, err := jsonrpc.ValidateFloat(items[len(items)-2], "array[-2]")
			if err != nil {
				return 0, err
			}
			b, err := jsonrpc.ValidateFloat(items[len(items)-1], "array[-1]")
			if err != nil {
				return 0, err
			}
			items = items[0 : len(items)-2]
			v := a + b
			items = append(items, v)
			return v, nil

		},
		dispatch.WithHelp("pop two integers from the array, add them and push the result to array"))

	disp.On("metrics.collect",
		func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {

			lines := []string{
				"# TYPE array_size gauge",
				"# HELP array_size array size",
				fmt.Sprintf(`array_size{collect="example.array"} %d`, len(items)),
			}
			return lines, nil
		})

	disp.OnDefault(func(req *dispatch.RPCRequest, method string, params []interface{}) (interface{}, error) {
		return "I don't know how to respond", nil
	})

	err := rpcClient.Connect()
	if err != nil {
		return err
	}
	return rpcClient.Live(context.Background(), disp)
} // end of ExampleArray
