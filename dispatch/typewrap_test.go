package dispatch

import (
	//"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"reflect"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

type TTT struct {
	A string `json:"a"`
	B int `json:"b"`
}

func TestInterfaceAndValue(t *testing.T) {
	assert := assert.New(t)
	tp := reflect.TypeOf(TTT{})

	m := map[string](interface{}){
		"a" : "hello",
		"b" : 5,
	}
	val, err := InterfaceToValue(tp, m)
	assert.Nil(err)

	st, ok := val.Interface().(TTT)
	assert.True(ok)
	assert.Equal("hello", st.A)
	assert.Equal(5, st.B)

	v1, err := ValueToInterface(tp, val)
	assert.Nil(err)
	m1, ok := v1.(map[string]interface{})
	assert.Equal("hello", m1["a"])
	assert.Equal(5, m1["b"])
}

func TestPtrInterfaceAndValue(t *testing.T) {
	assert := assert.New(t)
	tp := reflect.TypeOf(&TTT{})

	m := map[string](interface{}){
		"a" : "hello",
		"b" : 5,
	}
	val, err := InterfaceToValue(tp, m)
	assert.Nil(err)
	st, ok := val.Interface().(*TTT)
	assert.True(ok)
	assert.Equal("hello", st.A)
	assert.Equal(5, st.B)

	v1, err := ValueToInterface(tp, val)
	assert.Nil(err)
	m1, ok := v1.(map[string]interface{})
	assert.True(ok)	
	assert.Equal("hello", m1["a"])
	assert.Equal(5, m1["b"])
}
