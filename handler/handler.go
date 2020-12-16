package handler

import (
	"errors"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	//tube "github.com/superisaac/rpctube/tube"
)

func (self *HandlerManager) InitHandlerManager() {
	self.ChResultMsg = make(chan *jsonrpc.RPCMessage)	
	self.MethodHandlers = make(map[string](MethodHandler))
}

func (self *HandlerManager) On(method string, handler HandlerFunc, opts ...func(*MethodHandler)) {
	h := MethodHandler{function: handler, Concurrent: false}
	for _, opt := range opts {
		opt(&h)
	}

	_, found := self.MethodHandlers[method]
	self.MethodHandlers[method] = h

	if !found {
		self.OnHandlerChanged()
	}
}

func (self *HandlerManager) UnHandle(method string) bool {
	_, found := self.MethodHandlers[method]
	if found {
		delete(self.MethodHandlers, method)
		self.OnHandlerChanged()
	}
	return found
}


func (self *HandlerManager) OnHandlerChanged() {
	err := errors.New("On Handler Changed not implemented")
	panic(err)
}

func (self *HandlerManager) wrapHandlerResult(msg *jsonrpc.RPCMessage, res interface{}, err error) (*jsonrpc.RPCMessage, error) {
	if err != nil {
		if rpcErr, ok := err.(*jsonrpc.RPCError); ok {
			return rpcErr.ToMessage(msg.Id), nil
		}
		return nil, err
	} else if msg.IsRequest() {
		return jsonrpc.NewResultMessage(msg.Id, res), nil
	} else {
		return nil, nil
	}
}


func (self *HandlerManager) ReturnResultMessage(resmsg *jsonrpc.RPCMessage) {
	self.ChResultMsg <- resmsg
}

func (self *HandlerManager) HandleRequestMessage(msg *jsonrpc.RPCMessage) {
	handler, ok := self.MethodHandlers[msg.Method]
	var resmsg *jsonrpc.RPCMessage
	var err error
	if ok {
		req := &RPCRequest{Message: msg}
		params := msg.Params.MustArray()
		res, err := handler.function(req, params)
		resmsg, err = self.wrapHandlerResult(msg, res, err)
	} else if self.defaultHandler != nil {
		req := &RPCRequest{Message: msg}

		params := msg.Params.MustArray()
		res, err := self.defaultHandler(req, msg.Method, params)
		resmsg, err = self.wrapHandlerResult(msg, res, err)
	} else {
		resmsg, err = jsonrpc.ErrNoSuchMethod.ToMessage(msg.Id), nil
	}

	if err != nil {
		log.Warnf("bad up message %w", err)
		errmsg := jsonrpc.NewErrorMessage(msg.Id, 10401, "bad handler res", false)
		self.ReturnResultMessage(errmsg)
		return
	}
	if r := recover(); r != nil {
		log.Fatalf("error on handling request msg %+v", r)
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
