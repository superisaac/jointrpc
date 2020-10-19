package jsonrpc

import (
	"encoding/json"
	"errors"
	"github.com/bitly/go-simplejson"
	"strings"
)

func ParseMessage(data []byte) (RPCMessage, error) {
	// Reserved for other data format
	parsed, err := simplejson.NewJson(data)
	if err != nil {
		return RPCMessage{}, err
	}
	return NewRPCMessage(parsed), nil
}

func NewRPCMessage(data *simplejson.Json) RPCMessage {
	//msg := new(RPCMessage)
	msg := RPCMessage{Initialized: true}
	msg.Id = data.Get("id").Interface()
	method, err := data.Get("method").String()
	if err == nil {
		arr := strings.SplitN(method, "::", 2)
		if len(arr) == 2 {
			msg.ServiceName = arr[0]
			msg.Method = arr[1]
		} else if len(arr) == 1 {
			msg.Method = arr[0]
		}
	}
	msg.Params = data.Get("params")
	msg.Result = data.Get("result")
	msg.Error = data.Get("error")
	msg.Raw = data
	return msg
}

func NewResultMessage(id interface{}, result interface{}) RPCMessage {
	resultJson := simplejson.New()
	resultJson.Set("id", id)
	resultJson.Set("result", result)
	return NewRPCMessage(resultJson)
}

func NewNotifyMessage(serviceName string, method string, params []interface{}) RPCMessage {
	notifyJson := simplejson.New()
	notifyJson.Set("method", serviceName+"::"+method)
	notifyJson.Set("params", params)
	return NewRPCMessage(notifyJson)
}

func NewErrorMessage(id interface{}, code int, message string) RPCMessage {
	jsonData := NewErrorJSON(id, code, message)
	return NewRPCMessage(jsonData)
}

func NewErrorJSON(id interface{}, code int, message string) *simplejson.Json {
	errJson := simplejson.New()
	errJson.Set("code", code)
	errJson.Set("message", message)
	body := simplejson.New()
	body.Set("id", id)
	body.Set("error", errJson.Interface())
	return body
}

func (self RPCMessage) GetIntId() (int64, error) {
	msgId, ok := self.Id.(json.Number)
	if !ok {
		return 0, errors.New("not a number")
	}
	v, e := msgId.Int64()
	return v, e
}

func (self RPCMessage) IsRequest() bool {
	return self.Id != nil &&
		self.ServiceName != "" &&
		self.Method != ""
}

func (self RPCMessage) IsNotify() bool {
	return self.Id == nil &&
		self.ServiceName != "" &&
		self.Method != ""
}

func (self RPCMessage) IsResult() bool {
	return (self.Id != nil &&
		self.ServiceName == "" &&
		self.Method == "" &&
		self.Result.Interface() != nil)
}

func (self RPCMessage) IsError() bool {
	return (self.Id != nil &&
		self.ServiceName == "" &&
		self.Method == "" &&
		self.Error.Interface() != nil)
}

func (self RPCMessage) IsResultOrError() bool {
	return (self.Id != nil &&
		self.ServiceName == "" &&
		self.Method == "" &&
		(self.Result.Interface() != nil || self.Error.Interface() != nil))
}

func (self RPCMessage) IsValid() bool {
	return self.IsRequest() || self.IsNotify() || self.IsResultOrError()
}

func (self RPCMessage) GetParams() []interface{} {
	return self.Params.MustArray()
}


