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

	msg, err := Parse(js)
	assert.Nil(err)

	assert.True(msg.IsRequest())
	assert.False(msg.IsNotify())
	assert.False(msg.IsError())
	assert.False(msg.IsResult())

	assert.Equal(json.Number("100"), msg.MustId())
	assert.Equal("abc::add", msg.MustMethod())
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

	msg, err := Parse(js)
	assert.Nil(err)

	assert.Equal("abc::add", msg.MustMethod())

	assert.False(msg.IsRequest())
	assert.True(msg.IsNotify())
	assert.False(msg.IsError())
	assert.False(msg.IsResult())

	params := msg.MustParams()
	assert.Equal(len(params), 3)
	assert.Equal(params[0], json.Number("13"))
	assert.Equal(params[1], json.Number("4"))
	assert.Equal(params[2], "hello")

	arr := [](interface{}){3, "uu"}
	msg = NewNotifyMessage("hahaha", arr, nil)
	assert.Equal("hahaha", msg.MustMethod())
	params = msg.MustParams()

	assert.Equal(len(params), 2)
	assert.Equal(params[1], "uu")
}

func TestGuessJson(t *testing.T) {
	assert := assert.New(t)

	v1, err := GuessJson("5")
	assert.Equal(int64(5), v1)

	v1_1, err := GuessJson("-5")
	assert.Equal(int64(-5), v1_1)

	v1_2, err := GuessJson("-5.78389383")
	assert.InDelta(float64(-5.78389383), v1_2, 0.0001)

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
