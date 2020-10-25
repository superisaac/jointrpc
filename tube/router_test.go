package tube

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

// implements ConnT
type TestConnT struct {
	ch MsgChannel
}

func (self TestConnT) RecvChannel() MsgChannel {
	return self.ch
}

func (self TestConnT) CanBroadcast() bool {
	return true
}

func TestJoinConn(t *testing.T) {
	assert := assert.New(t)
	router := NewRouter()

	cid := CID(1002)
	ch := make(MsgChannel, 100)
	router.Join(cid, ch)
	assert.Equal(1, len(router.ConnMap))

	router.RegisterMethod(cid, "abc")
	methods := router.GetMethods(cid)
	assert.Equal(1, len(methods))
	assert.Equal("abc", methods[0])
}
