package client

import (
	"context"
	"github.com/stretchr/testify/assert"
	server "github.com/superisaac/jointrpc/server"
	"testing"
	"time"
)

func TestServerClientRound(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.StartServer(ctx, "127.0.0.1:10001")

	time.Sleep(1 * time.Second)

	c := NewRPCClient(ServerEntry{"127.0.0.1:10001", ""})

	err := c.Connect()
	assert.Nil(err)

	res, err := client.CallRPC(ctx, ".listMethods", [](interface{}){})
	assert.Nil(res)
	fmt.Printf("res %+v", res)
}
