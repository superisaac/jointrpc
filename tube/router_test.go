package tube

import (
	"github.com/stretchr/testify/assert"
	"testing"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

// implements ConnT
type TestConnT struct {
	ch MsgChannel
}

func (self TestConnT) RecvChannel() MsgChannel{
	return self.ch
}

func (self TestConnT) CanBroadcast() bool {
	return true
}

func TestJoinConn(t *testing.T) {
	assert := assert.New(t)
	router := NewRouter()

	cid := jsonrpc.CID(1002)
	ch := make(MsgChannel, 100)
	c := TestConnT{ch: ch}
	router.JoinConn(cid, c)
	assert.Equal(len(router.ConnMap), 1)
}
