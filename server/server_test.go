package server

import (
	//"fmt"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	client "github.com/superisaac/jointrpc/client"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	handler "github.com/superisaac/jointrpc/rpcrouter/handler"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	assert := assert.New(t)

	root, cancelRoot := context.WithCancel(context.Background())
	defer cancelRoot()

	c1 := context.WithValue(root, "key1", "value1")
	assert.Equal("value1", c1.Value("key1"))
	c2, cancelC2 := context.WithCancel(c1)
	defer cancelC2()
	assert.Equal("value1", c2.Value("key1"))
}

func TestServerClientRound(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartServer(ctx, "127.0.0.1:10001", nil)

	time.Sleep(100 * time.Millisecond)

	c := client.NewRPCClient(client.ServerEntry{"127.0.0.1:10001", ""})
	err := c.Connect()
	assert.Nil(err)

	res, err := c.CallRPC(ctx, ".echo", [](interface{}){"nice"})
	assert.Nil(err)

	assert.True(res.IsResult())
	m, ok := res.MustResult().(map[string]interface{})
	assert.True(ok)
	assert.Equal("nice", m["echo"])

	res1, err := c.CallRPC(ctx, ".echo", [](interface{}){1})
	assert.Nil(err)
	assert.True(res1.IsError())
	errbody, ok := res1.MustError().(map[string]interface{})
	assert.True(ok)
	assert.Equal("Validation Error: .params[0] data is not string", errbody["reason"])
}

func StartTestServe(rootCtx context.Context, serverAddress string) {
	c := client.NewRPCClient(client.ServerEntry{serverAddress, ""})
	c.On("add2int", func(req *handler.RPCRequest, params []interface{}) (interface{}, error) {
		a := jsonrpc.MustInt(params[0], "params[0]")
		b := jsonrpc.MustInt(params[1], "params[1]")
		return a + b, nil
	}, handler.WithSchema(`{"type": "method", "params": [{"type": "number"}, {"type": "number"}], "returns": {"type": "number"}}`))
	c.Connect()
	c.Handle(rootCtx)
}

func TestClientAsServe(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	go StartServer(ctx, "127.0.0.1:10002", nil)
	time.Sleep(100 * time.Millisecond)

	go StartTestServe(ctx, "127.0.0.1:10002")
	time.Sleep(100 * time.Millisecond)

	c := client.NewRPCClient(client.ServerEntry{"127.0.0.1:10002", ""})
	err := c.Connect()
	assert.Nil(err)

	res, err := c.CallRPC(ctx, "add2int", [](interface{}){5, 6})
	assert.Nil(err)
	assert.True(res.IsResult())
	assert.Equal(json.Number("11"), res.MustResult())

}
