package jsonrpc

import (
	"github.com/bitly/go-simplejson"
)

func (self BaseMessage) IsRequest() bool {
	return self.messageType == MTRequest
}

func (self BaseMessage) IsNotify() bool {
	return self.messageType == MTNotify
}

func (self BaseMessage) IsRequestOrNotify() bool {
	return self.IsRequest() || self.IsNotify()
}

func (self BaseMessage) IsResult() bool {
	return self.messageType == MTResult
}
func (self BaseMessage) IsError() bool {
	return self.messageType == MTError
}
func (self BaseMessage) IsResultOrError() bool {
	return self.IsResult() || self.IsError()
}

func (self BaseMessage) EncodePretty() (string, error) {
	bytes, err := self.raw.EncodePretty()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (self BaseMessage) Interface() interface{} {
	return self.raw.Interface()
}

func (self BaseMessage) MustString() string {
	bytes, err := self.Bytes()
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (self BaseMessage) Bytes() ([]byte, error) {
	return self.raw.MarshalJSON()
}

// Must methods
// MustId
func (self RequestMessage) MustId() interface{} {
	return self.Id
}
func (self NotifyMessage) MustId() interface{} {
	panic(ErrMessageType)
}
func (self ResultMessage) MustId() interface{} {
	return self.Id
}
func (self ErrorMessage) MustId() interface{} {
	return self.Id
}

// MustMethod
func (self RequestMessage) MustMethod() string {
	return self.Method
}
func (self NotifyMessage) MustMethod() string {
	return self.Method
}
func (self ResultMessage) MustMethod() string {
	panic(ErrMessageType)
}

func (self ErrorMessage) MustMethod() string {
	panic(ErrMessageType)
}

// MustParams
func (self RequestMessage) MustParams() []interface{} {
	return self.Params
}
func (self NotifyMessage) MustParams() []interface{} {
	return self.Params
}
func (self ResultMessage) MustParams() []interface{} {
	panic(ErrMessageType)
}
func (self ErrorMessage) MustParams() []interface{} {
	panic(ErrMessageType)
}

// MustResult
func (self RequestMessage) MustResult() interface{} {
	panic(ErrMessageType)
}
func (self NotifyMessage) MustResult() interface{} {
	panic(ErrMessageType)
}
func (self ResultMessage) MustResult() interface{} {
	return self.Result
}
func (self ErrorMessage) MustResult() interface{} {
	panic(ErrMessageType)
}

// MustError
func (self RequestMessage) MustError() interface{} {
	panic(ErrMessageType)
}
func (self NotifyMessage) MustError() interface{} {
	panic(ErrMessageType)
}
func (self ResultMessage) MustError() interface{} {
	panic(ErrMessageType)
}
func (self ErrorMessage) MustError() interface{} {
	return self.Error
}

func NewRequestMessage(id interface{}, method string, params []interface{}, raw *simplejson.Json) *RequestMessage {
	if id == nil {
		panic(&RPCError{10400, "request message cannot have nil id", false})
	}
	if method == "" {
		panic(&RPCError{10400, "request message method cannot be empty", false})
	}

	if raw == nil {
		raw = simplejson.New()
		raw.Set("version", "2.0")
		raw.Set("id", id)
		raw.Set("method", method)
		raw.Set("params", params)
	}
	msg := &RequestMessage{}
	msg.messageType = MTRequest
	msg.raw = raw
	msg.Id = id
	msg.Method = method
	msg.Params = params
	return msg
}

func NewNotifyMessage(method string, params []interface{}, raw *simplejson.Json) *NotifyMessage {
	if method == "" {
		panic(&RPCError{10400, "notify message method cannot be empty", false})
	}

	if raw == nil {
		raw = simplejson.New()
		raw.Set("version", "2.0")
		raw.Set("method", method)
		raw.Set("params", params)
	}
	msg := &NotifyMessage{}
	msg.messageType = MTNotify
	msg.raw = raw
	msg.Method = method
	msg.Params = params
	return msg
}

func NewResultMessage(id interface{}, result interface{}, raw *simplejson.Json) *ResultMessage {
	if id == nil {
		panic(&RPCError{10400, "result message id cannot be nil", false})
	}

	if raw == nil {
		raw = simplejson.New()
		raw.Set("version", "2.0")
		raw.Set("id", id)
		raw.Set("result", result)
	}
	msg := &ResultMessage{}
	msg.messageType = MTResult
	msg.raw = raw
	msg.Id = id
	msg.Result = result
	return msg
}

func NewErrorMessage(id interface{}, errbody interface{}, raw *simplejson.Json) *ErrorMessage {
	if id == nil {
		panic(&RPCError{10400, "error message id cannot be nil", false})
	}

	if raw == nil {
		raw = simplejson.New()
		raw.Set("version", "2.0")
		raw.Set("id", id)
		raw.Set("error", errbody)
	}

	msg := &ErrorMessage{}
	msg.messageType = MTError
	msg.raw = raw
	msg.Id = id
	msg.Error = errbody
	return msg
}

func RPCErrorMessage(id interface{}, code int, reason string, retryable bool) *ErrorMessage {
	errbody := map[string](interface{}){
		"code":      code,
		"reason":    reason,
		"retryable": retryable}
	return NewErrorMessage(id, errbody, nil)
}

func ParseBytes(data []byte) (IMessage, error) {
	parsed, err := simplejson.NewJson(data)
	if err != nil {
		return nil, err
	}
	return Parse(parsed)
}
func Parse(parsed *simplejson.Json) (IMessage, error) {
	id := parsed.Get("id").Interface()
	method, err := parsed.Get("method").String()
	if err != nil {
		method = ""
	}

	if id != nil {
		if method != "" {
			// request
			params := parsed.Get("params").MustArray()
			return NewRequestMessage(id, method, params, parsed), nil
		}
		if errIntf := parsed.Get("error").Interface(); errIntf != nil {
			return NewErrorMessage(id, errIntf, parsed), nil
		}
		res := parsed.Get("result").Interface()
		return NewResultMessage(id, res, parsed), nil
	} else if method != "" {
		params := parsed.Get("params").MustArray()
		return NewNotifyMessage(method, params, parsed), nil
	} else {
		return nil, &RPCError{10402, "parse JSONRPC error", false}
	}
}
