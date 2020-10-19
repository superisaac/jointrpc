package jsonrpc

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestMsg(t *testing.T) {
	assert := assert.New(t)

	j1 := `{
  "id": 100,
  "method": "abc::add",
  "params": [3, 4, 5]
  }`

	js, _ := simplejson.NewJson([]byte(j1))

	assert.Equal(js.Get("id").MustInt(), 100)
	assert.Equal(js.Get("id").MustString(), "")

	msg := NewRPCMessage(js)

	intMsgId, err := msg.GetIntId()
	assert.Nil(err)
	assert.Equal(intMsgId, int64(100))

	assert.Equal(msg.Method, "add")

	assert.Equal(msg.ServiceName, "abc")
	assert.True(msg.IsValid())
	assert.True(msg.IsRequest())
	assert.False(msg.IsNotify())
	assert.False(msg.IsError())
	assert.False(msg.IsResult())
}

func TestNotifyMsg(t *testing.T) {
	assert := assert.New(t)

	j1 := `{
  "method": "abc::add",
  "params": [13, 4, "hello"]
  }`

	js, _ := simplejson.NewJson([]byte(j1))

	assert.Equal(js.Get("id").MustInt(), 0)
	assert.Equal(js.Get("id").MustString(), "")

	msg := NewRPCMessage(js)
	intMsgId, err := msg.GetIntId()
	assert.NotNil(err)
	assert.Equal(intMsgId, int64(0))

	assert.Equal(msg.Method, "add")

	assert.Equal(msg.ServiceName, "abc")
	assert.True(msg.IsValid())
	assert.False(msg.IsRequest())
	assert.True(msg.IsNotify())
	assert.False(msg.IsError())
	assert.False(msg.IsResult())

	params := msg.GetParams()
	assert.Equal(len(params), 3)
	assert.Equal(params[0], json.Number("13"))
	assert.Equal(params[1], json.Number("4"))
	assert.Equal(params[2], "hello")

	arr := [](interface{}){3, "uu"}
	msg = NewNotifyMessage("ttt", "hahaha", arr)
	assert.Equal(msg.ServiceName, "ttt")
	assert.Equal(msg.Method, "hahaha")
	params = msg.GetParams()

	assert.Equal(len(params), 2)
	assert.Equal(params[1], "uu")
}
