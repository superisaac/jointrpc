package dispatch

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	misc "github.com/superisaac/jointrpc/misc"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
	"github.com/superisaac/jsonz"
	schema "github.com/superisaac/jsonz/schema"
)

const (
	recoverPanic = false
)

// methods of RPCRequest
func NewRPCRequest(ctx context.Context, cmdMsg rpcrouter.CmdMsg) *RPCRequest {
	return &RPCRequest{Context: ctx, CmdMsg: cmdMsg}
}

func (self *RPCRequest) WithData(data interface{}) *RPCRequest {
	self.Data = data
	return self
}

func WithRequestData(data interface{}) func(*RPCRequest) {
	return func(req *RPCRequest) {
		req.WithData(data)
	}
}

// methods of dispatcher
func NewDispatcher() *Dispatcher {
	disp := new(Dispatcher)
	disp.methodHandlers = make(map[string](MethodHandler))
	disp.changeHandlers = make([]OnChangeFunc, 0)
	return disp
}

func (self Dispatcher) HasMethod(method string) bool {
	_, ok := self.methodHandlers[method]
	return ok
}

func (self Dispatcher) GetMethodInfos() []rpcrouter.MethodInfo {
	minfos := make([]rpcrouter.MethodInfo, 0)
	for m, info := range self.methodHandlers {
		minfo := rpcrouter.MethodInfo{
			Name:       m,
			Help:       info.Help,
			SchemaJson: info.SchemaJson,
		}
		minfos = append(minfos, minfo)
	}
	return minfos
}

func (self Dispatcher) GetPublicMethodInfos() []rpcrouter.MethodInfo {
	minfos := make([]rpcrouter.MethodInfo, 0)
	for _, minfo := range self.GetMethodInfos() {
		if jsonz.IsPublicMethod(minfo.Name) {
			minfos = append(minfos, minfo)
		}
	}
	return minfos
}

func (self *Dispatcher) On(method string, handler HandlerFunc, opts ...func(*MethodHandler)) {
	if !jsonz.IsMethod(method) {
		panic(errors.New("invalid method name"))
	}
	h := MethodHandler{function: handler}
	for _, opt := range opts {
		opt(&h)
	}

	_, found := self.methodHandlers[method]
	self.methodHandlers[method] = h

	if !found && len(self.changeHandlers) == 0 {
		self.TriggerChange()
	}
}

func (self *Dispatcher) OnTyped(method string, typedHandler interface{}, opts ...func(*MethodHandler)) {
	handler, err := WrapTyped(typedHandler)
	if err != nil {
		panic(err)
	}
	self.On(method, handler, opts...)
}

func (self *Dispatcher) SetSpawnExec(v bool) {
	self.spawnExec = v
}

func (self *Dispatcher) OnChange(onChange OnChangeFunc) {
	self.changeHandlers = append(self.changeHandlers, onChange)
}

func (self *Dispatcher) TriggerChange() {
	for _, changeFunc := range self.changeHandlers {
		changeFunc()
	}
}

func (self *Dispatcher) UnHandle(method string) bool {
	_, found := self.methodHandlers[method]
	if found {
		delete(self.methodHandlers, method)
		self.TriggerChange()
	}
	return found
}

func (self *Dispatcher) wrapHandlerResult(reqmsg *jsonz.RequestMessage, res interface{}, err error) (jsonz.Message, error) {
	if err != nil {
		var rpcErr *jsonz.RPCError
		if errors.As(err, &rpcErr) {
			return rpcErr.ToMessage(reqmsg), nil
		}
		reqmsg.Log().Warnf("error %s", err.Error())
		errmsg := jsonz.ErrServerError.ToMessage(reqmsg)
		return errmsg, nil
		//return , err
	}

	if resMsg, ok := res.(jsonz.Message); ok {
		// TODO: assert resMsg is res and resId matches
		return resMsg, nil
	}
	return jsonz.NewResultMessage(reqmsg, res), nil
}

func (self *Dispatcher) ReturnResultMessage(resmsg jsonz.Message, req rpcrouter.CmdMsg, chResult chan ResultT) {
	chResult <- ResultT{
		ResMsg:    resmsg,
		ReqCmdMsg: req,
	}
}

func (self *Dispatcher) Expect(ctx context.Context, cmdMsg rpcrouter.CmdMsg, opts ...func(*RPCRequest)) jsonz.Message {
	chResult := make(chan ResultT, 2)
	self.Feed(ctx, cmdMsg, chResult, opts...)
	res := <-chResult
	return res.ResMsg
}

func (self *Dispatcher) Feed(ctx context.Context, cmdMsg rpcrouter.CmdMsg, chResult chan ResultT, opts ...func(*RPCRequest)) {
	req := NewRPCRequest(ctx, cmdMsg)
	for _, opt := range opts {
		opt(req)
	}
	if self.spawnExec {
		go self.feed(req, chResult)
	} else {
		self.feed(req, chResult)
	}
}

func (self *Dispatcher) feed(req *RPCRequest, chResult chan ResultT) {
	if req.CmdMsg.Msg.IsRequest() {
		self.feedRequest(req, chResult)
	} else {
		misc.Assert(req.CmdMsg.Msg.IsNotify(), "invalid msg type")
		self.feedNotify(req, chResult)
	}
}

func (self *Dispatcher) feedRequest(req *RPCRequest, chResult chan ResultT) {
	reqmsg, ok := req.CmdMsg.Msg.(*jsonz.RequestMessage)

	misc.Assert(ok, "msg is not request")

	handler, ok := self.methodHandlers[reqmsg.Method]

	if recoverPanic {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("recovered r %+v\n", r)

				if err, ok := r.(error); ok {
					var rpcError *jsonz.RPCError
					if errors.Is(err, Deferred) {
						reqmsg.Log().Debugf("handler is deferred")
						return
						//} else if rpcError, ok := r.(*jsonz.RPCError); ok {
					} else if errors.As(err, &rpcError) {
						errmsg := rpcError.ToMessage(reqmsg)
						self.ReturnResultMessage(errmsg, req.CmdMsg, chResult)
						return
					}
				}
				reqmsg.Log().Errorf("Recovered ERROR on handling request msg %+v", r)
			}
		}()
	}

	var resmsg jsonz.Message
	var err error
	if ok {
		res, err := handler.function(req, reqmsg.Params)
		//log.Debugf("handler function returns %+v, %+v", reqmsg, res)
		resmsg, err = self.wrapHandlerResult(reqmsg, res, err)
	} else if self.defaultHandler != nil {
		res, err := self.defaultHandler(req, reqmsg.Method, reqmsg.Params)
		resmsg, err = self.wrapHandlerResult(reqmsg, res, err)
	} else {
		resmsg, err = jsonz.ErrMethodNotFound.ToMessage(reqmsg), nil
	}

	//log.Debugf("handle request method %+v, resmsg %+v, error %+v", msg, resmsg, err)
	if errors.Is(err, Deferred) {
		log.Infof("handler is deferred")
		return
	}
	if err != nil {
		log.Warnf("bad up message %s", err)
		errmsg := jsonz.ErrBadResource.ToMessage(reqmsg)
		self.ReturnResultMessage(errmsg, req.CmdMsg, chResult)
		return
	}
	if resmsg != nil {
		self.ReturnResultMessage(resmsg, req.CmdMsg, chResult)
	}
}

func (self *Dispatcher) feedNotify(req *RPCRequest, chResult chan ResultT) {
	ntfmsg, ok := req.CmdMsg.Msg.(*jsonz.NotifyMessage)
	misc.Assert(ok, "message is not ok")
	handler, ok := self.methodHandlers[ntfmsg.Method]

	if recoverPanic {
		defer func() {
			if r := recover(); r != nil {
				if r == Deferred {
					ntfmsg.Log().Debugf("handler is deferred")
					return
				} else if rpcError, ok := r.(*jsonz.RPCError); ok {
					ntfmsg.Log().Warnf("RPCError code=%d, message=%s", rpcError.Code, rpcError.Message)
					return
				} else {
					ntfmsg.Log().Errorf("Recovered ERROR on handling notify msg %+v", r)
				}
			}
		}()
	}

	var res interface{}
	var err error
	if ok {
		res, err = handler.function(req, ntfmsg.Params)
		if res != nil {
			ntfmsg.Log().Infof("res is not nil, %+v", res)
		}
	} else if self.defaultHandler != nil {

		res, err = self.defaultHandler(req, ntfmsg.Method, ntfmsg.Params)
		if res != nil {
			ntfmsg.Log().Infof("res of default handler is not nil, %+v", res)
		}
	}

	//log.Debugf("handle request method %+v, resmsg %+v, error %+v", msg, resmsg, err)
	if errors.Is(err, Deferred) {
		log.Infof("handler is deferred")
	} else if err != nil {
		log.Warnf("bad up message %s", err)
	}
}

// MethodHandler Helper methods
func WithHelp(help string) func(*MethodHandler) {
	return func(h *MethodHandler) {
		h.Help = help
	}
}

func WithSchema(schemaJson string) func(*MethodHandler) {
	return func(h *MethodHandler) {
		if schemaJson != "" {
			// TODO: build schema
			builder := schema.NewSchemaBuilder()
			_, err := builder.BuildBytes([]byte(schemaJson))
			if err != nil {
				panic(err)
			}
		}
		h.SchemaJson = schemaJson
	}
}

func (self *Dispatcher) OnDefault(handler DefaultHandlerFunc, opts ...func(*Dispatcher)) {
	self.defaultHandler = handler
	for _, opt := range opts {
		opt(self)
	}
}
