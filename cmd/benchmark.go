package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
	//"strings"

	//intf "github.com/superisaac/jointrpc/intf/jointrpc"
	client "github.com/superisaac/jointrpc/client"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	//example "github.com/superisaac/jointrpc/client/example"
	//grpc "google.golang.org/grpc"
)

// Call benchmark
func CommandCallBenchmark() {
	callFlags := flag.NewFlagSet("benchmark", flag.ExitOnError)
	serverFlag := client.NewServerFlag(callFlags)
	pBroadcast := callFlags.Bool("broadcast", false, "broadcast the notify to all listeners")
	pTraceId := callFlags.String("traceid", "", "trace id during the workflow")

	pConcurrency := callFlags.Int("con", 10, "the number of concurrent clients")
	pNum := callFlags.Int("num", 10, "the number of calls each client can call")

	callFlags.Parse(os.Args[2:])
	// TODO, check the sanity agains traceId

	if callFlags.NArg() < 1 {
		fmt.Println("method params...")
		os.Exit(1)
	}

	args := callFlags.Args()
	method := args[0]
	clParams := args[1:len(args)]

	params, err := jsonrpc.GuessJsonArray(clParams)
	if err != nil {
		panic(err)
	}

	RunCallBenchmark(
		serverFlag.Get(),
		method, params,
		*pConcurrency, *pNum,
		client.WithBroadcast(*pBroadcast), client.WithTraceId(*pTraceId))
}

func toS(ns int) float64 {
	return float64(ns) / float64(time.Second)
}

func RunCallBenchmark(serverEntry client.ServerEntry, method string, params []interface{}, concurrency int, num int, opts ...client.CallOptionFunc) {
	chResults := make(chan int, concurrency*num)
	results := make([]int, concurrency*num)
	var sum int = 0

	for a := 0; a < concurrency; a++ {
		go callNTimes(chResults, serverEntry, method, params, num, opts...)
	}

	for i := 0; i < concurrency*num; i++ {
		usedTime := <-chResults
		sum += usedTime
		results[i] = usedTime
	}

	sort.Ints(results)

	avg := sum / len(results)
	pos95 := int(0.95 * float64(len(results)))
	t95 := results[pos95]
	maxv := results[len(results)-1]
	minv := results[0]
	fmt.Printf("avg=%g, min=%g, p95=%g, max=%g\n", toS(avg), toS(minv), toS(t95), toS(maxv))
}

func callNTimes(chResults chan int, serverEntry client.ServerEntry, method string, params []interface{}, num int, opts ...client.CallOptionFunc) error {
	ctx := context.Background()
	c := client.NewRPCClient(serverEntry)
	err := c.Connect()
	if err != nil {
		return err
	}

	for i := 0; i < num; i++ {
		startTime := time.Now()
		_, err := c.CallRPC(ctx, method, params, opts...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad results %+v\n", err)
		}
		endTime := time.Now()
		chResults <- int(endTime.Sub(startTime))
	}
	return nil
}
