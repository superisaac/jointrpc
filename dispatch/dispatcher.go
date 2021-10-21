package dispatch

import (
	"context"
	"errors"
	//"fmt"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	misc "github.com/superisaac/jointrpc/misc"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
)

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

func (self *Dispatcher) On(method string, handler HandlerFunc, opts ...func(*MethodHandler)) {
	if !jsonrpc.IsMethod(method) {
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

func (self *Dispatcher) wrapHandlerResult(msg jsonrpc.IMessage, res interface{}, err error) (jsonrpc.IMessage, error) {
	if err != nil {
		if rpcErr, ok := err.(*jsonrpc.RPCError); ok {
			return rpcErr.ToMessage(msg), nil
		}
		msg.Log().Warnf("error %s", err.Error())
		errmsg := jsonrpc.ErrServerError.ToMessage(msg)
		return errmsg, nil
		//return , err
	} else if msg.IsRequest() {
		if resMsg, ok := res.(jsonrpc.IMessage); ok {
			// TODO: assert resMsg is res and resId matches
			return resMsg, nil
		}
		return jsonrpc.NewResultMessage(msg, res, nil), nil
	} else {
		return nil, nil
	}
}

func (self *Dispatcher) ReturnResultMessage(resmsg jsonrpc.IMessage, req rpcrouter.MsgVec, chResult chan ResultT) {
	chResult <- ResultT{
		ResMsg:    resmsg,
		ReqMsgVec: req,
	}
}

func (self *Dispatcher) Expect(ctx context.Context, msgvec rpcrouter.MsgVec) jsonrpc.IMessage {
	chResult := make(chan ResultT, 2)
	self.Feed(ctx, msgvec, chResult)
	res := <-chResult
	return res.ResMsg
}

func (self *Dispatcher) Feed(ctx context.Context, msgvec rpcrouter.MsgVec, chResult chan ResultT) {
	if self.spawnExec {
		go self.feed(ctx, msgvec, chResult)
	} else {
		self.feed(ctx, msgvec, chResult)
	}
}

func (self *Dispatcher) feed(ctx context.Context, msgvec rpcrouter.MsgVec, chResult chan ResultT) {
	msg := msgvec.Msg
	namespace := msgvec.Namespace
	misc.Assert(namespace != "", "empty namespace")

	handler, ok := self.methodHandlers[msg.MustMethod()]

	defer func() {
		if r := recover(); r != nil {

			if r == Deferred {
				msg.Log().Infof("handler is deferred")
				return
			} else if rpcError, ok := r.(*jsonrpc.RPCError); ok {
				if msg.IsRequest() {
					errmsg := rpcError.ToMessage(msg)
					self.ReturnResultMessage(errmsg, msgvec, chResult)
				} else {
					msg.Log().Warnf("RPCError code=%s, message=%s", rpcError.Code, rpcError.Message)
					//fmt.Printf("RPCError code=%s, message=%s\n", rpcError.Code, rpcError.Message)
				}
				return
			} else {
				msg.Log().Errorf("Recovered ERROR on handling request msg %+v", r)
				if msg.IsRequest() {
					errmsg := jsonrpc.ErrServerError.ToMessage(msg)
					self.ReturnResultMessage(errmsg, msgvec, chResult)
				}
			}
		}
	}()

	var resmsg jsonrpc.IMessage
	var err error
	if ok {
		req := &RPCRequest{Context: ctx, MsgVec: msgvec}
		params := msg.MustParams()
		res, err := handler.function(req, params)
		log.Debugf("handler function returns %+v, %+v", msg, res)
		resmsg, err = self.wrapHandlerResult(msg, res, err)
	} else if self.defaultHandler != nil {
		req := &RPCRequest{Context: ctx, MsgVec: msgvec}

		params := msg.MustParams()
		res, err := self.defaultHandler(req, msg.MustMethod(), params)
		resmsg, err = self.wrapHandlerResult(msg, res, err)
	} else {
		resmsg, err = jsonrpc.ErrMethodNotFound.ToMessage(msg), nil
	}

	//log.Debugf("handle request method %+v, resmsg %+v, error %+v", msg, resmsg, err)
	if err == Deferred {
		log.Infof("handler is deferred")
		return
	}
	if err != nil {
		log.Warnf("bad up message %w", err)
		errmsg := jsonrpc.ErrBadResource.ToMessage(msg)
		self.ReturnResultMessage(errmsg, msgvec, chResult)
		return
	}
	if resmsg != nil {
		self.ReturnResultMessage(resmsg, msgvec, chResult)
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
