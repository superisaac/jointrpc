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
	"github.com/superisaac/jsonz"
	//example "github.com/superisaac/jointrpc/client/example"
	//grpc "google.golang.org/grpc"
)

// Call benchmark
func CommandCallBenchmark() {
	callFlags := flag.NewFlagSet("benchmark", flag.ExitOnError)
	serverFlag := client.NewServerFlag(callFlags)
	pTraceId := callFlags.String("traceid", "", "trace id during the workflow")

	pConcurrency := callFlags.Uint("con", 10, "the number of concurrent clients")
	pNum := callFlags.Uint("n", 10, "the number of calls each client can call")

	callFlags.Parse(os.Args[2:])
	// TODO, check the sanity agains traceId

	if callFlags.NArg() < 1 {
		fmt.Println("method params...")
		os.Exit(1)
	}

	args := callFlags.Args()
	method := args[0]
	clParams := args[1:len(args)]

	params, err := jsonz.GuessJsonArray(clParams)
	if err != nil {
		panic(err)
	}

	RunCallBenchmark(
		serverFlag.Get(),
		method, params,
		*pConcurrency, *pNum,
		client.WithTraceId(*pTraceId))
}

func toS(ns uint) float64 {
	return float64(ns) / float64(time.Second)
}

func RunCallBenchmark(serverEntry client.ServerEntry, method string, params []interface{}, concurrency uint, num uint, opts ...client.CallOptionFunc) {
	chResults := make(chan uint, concurrency*num)
	results := make([]uint, concurrency*num)
	var sum uint = 0

	for a := uint(0); a < concurrency; a++ {
		go callNTimes(chResults, serverEntry, method, params, num, opts...)
	}

	for i := uint(0); i < concurrency*num; i++ {
		usedTime := <-chResults
		sum += usedTime
		results[i] = usedTime
	}

	//sort.Uints(results)
	sort.Slice(results, func(i, j int) bool { return results[i] < results[j] })

	avg := sum / uint(len(results))
	pos95 := int(0.95 * float64(len(results)))
	t95 := results[pos95]
	maxv := results[len(results)-1]
	minv := results[0]
	fmt.Printf("avg=%gs, min=%gs, p95=%gs, max=%gs\n", toS(avg), toS(minv), toS(t95), toS(maxv))
}

func callNTimes(chResults chan uint, serverEntry client.ServerEntry, method string, params []interface{}, num uint, opts ...client.CallOptionFunc) error {
	ctx := context.Background()
	c := client.NewRPCClient(serverEntry)
	err := c.Connect()
	if err != nil {
		return err
	}

	for i := uint(0); i < num; i++ {
		startTime := time.Now()
		_, err := c.CallRPC(ctx, method, params, opts...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad results %d %s\n", i, err)
		}
		endTime := time.Now()
		chResults <- uint(endTime.Sub(startTime))
		//time.Sleep(10 * time.Millisecond)
	}
	return nil
}
