package builtin

import (
	//"fmt"
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	client "github.com/superisaac/jointrpc/client"
	server "github.com/superisaac/jointrpc/server"
	service "github.com/superisaac/jointrpc/service"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestBuiltinMethods(t *testing.T) {
	assert := assert.New(t)

	ctx0, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx := server.ServerContext(ctx0, nil)
	go server.StartGRPCServer(ctx, "127.0.0.1:10001")

	time.Sleep(100 * time.Millisecond)
	srv := NewBuiltinService()

	service.TryStartService(ctx, srv)

	time.Sleep(1000 * time.Millisecond)

	c := client.NewRPCClient(client.ServerEntry{"h2c://127.0.0.1:10001", ""})
	err := c.Connect()
	assert.Nil(err)
	res, err := c.CallRPC(ctx, "_echo", [](interface{}){"nice"}, client.WithTraceId("trace1"))
	assert.Nil(err)
	assert.Equal("trace1", res.TraceId())

	//fmt.Printf("sssres %+v\n", res)
	assert.True(res.IsResult())
	m, ok := res.MustResult().(map[string]interface{})
	assert.True(ok)
	assert.Equal("nice", m["echo"])

	res1, err := c.CallRPC(ctx, "_echo", [](interface{}){1}, client.WithTraceId("trace2"))
	assert.Nil(err)
	assert.Equal("trace2", res1.TraceId())
	assert.True(res1.IsError())
	errbody := res1.MustError()
	assert.Equal("Validation Error: .params[0] data is not string", errbody.Message)
}
