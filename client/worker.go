package client

import (
	"context"
	"errors"
	//"flag"
	//"fmt"
	"github.com/gorilla/websocket"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/mitchellh/mapstructure"

	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"

	"github.com/superisaac/jointrpc/jsonrpc"
	"io"
	//"net/url"
	//"os"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/msgutil"
	"github.com/superisaac/jointrpc/rpcrouter"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"time"
)

// Override Handler.OnHandlerChanged
func (self *RPCClient) OnHandlerChanged(disp *dispatch.Dispatcher) {
	//if self.grpcClient != nil && self.grpcStream != nil {
	if self.connected {
		self.declareMethods(context.Background(), disp)
	}
}

func (self *RPCClient) declareMethods(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	upMethods := make([](map[string](interface{})), 0)
	for _, minfo := range disp.GetMethodInfos() {
		infoDict := make(map[string](interface{}))
		err := mapstructure.Decode(minfo, &infoDict)
		if err != nil {
			return err
		}
		upMethods = append(upMethods, infoDict)
	}

	reqId := misc.NewUuid()
	params := make([]interface{}, 0)
	params = append(params, upMethods)
	reqmsg := jsonrpc.NewRequestMessage(reqId, "_stream.declareMethods", params)

	return self.CallInStream(rootCtx, reqmsg, func(res jsonrpc.IMessage) {
		log.Debugf("declared methods")
	})
}

func (self *RPCClient) Worker(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	disp.OnChange(func() {
		self.OnHandlerChanged(disp)
	})

	for i := 0; i < self.WorkerRetryTimes; i++ {
		log.Debugf("Worker connect %d times", i)
		var err error
		if self.IsHttp() {
			err = self.runHTTPWorker(rootCtx, disp)
		} else {
			misc.Assert(self.IsH2(), "rpc client is not via grpc")
			err = self.runGRPCWorker(rootCtx, disp)
		}

		self.connected = false
		if err != nil {
			return err
		}
		if self.onConnectionLost != nil {
			self.onConnectionLost()
		}
		// wait to retry
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (self *RPCClient) sendUpGRPC(ctx context.Context, stream intf.JointRPC_WorkerClient, disp *dispatch.Dispatcher) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-self.chSendUp:
			if !ok {
				log.Warnf("send up channel closed")
				return
			}
			err := msgutil.GRPCClientSend(stream, msg)
			if err != nil {
				msg.Log().Warnf("send failed, %s", err)
				panic(err)
			}
		case result, ok := <-self.chResult:
			if !ok {
				log.Warnf("result msg closed")
				return
			}
			msg := result.ResMsg
			err := msgutil.GRPCClientSend(stream, msg)
			if err != nil {
				msg.Log().Warnf("send failed, %s", err)
				panic(err)
			}
		}
	}
}

func (self *RPCClient) sendUpWS(ctx context.Context, ws *websocket.Conn, disp *dispatch.Dispatcher) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-self.chSendUp:
			if !ok {
				log.Warnf("send up channel closed")
				return
			}
			err := msgutil.WSSend(ws, msg)
			if err != nil {
				msg.Log().Warnf("send failed, %s", err)
				panic(err)
			}
		case result, ok := <-self.chResult:
			if !ok {
				log.Warnf("result msg closed")
				return
			}
			msg := result.ResMsg
			err := msgutil.WSSend(ws, msg)
			if err != nil {
				msg.Log().Warnf("send failed, %s", err)
				panic(err)
			}
		}
	}
}

func (self *RPCClient) OnConnected(cb ConnectedCallback) {
	self.onConnected = cb
}
func (self *RPCClient) OnConnectionLost(cb ConnectionLostCallback) {
	self.onConnectionLost = cb
}

func (self *RPCClient) NewAuthRequest() *jsonrpc.RequestMessage {
	reqId := misc.NewUuid()
	auth := self.ClientAuth()
	params := [](interface{}){auth.Username, auth.Password}
	return jsonrpc.NewRequestMessage(reqId, "_stream.authorize", params)
}

func (self *RPCClient) handleDownRequest(ctx context.Context, msg jsonrpc.IMessage, traceId string, disp *dispatch.Dispatcher, namespace string) {
	msgvec := rpcrouter.MsgVec{
		Msg:        msg,
		Namespace:  namespace,
		FromConnId: 0}
	disp.Feed(ctx, msgvec, self.chResult)
}

// transport specific workers
func (self *RPCClient) runHTTPWorker(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	if self.connected {
		return errors.New("worker stream already connected")
	}

	ws, _, err := websocket.DefaultDialer.Dial(self.WebsocketUrlString(), nil)
	if err != nil {
		log.Warnf("error on dailing websocket %s", err)
		return err
	}

	defer ws.Close()

	self.connected = true

	if self.onConnected != nil {
		self.onConnected()
	}

	authreq := self.NewAuthRequest()
	err = msgutil.WSSend(ws, authreq)
	if err != nil {
		return err
	}

	// wait for auth response
	authRes, err := msgutil.WSRecv(ws)
	if err == io.EOF {
		log.Infof("websocket conn failed")
		return nil
	} else if err != nil {
		return err
	}

	if authRes.IsError() {
		rpcError := authRes.MustError()
		return &RPCStatusError{
			Method: "_stream.authorize",
			Code:   rpcError.Code,
			Reason: rpcError.Message,
		}
	}
	misc.Assert(authRes.IsResult(), "authres is not request")

	namespace, ok := authRes.MustResult().(string)
	misc.Assert(ok, "authres.result is not string")

	// startup sendup goroutine
	sendCtx, sendCancel := context.WithCancel(ctx)
	defer sendCancel()
	go self.sendUpWS(sendCtx, ws, disp)
	disp.TriggerChange()
	for {
		msg, err := msgutil.WSRecv(ws)
		if err == io.EOF {
			log.Infof("websocket conn failed")
			return nil
		} else if err != nil {
			return err
		}

		if msg.IsRequestOrNotify() {
			self.handleDownRequest(ctx, msg, msg.TraceId(), disp, namespace)
		} else {
			self.handleWireResult(msg)
		}
	}
	return nil
}

func (self *RPCClient) runGRPCWorker(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	if self.connected {
		return errors.New("worker stream already exist")
	}

	stream, err := self.grpcClient.Worker(ctx, grpc_retry.WithMax(500))
	if err == io.EOF {
		log.Infof("cannot connect stream")
		return nil
	} else if grpc.Code(err) == codes.Unavailable {
		log.Debugf("connect closed retrying")
		return nil
	} else if err != nil {
		log.Warnf("error on handle %v", err)
		return err
	}

	self.connected = true
	if self.onConnected != nil {
		self.onConnected()
	}

	authmsg := self.NewAuthRequest()
	err = msgutil.GRPCClientSend(stream, authmsg)
	if err != nil {
		return err
	}

	// wait for auth response
	authRes, err := msgutil.GRPCClientRecv(stream)
	if err == io.EOF {
		log.Infof("client stream closed")
		return nil
	} else if grpc.Code(err) == codes.Unavailable {
		log.Debugf("connect closed retrying")
		return nil
	} else if err != nil {
		log.Debugf("down pack error code=%d, %+v", grpc.Code(err), err)
		return err
	}

	if authRes.IsError() {
		rpcError := authRes.MustError()
		return &RPCStatusError{
			Method: "_stream.authorize",
			Code:   rpcError.Code,
			Reason: rpcError.Message,
		}
	}
	misc.Assert(authRes.IsResult(), "authres is not request")

	namespace, ok := authRes.MustResult().(string)
	misc.Assert(ok, "authres.result is not string")

	// startup sendup goroutine
	sendCtx, sendCancel := context.WithCancel(rootCtx)
	defer sendCancel()
	go self.sendUpGRPC(sendCtx, stream, disp)
	disp.TriggerChange()
	for {
		msg, err := msgutil.GRPCClientRecv(stream)
		if err == io.EOF {
			log.Infof("client stream closed")
			return nil
		} else if grpc.Code(err) == codes.Unavailable {
			log.Debugf("connect closed retrying")
			return nil
		} else if err != nil {
			log.Debugf("down pack error code=%d, %+v", grpc.Code(err), err)
			return err
		}

		if msg.IsRequestOrNotify() {
			self.handleDownRequest(rootCtx, msg, msg.TraceId(), disp, namespace)
		} else {
			self.handleWireResult(msg)
		}
	}
	return nil
}
