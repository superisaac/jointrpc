package server

import (
	//"fmt"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	client "github.com/superisaac/jointrpc/client"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	//mirror "github.com/superisaac/jointrpc/cluster/mirror"
	datadir "github.com/superisaac/jointrpc/datadir"
	handler "github.com/superisaac/jointrpc/rpcrouter/handler"
	"testing"
	"time"
)

func TestMirrorRun(t *testing.T) {
	assert := assert.New(t)

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// start server1
	go StartServer(rootCtx, "localhost:10010", nil)

	// start server2
	cfg := datadir.NewConfig()
	cfg.Cluster.StaticPeers = []datadir.PeerConfig{{"h2c://localhost:10010", ""}}
	go StartServer(rootCtx, "localhost:10011", cfg)
	time.Sleep(100 * time.Millisecond)

	// start server3
	go StartServer(rootCtx, "localhost:10012", nil)

	// start client1, the serve of add2int()
	c1 := client.NewRPCClient(client.ServerEntry{"h2c://localhost:10010", ""})
	c1.On("add2int", func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
		a := jsonrpc.MustInt(params[0], "params[0]")
		b := jsonrpc.MustInt(params[1], "params[1]")
		return a + b, nil
	}, handler.WithSchema(addSchema))
	err := c1.Connect()
	assert.Nil(err)
	cCtx, cancelClient := context.WithCancel(context.Background())
	//defer cancelClient()
	go c1.Handle(cCtx)

	// start c2, the add2int() caller to server2
	time.Sleep(100 * time.Millisecond)
	c2 := client.NewRPCClient(client.ServerEntry{"h2c://localhost:10011", ""})
	err = c2.Connect()
	assert.Nil(err)

	// call rpc from server2 which delegates server1
	time.Sleep(100 * time.Millisecond)
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	delegates, err := c2.ListDelegates(ctx1)
	assert.Nil(err)
	assert.Equal([]string{"add2int"}, delegates)
	res, err := c2.CallRPC(ctx1, "add2int", [](interface{}){5, 6}, client.WithTraceId("trace1"))
	assert.Nil(err)
	assert.Equal("trace1", res.TraceId())
	assert.True(res.IsResult())
	assert.Equal(json.Number("11"), res.MustResult())

	// start client3
	c3 := client.NewRPCClient(client.ServerEntry{"h2c://localhost:10012", ""})
	err = c3.Connect()
	assert.Nil(err)

	// call rpc from server3 which doesnot delegates server1
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	res3, err := c3.CallRPC(ctx2, "add2int", [](interface{}){15, 16}, client.WithTraceId("test21"))
	assert.Nil(err)
	assert.Equal("test21", res3.TraceId())
	assert.True(res3.IsError())
	errBody, ok := res3.MustError().(map[string]interface{})
	assert.True(ok)
	assert.Equal(json.Number("404"), errBody["code"])
	assert.Equal("method not found", errBody["reason"])

	// close client serve, the add2int() provider
	cancelClient()
	time.Sleep(100 * time.Millisecond)
	d1, err := c2.ListDelegates(ctx1)
	assert.Nil(err)
	assert.Equal(0, len(d1))
	//assert.Equal([]string{}, d1)

	// call rpc from server2 which doesnot delegates server1
	ctx4, cancel4 := context.WithCancel(context.Background())
	defer cancel4()
	res4, err := c2.CallRPC(ctx4, "add2int", [](interface{}){15, 16}, client.WithTraceId("trace13"))
	assert.Nil(err)
	assert.Equal("trace13", res4.TraceId())
	assert.True(res4.IsError())
	errBody4, ok := res4.MustError().(map[string]interface{})
	assert.True(ok)
	assert.Equal(json.Number("404"), errBody4["code"])
	assert.Equal("method not found", errBody4["reason"])
}
