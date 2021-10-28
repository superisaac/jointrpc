package neighbor

import (
	//"fmt"
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/superisaac/jointrpc/client"
	"github.com/superisaac/jointrpc/datadir"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/rpcrouter"
	"github.com/superisaac/jointrpc/server"
	"github.com/superisaac/jointrpc/service"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const addSchema = `
{
  "type": "method",
  "params": [
    {
      "type": "number"
    },
    {
      "type": "number"
    }
  ],
  "returns": {
    "type": "number"
  }
}
`

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestNeighborRun(t *testing.T) {
	assert := assert.New(t)

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start server1
	go server.StartGRPCServer(rootCtx, "localhost:10010")

	// start server2

	rootCtx1 := server.ServerContext(rootCtx, nil)
	factory := rpcrouter.RouterFactoryFromContext(rootCtx1)
	factory.Config.Neighbors = make(map[string]datadir.NeighborConfig)

	factory.Config.Neighbors["default"] = datadir.NeighborConfig{
		Peers: []datadir.PeerConfig{{"h2c://localhost:10010", ""}},
	}

	go server.StartGRPCServer(rootCtx1, "localhost:10011")
	time.Sleep(100 * time.Millisecond)

	srv := NewNeighborService()
	service.TryStartService(rootCtx1, srv)
	time.Sleep(100 * time.Millisecond)

	// start server3
	go server.StartGRPCServer(rootCtx, "localhost:10012")

	// start client1, the serve of add2int()
	c1 := client.NewRPCClient(client.ServerEntry{"h2c://localhost:10010", ""})
	disp1 := dispatch.NewDispatcher()
	disp1.On("add2int", func(req *dispatch.RPCRequest, params []interface{}) (interface{}, error) {
		a := jsonrpc.MustInt(params[0], "params[0]")
		b := jsonrpc.MustInt(params[1], "params[1]")
		return a + b, nil
	}, dispatch.WithSchema(addSchema))

	err := c1.Connect()
	assert.Nil(err)
	cCtx, cancelClient := context.WithCancel(context.Background())
	//defer cancelClient()
	go c1.Live(cCtx, disp1)

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
	errbody := res3.MustError()
	assert.Equal(jsonrpc.ErrMethodNotFound.Code, errbody.Code)
	assert.Equal("method not found", errbody.Message)

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
	errbody4 := res4.MustError()
	assert.Equal(jsonrpc.ErrMethodNotFound.Code, errbody4.Code)
	assert.Equal("method not found", errbody4.Message)
}
