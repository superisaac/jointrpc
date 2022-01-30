package msgutil

import (
	"io/ioutil"
	"os"
	//"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/superisaac/jsonz"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestMsgEncoding(t *testing.T) {
	assert := assert.New(t)

	req := jsonz.NewRequestMessage(1, "test", []interface{}{})
	req.SetTraceId("hello1")
	assert.Equal("hello1", req.TraceId())

	envo := MessageToEnvolope(req)
	copied, err := MessageFromEnvolope(envo)
	assert.Nil(err)

	assert.Equal("hello1", copied.TraceId())
}
