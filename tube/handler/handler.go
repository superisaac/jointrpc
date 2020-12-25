package handler

import (
	//"errors"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	tube "github.com/superisaac/rpctube/tube"
)

func (self *HandlerManager) InitHandlerManager() {
	self.ChResultMsg = make(chan jsonrpc.IMessage, 100)
	self.MethodHandlers = make(map[string](MethodHandler))
}

func (self *HandlerManager) On(method string, handler HandlerFunc, opts ...func(*MethodHandler)) {
	h := MethodHandler{function: handler, Concurrent: false}
	for _, opt := range opts {
		opt(&h)
	}

	_, found := self.MethodHandlers[method]
	self.MethodHandlers[method] = h

	if !found && self.onChange != nil {
		self.onChange()
	}
}

func (self *HandlerManager) OnChange(onchange OnChangeFunc) {
	self.onChange = onchange
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
		return jsonrpc.NewResultMessage(msg.MustId(), res, nil), nil
	} else {
		return nil, nil
	}
}

func (self *HandlerManager) ReturnResultMessage(resmsg jsonrpc.IMessage) {
	self.ChResultMsg <- resmsg
}

func (self *HandlerManager) HandleRequestMessage(msgvec tube.MsgVec) {
	msg := msgvec.Msg
	handler, ok := self.MethodHandlers[msg.MustMethod()]

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered ERROR on handling request msg %+v", r)
			if rpcError, ok := r.(*jsonrpc.RPCError); ok {
				errmsg := rpcError.ToMessage(msg.MustId())
				self.ReturnResultMessage(errmsg)
				return
			} else {
				errmsg := jsonrpc.ErrServerError.ToMessage(msg.MustId())
				self.ReturnResultMessage(errmsg)
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
	if err != nil {
		log.Warnf("bad up message %w", err)
		errmsg := jsonrpc.RPCErrorMessage(msg.MustId(), 10401, "bad handler res", false)
		self.ReturnResultMessage(errmsg)
		return
	}
	if resmsg != nil {
		self.ReturnResultMessage(resmsg)
	}
}

// MethodHandler Helper methods
func WithHelp(help string) func(*MethodHandler) {
	return func(h *MethodHandler) {
		h.Help = help
	}
}

func WithSchema(schema string) func(*MethodHandler) {
	return func(h *MethodHandler) {
		// TODO: parse schema
		h.Schema = schema
	}
}

func WithConcurrent(c bool) func(*MethodHandler) {
	return func(h *MethodHandler) {
		// TODO: parse schema
		h.Concurrent = c
	}
}

func (self *HandlerManager) OnDefault(handler DefaultHandlerFunc, opts ...func(*HandlerManager)) {
	self.defaultHandler = handler
	for _, opt := range opts {
		opt(self)
	}
}

func (self HandlerManager) CanRunConcurrent(method string) bool {
	handler, ok := self.MethodHandlers[method]
	if ok {
		return handler.Concurrent
	} else if self.defaultConcurrent {
		return true
	}
	return false
}

func WithDefaultConcurrent(c bool) func(*HandlerManager) {
	return func(h *HandlerManager) {
		// TODO: parse schema
		h.defaultConcurrent = c
	}
}
