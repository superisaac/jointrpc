package jsonrpc

import (
	//"encoding/json"
	"errors"
	"github.com/bitly/go-simplejson"
)

func ParseMessage(data []byte) (*RPCMessage, error) {
	// Reserved for other data format
	parsed, err := simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return NewRPCMessage(parsed), nil
}

func MarshalJson(json_data *simplejson.Json) (string, error) {
	bytes, err := json_data.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func NewRPCMessage(data *simplejson.Json) *RPCMessage {
	//msg := new(RPCMessage)
	msg := &RPCMessage{Initialized: true}
	//msg.Id = data.Get("id").Interface()
	msgId, err := data.Get("id").Int64()
	if err != nil {
		// TODO: print msg.Id
		msg.Id = 0
	} else {
		msg.Id = UID(msgId)
	}
		
	method, err := data.Get("method").String()
	if err == nil {
		msg.Method = method
	}
	msg.Params = data.Get("params")
	msg.Result = data.Get("result")
	msg.Error = data.Get("error")
	msg.Raw = data
	return msg
}

func NewResultMessage(id interface{}, result interface{}) *RPCMessage {
	resultJson := simplejson.New()
	resultJson.Set("id", id)
	resultJson.Set("result", result)
	return NewRPCMessage(resultJson)
}

func NewNotifyMessage(method string, params []interface{}) *RPCMessage {
	notifyJson := simplejson.New()
	notifyJson.Set("method", method)
	notifyJson.Set("params", params)
	return NewRPCMessage(notifyJson)
}

func NewErrorMessage(id interface{}, code int, message string) *RPCMessage {
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

func (self RPCMessage) GetIntId() (UID, error) {
	//msgId, ok := self.Id.(json.Number)
	if self.Id == 0 {
		return 0, errors.New("not an int64 id")
	}
	return self.Id, nil
	
	// if !ok {
	// 	return 0, errors.New("not a number")
	// }
	// v, e := msgId.Int64()
	// return v, e
}

func (self RPCMessage) IsRequest() bool {
	//return self.Id != nil && self.Method != ""
	return self.Id != 0 && self.Method != ""
}

func (self RPCMessage) IsNotify() bool {
	return self.Id == 0 && self.Method != ""
}

func (self RPCMessage) IsResult() bool {
	return (self.Id != 0 &&
		self.Method == "" &&
		self.Result.Interface() != nil)
}

func (self RPCMessage) IsError() bool {
	return (self.Id != 0 &&
		self.Method == "" &&
		self.Error.Interface() != nil)
}

func (self RPCMessage) IsResultOrError() bool {
	return (self.Id != 0 &&
		self.Method == "" &&
		(self.Result.Interface() != nil || self.Error.Interface() != nil))
}

func (self RPCMessage) IsValid() bool {
	return self.IsRequest() || self.IsNotify() || self.IsResultOrError()
}

func (self RPCMessage) GetParams() []interface{} {
	return self.Params.MustArray()
}
