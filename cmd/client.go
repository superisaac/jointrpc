package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	//"strings"

	simplejson "github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	//intf "github.com/superisaac/jointrpc/intf/jointrpc"
	client "github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/rpcrouter"
	"github.com/superisaac/jsonz"
	"os"
	//example "github.com/superisaac/jointrpc/client/example"
	//grpc "google.golang.org/grpc"
)

// Send Notify
func CommandSendNotify() {
	callFlags := flag.NewFlagSet("notify", flag.ExitOnError)
	serverFlag := client.NewServerFlag(callFlags)
	pBroadcast := callFlags.Bool("broadcast", false, "broadcast the notify to all listeners")
	pTraceId := callFlags.String("traceid", "", "trace id during the workflow")

	callFlags.Parse(os.Args[2:])

	if callFlags.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "method params...\n")
		os.Exit(1)
	}

	args := callFlags.Args()
	method := args[0]
	clParams := args[1:len(args)]

	params, err := jsonz.GuessJsonArray(clParams)
	if err != nil {
		//panic(err)
		log.Errorf("params error %s", err)
		os.Exit(1)
	}

	err = RunSendNotify(serverFlag.Get(), method, params,
		client.WithBroadcast(*pBroadcast), client.WithTraceId(*pTraceId))
	if err != nil {
		//panic(err)
		log.Errorf("%s", err)
		os.Exit(1)
	}
}

func RunSendNotify(serverEntry client.ServerEntry, method string, params []interface{}, opts ...client.CallOptionFunc) error {
	c := client.NewRPCClient(serverEntry)
	err := c.Connect()
	if err != nil {
		return err
	}
	err = c.SendNotify(context.Background(), method, params, opts...)
	if err != nil {
		return err
	}
	return nil
}

// Call RPC
func CommandCallRPC(subcmd string) {
	callFlags := flag.NewFlagSet(subcmd, flag.ExitOnError)
	serverFlag := client.NewServerFlag(callFlags)
	pBroadcast := callFlags.Bool("broadcast", false, "broadcast the notify to all listeners")
	pTraceId := callFlags.String("traceid", "", "trace id during the workflow")

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
		//panic(err)
		log.Errorf("fail to parse json %s", err)
		os.Exit(1)
	}

	err = RunCallRPC(serverFlag.Get(), method, params,
		client.WithBroadcast(*pBroadcast), client.WithTraceId(*pTraceId))
	if err != nil {
		//panic(err)
		log.Errorf("fail to call rpc %s", err)
		os.Exit(1)
	}
}

func RunCallRPC(serverEntry client.ServerEntry, method string, params []interface{}, opts ...client.CallOptionFunc) error {
	c := client.NewRPCClient(serverEntry)
	err := c.Connect()
	if err != nil {
		return err
	}

	res, err := c.CallRPC(context.Background(), method, params, opts...)
	if err != nil {
		return err
	}
	log.Debugf("result got trace id %s", res.TraceId())
	repr, err := jsonz.EncodePretty(res)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", repr)
	//repr = res.MustString()
	//fmt.Printf("%s\n", repr)
	return nil
}

// Call ListMethods
func CommandListMethods() {
	aFlags := flag.NewFlagSet("methods", flag.ExitOnError)
	serverFlag := client.NewServerFlag(aFlags)
	aFlags.Parse(os.Args[2:])

	err := RunListMethods(serverFlag.Get())
	if err != nil {
		panic(err)
	}
}

func RunListMethods(serverEntry client.ServerEntry) error {
	c := client.NewRPCClient(serverEntry)
	err := c.Connect()
	if err != nil {
		return err
	}
	methodInfos, err := c.ListMethods(context.Background())
	if err != nil {
		return nil
	}

	fmt.Printf("available methods:\n")
	for _, minfo := range methodInfos {
		fmt.Printf("  %s\t%s\n", minfo.Name, minfo.Help)
	}
	return nil
}

// Call ListDelegates
func CommandListDelegates() {
	aFlags := flag.NewFlagSet("delegates", flag.ExitOnError)
	serverFlag := client.NewServerFlag(aFlags)
	aFlags.Parse(os.Args[2:])

	err := RunListDelegates(serverFlag.Get())
	if err != nil {
		panic(err)
	}
}

func RunListDelegates(serverEntry client.ServerEntry) error {
	c := client.NewRPCClient(serverEntry)
	err := c.Connect()
	if err != nil {
		return err
	}
	delegates, err := c.ListDelegates(context.Background())
	if err != nil {
		return nil
	}

	fmt.Printf("delegated methods:\n")
	for _, name := range delegates {
		fmt.Printf("  %s\n", name)
	}
	return nil
}

// Watch notify
func CommandWatch() {
	subFlags := flag.NewFlagSet("watchnotify", flag.ExitOnError)
	serverFlag := client.NewServerFlag(subFlags)
	subFlags.Parse(os.Args[2:])

	notifyNames := subFlags.Args()
	if len(notifyNames) < 1 {
		panic(errors.New("No notify methods specified to watch"))
	}

	rpcClient := client.NewRPCClient(serverFlag.Get())

	disp := dispatch.NewDispatcher()

	for _, notifyName := range notifyNames {
		disp.On(notifyName, func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
			msg := req.CmdMsg.Msg
			repr, err := jsonz.EncodePretty(msg)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", repr)
			return nil, nil
		})
	}
	err := rpcClient.Connect()
	if err != nil {
		panic(err)
	}

	err = rpcClient.Live(context.Background(), disp)
	if err != nil {
		panic(err)
	}
}

// Watch methods update
func CommandWatchState() {
	subFlags := flag.NewFlagSet("watchstate", flag.ExitOnError)
	serverFlag := client.NewServerFlag(subFlags)
	pVerbose := subFlags.Bool("verbose", false, "show method info")
	subFlags.Parse(os.Args[2:])

	rpcClient := client.NewRPCClient(serverFlag.Get())
	stateListener := dispatch.NewStateListener()

	if *pVerbose {
		stateListener.OnStateChange(printMethodInfos)
	} else {
		stateListener.OnStateChange(printMethodNames)
	}

	rpcClient.OnAuthorized(func() {
		req := rpcClient.NewWatchStateRequest()
		rpcClient.LiveCall(context.Background(), req,
			func(res jsonz.Message) {
				log.Infof("watch state")
			})

	})

	err := rpcClient.Connect()
	if err != nil {
		panic(err)
	}
	disp := dispatch.NewDispatcher()
	client.OnStateChanged(disp, stateListener)
	//rpcClient.SubscribeState(context.Background(), stateListener)
	rpcClient.Live(context.Background(), disp)
}

func printMethodInfos(state *rpcrouter.ServerState) {
	arr := make([](map[string](interface{})), 0)
	for _, info := range state.Methods {
		mapInfo := map[string](interface{}){
			"name": info.Name,
		}
		if info.Help != "" {
			mapInfo["help"] = info.Help
		}

		if info.Schema() != nil {
			mapInfo["schema"] = info.Schema().RebuildType()
		}
		// TODO: schema
		arr = append(arr, mapInfo)
	}
	jarr := simplejson.New()
	jarr.Set("methods", arr)
	repr, err := jarr.MarshalJSON()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(repr))
}

func printMethodNames(state *rpcrouter.ServerState) {
	arr := make([]string, 0)
	for _, info := range state.Methods {
		arr = append(arr, info.Name)
	}
	jarr := simplejson.New()
	//jarr.SetPath(nil, arr)
	jarr.Set("methods", arr)
	repr, err := jarr.MarshalJSON()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(repr))
}
