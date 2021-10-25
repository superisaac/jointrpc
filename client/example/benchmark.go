package example

import (
	"context"
	//"fmt"
	client "github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/dispatch"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
)

const (
	benchmarkEchoSchema = `
{
  "type": "method",
  "params": ["string"],
  "returns": "string"
}

`
)

func ExampleBenchmark(serverEntry client.ServerEntry) error {
	disp := dispatch.NewDispatcher()

	rpcClient := client.NewRPCClient(serverEntry)

	// hooked methods
	disp.On("benchmark.echo", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		src := jsonrpc.ConvertString(params[0])
		return src, nil
	}, dispatch.WithSchema(benchmarkEchoSchema))

	err := rpcClient.Connect()
	if err != nil {
		return err
	}
	return rpcClient.Worker(context.Background(), disp)
} // end of ExampleBenchmark
