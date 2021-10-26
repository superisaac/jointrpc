package example

import (
	"context"
	//"fmt"
	client "github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/dispatch"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"sync"
)

const (
	benchmarkEchoSchema = `
{
  "type": "method",
  "params": ["string"],
  "returns": "string"
}`
)

func ExampleBenchmark(serverEntry client.ServerEntry, concurrency int) error {
	disp := dispatch.NewDispatcher()
	disp.SetSpawnExec(true)
	// hooked methods
	disp.On("benchmark.echo", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		src := jsonrpc.ConvertString(params[0])
		return src, nil
	}, dispatch.WithSchema(benchmarkEchoSchema))

	wg := new(sync.WaitGroup)
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go runOneClient(wg, serverEntry, disp)
	}
	wg.Wait()
	return nil
}

func runOneClient(wg *sync.WaitGroup, serverEntry client.ServerEntry, disp *dispatch.Dispatcher) {
	rpcClient := client.NewRPCClient(serverEntry)
	err := rpcClient.Connect()
	if err != nil {
		panic(err)
	}
	err = rpcClient.Worker(context.Background(), disp)
	if err != nil {
		panic(err)
	}
	wg.Done()
} // end of ExampleBenchmark
