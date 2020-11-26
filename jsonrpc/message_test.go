package jsonrpc

import (
	//"fmt"
	json "encoding/json"
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
	assert.Equal(intMsgId, UID(100))

	assert.Equal(msg.Method, "abc::add")

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
	assert.Equal(intMsgId, UID(0))

	assert.Equal(msg.Method, "abc::add")

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
	msg = NewNotifyMessage("hahaha", arr)
	assert.Equal(msg.Method, "hahaha")
	params = msg.GetParams()

	assert.Equal(len(params), 2)
	assert.Equal(params[1], "uu")
}

func TestGuessJson(t * testing.T) {
	assert := assert.New(t)

	v1, err := GuessJson("5")
	assert.Equal(int64(5), v1)

	v2, err := GuessJson("false")
	assert.Equal(false, v2)

	_, err = GuessJson("[aaa")
	assert.Contains(err.Error(), "invalid character")

	v3, err := GuessJson(`{"abc": 5}`)
	map3 := v3.(map[string]interface{})
	assert.NotNil(map3)
	assert.Equal(json.Number("5"), map3["abc"])

	v4, err := GuessJsonArray([]string{"5", "hahah", `{"ccc": 6}`})
	assert.Equal(3, len(v4))
	assert.Equal(int64(5), v4[0])
	assert.Equal("hahah", v4[1])
}
