package playbook

import (
	//"fmt"
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	client "github.com/superisaac/jointrpc/client"
	server "github.com/superisaac/jointrpc/server"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

const PbSay = `---
version: 1.0.0
methods:
  say:
    description: say somthing using echo
    shell:
      command: jq '"echo " + .params[0]'
      env:
        - "AA=BB"
    schema:
      type: 'method'
      params:
        - type: 'string'
      returns:
        type: 'string'
`

func TestPlaybook(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start grpc server
	go server.StartGRPCServer(ctx, "127.0.0.1:11000")
	time.Sleep(100 * time.Millisecond)

	// create playbook instance and run
	pb := NewPlaybook()
	err := pb.Config.ReadConfigBytes([]byte(PbSay))
	assert.Nil(err)
	//	assert.Equal(
	method, ok := pb.Config.Methods["say"]
	assert.True(ok)
	assert.NotNil(method.innerSchema)
	assert.Equal("method", method.innerSchema.Type())

	go func() {
		err := pb.Run(client.ServerEntry{"h2c://127.0.0.1:11000", ""})
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	// create client
	c := client.NewRPCClient(client.ServerEntry{"h2c://127.0.0.1:11000", ""})
	err = c.Connect()
	assert.Nil(err)

	res, err := c.CallRPC(ctx, "say", [](interface{}){"nice"}, client.WithTraceId("testsay01"))
	assert.Nil(err)
	assert.Equal("testsay01", res.TraceId())
	assert.True(res.IsResult())
	assert.Equal("echo nice", res.MustResult())
}
