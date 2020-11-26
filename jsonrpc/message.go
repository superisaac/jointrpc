package jsonrpc

import (
	//"encoding/json"
	"errors"
	"strconv"
	simplejson "github.com/bitly/go-simplejson"
)

func GuessJson(input string) (interface{}, error) {
	if len(input) == 0 {
		return "", nil
	}
	if input == "true" || input == "false" {
		bv, err := strconv.ParseBool(input)
		if err != nil {
			return nil, err
		}
		return bv, nil
	}
	iv, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		return iv, nil
	}
	fv, err := strconv.ParseFloat(input, 64)
	if err == nil {
		return fv, nil
	}
	
	fc := input[0]
	if fc == '[' {
		parsed, err := simplejson.NewJson([]byte(input))
		if err != nil {
			return nil, err
		}
		return parsed.MustArray(), nil
	} else if fc == '{' {
		parsed, err := simplejson.NewJson([]byte(input))
		if err != nil {
			return nil, err
		}
		return parsed.MustMap(), nil
	} else {
		return input, nil
	}	
}

func GuessJsonArray(inputArr []string) ([]interface{}, error) {
	var arr []interface{}
	for _, input := range inputArr {
		v, err := GuessJson(input)
		if err != nil {
			return arr, err
		}
		arr = append(arr, v)
	}
	return arr, nil
}

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

func NewErrorMessage(id interface{}, code int, message string, retryable bool) *RPCMessage {
	jsonData := NewErrorJSON(id, code, message, retryable)
	return NewRPCMessage(jsonData)
}

func NewErrorJSON(id interface{}, code int, message string, retryable bool) *simplejson.Json {
	// Retryable indicates whether the client can retry the request using the same args
	// Usually the parameter is used in case of network failure.
	errJson := simplejson.New()
	errJson.Set("code", code)
	errJson.Set("message", message)

	errJson.Set("retryable", retryable)
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
