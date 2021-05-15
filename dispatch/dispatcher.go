package dispatch

import (
	"errors"
	log "github.com/sirupsen/logrus"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	misc "github.com/superisaac/jointrpc/misc"
	rpcrouter "github.com/superisaac/jointrpc/rpcrouter"
)

func NewDispatcher() *Dispatcher {
	disp := new(Dispatcher)
	disp.ChResult = make(chan ResultT, 100)
	disp.MethodHandlers = make(map[string](MethodHandler))
	disp.changeHandlers = make([]OnChangeFunc, 0)
	return disp
}

func (self *Dispatcher) On(method string, handler HandlerFunc, opts ...func(*MethodHandler)) {
	if !jsonrpc.IsMethod(method) {
		panic(errors.New("invalid method name"))
	}
	h := MethodHandler{function: handler}
	for _, opt := range opts {
		opt(&h)
	}

	_, found := self.MethodHandlers[method]
	self.MethodHandlers[method] = h

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
	_, found := self.MethodHandlers[method]
	if found {
		delete(self.MethodHandlers, method)
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
		//self.ReturnResultMessage(errmsg)
		return errmsg, nil
		//return , err
	} else if msg.IsRequest() {
		log.Debugf("msg is request")
		if resMsg, ok := res.(jsonrpc.IMessage); ok {
			// TODO: assert resMsg is res and resId matches
			return resMsg, nil
		}
		return jsonrpc.NewResultMessage(msg, res, nil), nil
	} else {
		return nil, nil
	}
}

func (self *Dispatcher) ReturnResultMessage(resmsg jsonrpc.IMessage, req rpcrouter.MsgVec) {
	self.ChResult <- ResultT{
		ResMsg:    resmsg,
		ReqMsgVec: req,
	}
}

func (self *Dispatcher) HandleRequestMessage(msgvec rpcrouter.MsgVec) {
	if self.spawnExec {
		go self.handleRequestMessage(msgvec)
	} else {
		self.handleRequestMessage(msgvec)
	}
}

func (self *Dispatcher) handleRequestMessage(msgvec rpcrouter.MsgVec) {
	msg := msgvec.Msg
	namespace := msgvec.Namespace
	misc.Assert(namespace != "", "empty namespace")

	handler, ok := self.MethodHandlers[msg.MustMethod()]

	defer func() {
		if r := recover(); r != nil {

			if r == Deferred {
				log.Infof("handler is deferred")
				return
			} else if rpcError, ok := r.(*jsonrpc.RPCError); ok {
				errmsg := rpcError.ToMessage(msg)
				self.ReturnResultMessage(errmsg, msgvec)
				return
			} else {
				log.Errorf("Recovered ERROR on handling request msg %+v", r)
				errmsg := jsonrpc.ErrServerError.ToMessage(msg)
				self.ReturnResultMessage(errmsg, msgvec)
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
		resmsg, err = jsonrpc.ErrNoSuchMethod.ToMessage(msg), nil
	}

	//log.Debugf("handle request method %+v, resmsg %+v, error %+v", msg, resmsg, err)
	if err == Deferred {
		log.Infof("handler is deferred")
		return
	}
	if err != nil {
		log.Warnf("bad up message %w", err)
		errMsg := jsonrpc.ErrBadResource.ToMessage(msg)
		self.ReturnResultMessage(errMsg, msgvec)
		return
	}
	if resmsg != nil {
		self.ReturnResultMessage(resmsg, msgvec)
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
