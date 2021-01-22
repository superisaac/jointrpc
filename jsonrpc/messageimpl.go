package jsonrpc

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
)

func NewErrMsgType(additional string) *RPCError {
	r := fmt.Sprintf("wrong message type %s", additional)
	return &RPCError{10403, r, false}
}

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

func (self *BaseMessage) SetTraceId(traceId string) {
	self.traceId = traceId
}

func (self BaseMessage) TraceId() string {
	return self.traceId
}

// Log
func (self RequestMessage) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"trace_id": self.traceId,
		"msg_type": "request",
		"msg_id":   self.Id,
		"method":   self.Method,
	})
}
func (self NotifyMessage) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"trace_id": self.traceId,
		"msg_type": "notify",
		"method":   self.Method,
	})
}
func (self ResultMessage) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"trace_id": self.traceId,
		"msg_type": "result",
		"msg_id":   self.Id,
	})
}

func (self ErrorMessage) Log() *log.Entry {
	return log.WithFields(log.Fields{
		"trace_id": self.traceId,
		"msg_type": "error",
		"msg_id":   self.Id,
	})
}

// Must methods

// MustId
func (self RequestMessage) MustId() interface{} {
	return self.Id
}
func (self NotifyMessage) MustId() interface{} {
	panic(NewErrMsgType("MustId"))
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
	panic(NewErrMsgType("MustMethod"))
}

func (self ErrorMessage) MustMethod() string {
	panic(NewErrMsgType("MustMethod"))
}

// MustParams
func (self RequestMessage) MustParams() []interface{} {
	return self.Params
}
func (self NotifyMessage) MustParams() []interface{} {
	return self.Params
}
func (self ResultMessage) MustParams() []interface{} {
	panic(NewErrMsgType("MustParams"))
}
func (self ErrorMessage) MustParams() []interface{} {
	panic(NewErrMsgType("MustParams"))
}

// MustResult
func (self RequestMessage) MustResult() interface{} {
	panic(NewErrMsgType("MustResult"))
}
func (self NotifyMessage) MustResult() interface{} {
	panic(NewErrMsgType("MustResult"))
}
func (self ResultMessage) MustResult() interface{} {
	return self.Result
}
func (self ErrorMessage) MustResult() interface{} {
	panic(NewErrMsgType("MustResult"))
}

// MustError
func (self RequestMessage) MustError() interface{} {
	panic(NewErrMsgType("MustError"))
}
func (self NotifyMessage) MustError() interface{} {
	panic(NewErrMsgType("MustError"))
}
func (self ResultMessage) MustError() interface{} {
	panic(NewErrMsgType("MustError"))
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

func (self RequestMessage) Clone(newId interface{}) *RequestMessage {
	newReq := NewRequestMessage(newId, self.Method, self.Params, nil)
	newReq.SetTraceId(self.traceId)
	return newReq
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

func rawResultMessage(id interface{}, result interface{}, raw *simplejson.Json) *ResultMessage {
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

func NewResultMessage(reqmsg IMessage, result interface{}, raw *simplejson.Json) *ResultMessage {
	resmsg := rawResultMessage(reqmsg.MustId(), result, raw)
	resmsg.SetTraceId(reqmsg.TraceId())
	return resmsg
}

func NewErrorMessage(reqmsg IMessage, errbody interface{}, raw *simplejson.Json) *ErrorMessage {
	errmsg := rawErrorMessage(reqmsg.MustId(), errbody, raw)
	errmsg.SetTraceId(reqmsg.TraceId())
	return errmsg
}

func rawErrorMessage(id interface{}, errbody interface{}, raw *simplejson.Json) *ErrorMessage {
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

func RPCErrorMessage(reqmsg IMessage, code int, reason string, retryable bool) *ErrorMessage {

	errbody := map[string](interface{}){
		"code":      code,
		"reason":    reason,
		"retryable": retryable}
	return NewErrorMessage(reqmsg, errbody, nil)
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
			return rawErrorMessage(id, errIntf, parsed), nil
		}
		res := parsed.Get("result").Interface()
		return rawResultMessage(id, res, parsed), nil
	} else if method != "" {
		params := parsed.Get("params").MustArray()
		return NewNotifyMessage(method, params, parsed), nil
	} else {
		return nil, &RPCError{10402, "parse JSONRPC error", false}
	}
}
