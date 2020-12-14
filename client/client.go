package client

import (
	"context"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	server "github.com/superisaac/rpctube/server"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"io"
	"time"
)

func NewRPCClient(serverAddress string) *RPCClient {
	methodHandlers := make(map[string](MethodHandler))
	sendUpChannel := make(chan *intf.JSONRPCUpPacket)
	return &RPCClient{ServerAddress: serverAddress, methodHandlers: methodHandlers, sendUpChannel: sendUpChannel}
}

func (self *RPCClient) Connect() error {
	conn, err := grpc.Dial(self.ServerAddress,
		grpc.WithInsecure())
	if err != nil {
		return err
	}
	self.TubeClient = intf.NewJSONRPCTubeClient(conn)
	return nil
}

func (self *RPCClient) wrapHandlerResult(msg *jsonrpc.RPCMessage, res interface{}, err error) (*jsonrpc.RPCMessage, error) {
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

func (self *RPCClient) updateMethods() {
	upMethods := make([](*intf.MethodInfo), 0)
	for m := range self.methodHandlers {
		minfo := &intf.MethodInfo{Name: m}
		upMethods = append(upMethods, minfo)
	}
	up := &intf.UpdateMethodsRequest{Methods: upMethods}
	payload := &intf.JSONRPCUpPacket_UpdateMethods{UpdateMethods: up}
	uppac := &intf.JSONRPCUpPacket{Payload: payload}
	self.sendUpChannel <- uppac
}

func (self *RPCClient) returnResult(resmsg *jsonrpc.RPCMessage) {
	rst, err := server.MessageToResult(resmsg)
	if err != nil {
		panic(err)
	}
	payload := &intf.JSONRPCUpPacket_Result{Result: rst}
	uppac := &intf.JSONRPCUpPacket{Payload: payload}

	self.sendUpChannel <- uppac
}

func (self *RPCClient) On(method string, handler HandlerFunc, opts ...func(*MethodHandler)) {
	h := MethodHandler{function: handler, schema: "", concurrent: false}
	for _, opt := range opts {
		opt(&h)
	}

	_, found := self.methodHandlers[method]
	self.methodHandlers[method] = h

	if !found && self.sendUpChannel != nil {
		self.updateMethods()
	}
}

func (self *RPCClient) UnHandle(method string) bool {
	_, found := self.methodHandlers[method]
	if found {
		delete(self.methodHandlers, method)
		if self.sendUpChannel != nil {
			self.updateMethods()
		}
	}
	return found
}

func WithSchema(schema string) func(*MethodHandler) {
	return func(h *MethodHandler) {
		// TODO: parse schema
		h.schema = schema
	}
}

func WithConcurrent(c bool) func(*MethodHandler) {
	return func(h *MethodHandler) {
		// TODO: parse schema
		h.concurrent = c
	}
}

func (self *RPCClient) OnDefault(handler DefaultHandlerFunc, opts ...func(*RPCClient)) {
	self.defaultHandler = handler
	for _, opt := range opts {
		opt(self)
	}
}

func WithDefaultConcurrent(c bool) func(*RPCClient) {
	return func(h *RPCClient) {
		// TODO: parse schema
		h.defaultConcurrent = c
	}
}

func (self *RPCClient) HandleRPC() error {
	for {
		err := self.handleRPC()
		if err != nil {
			if grpc.Code(err) == codes.Unavailable {
				log.Debugf("connect closed retrying")
			} else {
				return err
			}
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (self *RPCClient) sendUpResult(ctx context.Context, stream intf.JSONRPCTube_HandleClient) {
	// self.sendUpChannel = make(chan *intf.JSONRPCUpPacket)
	// defer func() {
	// 	self.sendUpChannel = nil
	// }()

	for {
		select {
		case <-ctx.Done():
			return
		case uppacket, ok := <-self.sendUpChannel:
			if !ok {
				log.Warnf("send up channel closed")
				return
			}
			stream.Send(uppacket)
		}
	}
}

func (self *RPCClient) DeliverUpPacket(uppack *intf.JSONRPCUpPacket) {
	self.sendUpChannel <- uppack
}

func (self *RPCClient) handleRPC() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := self.TubeClient.Handle(ctx, grpc_retry.WithMax(500))
	if err != nil {
		return err
	}

	sendCtx, sendCancel := context.WithCancel(context.Background())
	defer sendCancel()

	go self.sendUpResult(sendCtx, stream)
	self.updateMethods()

	for {
		downpac, err := stream.Recv()
		if err == io.EOF {
			log.Infof("eor close")
			return nil
		} else if err != nil {
			log.Debugf("down pack error %+v %d", err, grpc.Code(err))
			return err
		}

		// On Ping
		ping := downpac.GetPing()
		if ping != nil {
			// Send Pong
			pong := &intf.PONG{Text: ping.Text}
			payload := &intf.JSONRPCUpPacket_Pong{Pong: pong}
			uppac := &intf.JSONRPCUpPacket{Payload: payload}

			//stream.Send(uppac)
			self.sendUpChannel <- uppac
			continue
		}

		// Handle JSONRPC Request
		req := downpac.GetRequest()
		if req != nil {
			if self.CanRunConcurrent(req.Method) {
				go self.handleDownRequest(req)
			} else {
				self.handleDownRequest(req)
			}
			continue
		}
	}
	return nil
}

func (self RPCClient) CanRunConcurrent(method string) bool {
	handler, ok := self.methodHandlers[method]
	if ok {
		return handler.concurrent
	} else if self.defaultConcurrent {
		return true
	}
	return false
}

func (self *RPCClient) handleDownRequest(req *intf.JSONRPCRequest) {
	msg, err := server.RequestToMessage(req)
	if err != nil {
		log.Warnf("parse request message error %w", err)
		errmsg := jsonrpc.NewErrorMessage(req.Id, 10400, "parse message error", false)
		self.returnResult(errmsg)
		return
	}
	self.handleRequestMsg(msg)
}

func (self *RPCClient) handleRequestMsg(msg *jsonrpc.RPCMessage) {
	handler, ok := self.methodHandlers[msg.Method]
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
		self.returnResult(errmsg)
		return
	}
	if r := recover(); r != nil {
		log.Fatalf("error on handling request msg %+v", r)
		return
	}
	if resmsg != nil {
		self.returnResult(resmsg)
	}
}
