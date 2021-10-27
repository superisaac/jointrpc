package jsonrpc

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
)

func (self *RPCError) WithData(data interface{}) *RPCError {
	return &RPCError{self.Code, self.Message, data}
}

func (self *RPCError) Error() string {
	return fmt.Sprintf("code=%d, message=%s, data=%s", self.Code, self.Message, self.Data)
}

func (self RPCError) ToMessage(reqmsg IMessage) *ErrorMessage {
	return RPCErrorMessage(reqmsg, self.Code, self.Message, self.Data)
}

func (self RPCError) ToJson() *simplejson.Json {
	json := simplejson.New()
	json.Set("code", self.Code)
	json.Set("message", self.Message)
	if self.Data != nil {
		json.Set("data", self.Data)
	}
	return json
}

func (self RPCError) CodeString() string {
	return fmt.Sprintf("%d", self.Code)
}

func parseRPCError(errIntf *simplejson.Json) (*RPCError, error) {
	code, err := errIntf.Get("code").Int()
	if err != nil {
		return nil, err
	}

	message, err := errIntf.Get("message").String()
	if err != nil {
		return nil, err
	}

	data := errIntf.Get("data").Interface()
	return &RPCError{code, message, data}, nil
}

func NewErrMsgType(additional string) *RPCError {
	r := fmt.Sprintf("wrong message type %s", additional)
	return &RPCError{ErrMessageType.Code, r, false}
}

func (self *BaseMessage) SetRaw(raw *simplejson.Json) {
	self.raw = raw
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

// IMessage methods
func EncodePretty(msg IMessage) (string, error) {
	bytes, err := MessageJson(msg).EncodePretty()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func MessageJson(msg IMessage) *simplejson.Json {
	jsonObj := msg.GetJson()
	if traceId := msg.TraceId(); traceId != "" {
		jsonObj.Set("traceid", traceId)
	}
	return jsonObj
}

func MessageInterface(msg IMessage) interface{} {
	return MessageJson(msg).Interface()
}

func MessageString(msg IMessage) string {
	bytes, err := MessageBytes(msg)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func MessageBytes(msg IMessage) ([]byte, error) {
	return MessageJson(msg).MarshalJSON()
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
func (self RequestMessage) MustError() *RPCError {
	panic(NewErrMsgType("MustError"))
}
func (self NotifyMessage) MustError() *RPCError {
	panic(NewErrMsgType("MustError"))
}
func (self ResultMessage) MustError() *RPCError {
	panic(NewErrMsgType("MustError"))
}
func (self ErrorMessage) MustError() *RPCError {
	return self.Error
}

// Get Json
func (self *RequestMessage) GetJson() *simplejson.Json {
	if self.raw == nil {
		self.raw = simplejson.New()
		self.raw.Set("version", "2.0")
		self.raw.Set("id", self.Id)
		self.raw.Set("method", self.Method)
		self.raw.Set("params", self.Params)
	}
	return self.raw
}

func (self *NotifyMessage) GetJson() *simplejson.Json {
	if self.raw == nil {
		self.raw = simplejson.New()
		self.raw.Set("version", "2.0")
		self.raw.Set("method", self.Method)
		self.raw.Set("params", self.Params)
	}
	return self.raw
}

func (self *ResultMessage) GetJson() *simplejson.Json {
	if self.raw == nil {
		self.raw = simplejson.New()
		self.raw.Set("version", "2.0")
		self.raw.Set("id", self.Id)
		self.raw.Set("result", self.Result)
	}
	return self.raw
}

func (self *ErrorMessage) GetJson() *simplejson.Json {
	if self.raw == nil {
		self.raw = simplejson.New()
		self.raw.Set("version", "2.0")
		self.raw.Set("id", self.Id)
		self.raw.Set("error", self.Error.ToJson())
	}
	return self.raw
}

func NewRequestMessage(id interface{}, method string, params []interface{}) *RequestMessage {
	if id == nil {
		panic(ErrNilId)
	}
	if method == "" {
		panic(ErrEmptyMethod)
	}

	msg := &RequestMessage{}
	msg.messageType = MTRequest
	msg.Id = id
	msg.Method = method
	msg.Params = params
	return msg
}

func (self RequestMessage) Clone(newId interface{}) *RequestMessage {
	newReq := NewRequestMessage(newId, self.Method, self.Params)
	newReq.SetTraceId(self.traceId)
	return newReq
}

func NewNotifyMessage(method string, params []interface{}) *NotifyMessage {
	if method == "" {
		panic(ErrEmptyMethod)
	}

	msg := &NotifyMessage{}
	msg.messageType = MTNotify
	msg.Method = method
	msg.Params = params
	return msg
}

func rawResultMessage(id interface{}, result interface{}) *ResultMessage {
	msg := &ResultMessage{}
	msg.messageType = MTResult
	msg.Id = id
	msg.Result = result
	return msg
}

func NewResultMessage(reqmsg IMessage, result interface{}) *ResultMessage {
	resmsg := rawResultMessage(reqmsg.MustId(), result)
	resmsg.SetTraceId(reqmsg.TraceId())
	return resmsg
}

func NewErrorMessage(reqmsg IMessage, errbody *RPCError) *ErrorMessage {
	errmsg := rawErrorMessage(reqmsg.MustId(), errbody)
	errmsg.SetTraceId(reqmsg.TraceId())
	return errmsg
}

func rawErrorMessage(id interface{}, errbody *RPCError) *ErrorMessage {
	msg := &ErrorMessage{}
	msg.messageType = MTError
	msg.Id = id
	msg.Error = errbody
	return msg
}

func RPCErrorMessage(reqmsg IMessage, code int, message string, data interface{}) *ErrorMessage {
	errbody := &RPCError{code, message, data}
	return NewErrorMessage(reqmsg, errbody)
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

	traceId, err := parsed.Get("traceid").String()
	if err != nil {
		traceId = ""
	}

	if id != nil {
		if method != "" {
			// request
			params := parsed.Get("params").MustArray()
			reqmsg := NewRequestMessage(id, method, params)
			reqmsg.SetRaw(parsed)
			reqmsg.SetTraceId(traceId)
			return reqmsg, nil
		}
		if errIntf := parsed.Get("error"); errIntf != nil && errIntf.Interface() != nil {
			errbody, err := parseRPCError(errIntf)
			if err != nil {
				return nil, err
			}
			errmsg := rawErrorMessage(id, errbody)
			errmsg.SetRaw(parsed)
			errmsg.SetTraceId(traceId)
			return errmsg, nil
		}
		res := parsed.Get("result").Interface()
		rmsg := rawResultMessage(id, res)
		rmsg.SetRaw(parsed)
		rmsg.SetTraceId(traceId)
		return rmsg, nil
	} else if method != "" {
		params := parsed.Get("params").MustArray()
		ntfmsg := NewNotifyMessage(method, params)
		ntfmsg.SetRaw(parsed)
		ntfmsg.SetTraceId(traceId)
		return ntfmsg, nil
	} else {
		return nil, ErrParseMessage
	}
}
