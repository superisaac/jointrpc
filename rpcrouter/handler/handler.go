package handler

import (
	//"errors"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
)

// handler manager
func (self *HandlerManager) InitHandlerManager() {
	self.ChResult = make(chan MsgEnvo, 100)
	self.MethodHandlers = make(map[string](MethodHandler))
}

func (self *HandlerManager) On(method string, handler HandlerFunc, opts ...func(*MethodHandler)) {
	h := MethodHandler{function: handler}
	for _, opt := range opts {
		opt(&h)
	}

	_, found := self.MethodHandlers[method]
	self.MethodHandlers[method] = h

	if !found && self.onChange != nil {
		self.onChange()
	}
}

func (self *HandlerManager) OnChange(onChange OnChangeFunc) {
	self.onChange = onChange
}

func (self *HandlerManager) OnStateChange(onChange StateHandlerFunc) {
	self.StateHandler = onChange
}

func (self *HandlerManager) UnHandle(method string) bool {
	_, found := self.MethodHandlers[method]
	if found {
		delete(self.MethodHandlers, method)
		if self.onChange != nil {
			self.onChange()
		}
	}
	return found
}

func (self *HandlerManager) wrapHandlerResult(msg jsonrpc.IMessage, res interface{}, err error) (jsonrpc.IMessage, error) {
	if err != nil {
		if rpcErr, ok := err.(*jsonrpc.RPCError); ok {
			return rpcErr.ToMessage(msg.MustId()), nil
		}
		return nil, err
	} else if msg.IsRequest() {
		log.Debugf("msg is request")
		if resMsg, ok := res.(jsonrpc.IMessage); ok {
			// TODO: assert resMsg is res and resId matches
			return resMsg, nil
		}
		return jsonrpc.NewResultMessage(msg.MustId(), res, nil), nil
	} else {
		return nil, nil
	}
}

func (self *HandlerManager) ReturnResultMessage(resmsg jsonrpc.IMessage, traceId string) {
	self.ChResult <- MsgEnvo{Msg: resmsg, TraceId: traceId}
}

func (self *HandlerManager) HandleRequestMessage(msgvec rpcrouter.MsgVec) {
	msg := msgvec.Msg
	handler, ok := self.MethodHandlers[msg.MustMethod()]

	defer func() {
		if r := recover(); r != nil {

			if r == Deferred {
				log.Infof("handler is deferred")
				return
			} else if rpcError, ok := r.(*jsonrpc.RPCError); ok {
				errmsg := rpcError.ToMessage(msg.MustId())
				self.ReturnResultMessage(errmsg, msgvec.TraceId)
				return
			} else {
				log.Errorf("Recovered ERROR on handling request msg %+v", r)
				errmsg := jsonrpc.ErrServerError.ToMessage(msg.MustId())
				self.ReturnResultMessage(errmsg, msgvec.TraceId)
			}
		}
	}()

	var resmsg jsonrpc.IMessage
	var err error
	if ok {
		req := &RPCRequest{MsgVec: msgvec}
		params := msg.MustParams()
		res, err := handler.function(req, params)
		log.Debugf("handler function returns %+v, %+v", msg, res)
		resmsg, err = self.wrapHandlerResult(msg, res, err)
	} else if self.defaultHandler != nil {
		req := &RPCRequest{MsgVec: msgvec}

		params := msg.MustParams()
		res, err := self.defaultHandler(req, msg.MustMethod(), params)
		resmsg, err = self.wrapHandlerResult(msg, res, err)
	} else {
		resmsg, err = jsonrpc.ErrNoSuchMethod.ToMessage(msg.MustId()), nil
	}

	//log.Debugf("handle request method %+v, resmsg %+v, error %+v", msg, resmsg, err)
	if err == Deferred {
		log.Infof("handler is deferred")
		return
	}
	if err != nil {
		log.Warnf("bad up message %w", err)
		errmsg := jsonrpc.RPCErrorMessage(msg.MustId(), 10401, "bad handler res", false)
		self.ReturnResultMessage(errmsg, msgvec.TraceId)
		return
	}
	if resmsg != nil {
		self.ReturnResultMessage(resmsg, msgvec.TraceId)
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

func (self *HandlerManager) OnDefault(handler DefaultHandlerFunc, opts ...func(*HandlerManager)) {
	self.defaultHandler = handler
	for _, opt := range opts {
		opt(self)
	}
}
