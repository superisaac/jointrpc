package client

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"net"
	"reflect"

	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"

	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/msgutil"
	"github.com/superisaac/jointrpc/rpcrouter"
	"io"
	//grpc "google.golang.org/grpc"
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
	for _, minfo := range disp.GetPublicMethodInfos() {
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

	return self.LiveCall(rootCtx, reqmsg, func(res jsonrpc.IMessage) {
		log.Infof("declared methods %+v", upMethods)
	})
}

func (self *RPCClient) sendPing(ctx context.Context) {
	reqId := misc.NewUuid()
	ping := jsonrpc.NewRequestMessage(reqId, "_stream.ping", nil)
	self.LiveCall(ctx, ping, func(res jsonrpc.IMessage) {
		log.Debugf("pong received, %s", res.MustResult())
	})
}

func (self *RPCClient) Live(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	disp.OnChange(func() {
		self.OnHandlerChanged(disp)
	})

	self.retry = 0
	for {
		if self.retry >= self.LiveRetryTimes {
			break
		}

		log.Debugf("Live connect %d times", self.retry)
		var err error
		if self.IsHttp() {
			err = self.runHTTPLiveStream(rootCtx, disp)
		} else {
			misc.Assert(self.IsH2(), "rpc client is not via grpc")
			err = self.runGRPCLiveStream(rootCtx, disp)
		}

		self.connected = false
		self.retry++
		log.Infof("live connect failed %d times, retry", self.retry)

		if self.onConnectionLost != nil {
			self.onConnectionLost()
		}

		if err != nil {
			return err
		}
		// wait to retry
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (self *RPCClient) sendUpGRPC(ctx context.Context, stream intf.JointRPC_LiveClient, disp *dispatch.Dispatcher) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(15 * time.Second):
			self.sendPing(ctx)
		case <-time.After(5 * time.Second):
			self.cleanTimeoutLivecalls()
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
		case <- ctx.Done():
			return
		case <- time.After(15 * time.Second):
			self.sendPing(ctx)
		case <- time.After(5 * time.Second):
			self.cleanTimeoutLivecalls()
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

func (self *RPCClient) OnAuthorized(cb AuthorizedCallback) {
	self.onAuthorized = cb
}

func (self *RPCClient) NewAuthRequest() *jsonrpc.RequestMessage {
	reqId := misc.NewUuid()
	auth := self.ClientAuth()
	params := [](interface{}){auth.Username, auth.Password}
	return jsonrpc.NewRequestMessage(reqId, "_stream.authorize", params)
}

func (self *RPCClient) NewWatchStateRequest() *jsonrpc.RequestMessage {
	reqId := misc.NewUuid()
	params := [](interface{}){}
	return jsonrpc.NewRequestMessage(reqId, "_stream.watchState", params)
}

func (self *RPCClient) handleDownRequest(ctx context.Context, msg jsonrpc.IMessage, disp *dispatch.Dispatcher, namespace string) {
	cmdMsg := rpcrouter.CmdMsg{
		Msg:       msg,
		Namespace: namespace}
	disp.Feed(ctx, cmdMsg, self.chResult)
}

// transport specific lives
func (self *RPCClient) runHTTPLiveStream(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	if self.connected {
		return errors.New("live stream already connected")
	}

	ws, _, err := websocket.DefaultDialer.Dial(self.WebsocketUrlString(), nil)
	if err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) {
			log.Infof("close failed %s", opErr)
			return nil
		}

		log.Warnf("error on dailing websocket type %s, %s",
			reflect.TypeOf(err), err)
		return errors.Wrap(err, "websocket error")
	}

	defer ws.Close()

	self.connected = true
	self.retry = 0
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
	if err != nil {
		var closeErr *websocket.CloseError
		if errors.Is(err, io.EOF) {
			log.Infof("websocket conn failed")
			return nil
		} else if errors.As(err, &closeErr) {
			log.Infof("websocket close error %d %s", closeErr.Code, closeErr.Text)
			return nil

		}
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
	misc.Assert(authRes.IsResult(), fmt.Sprintf("authres is not request %+v", authRes))

	namespace, ok := authRes.MustResult().(string)
	misc.Assert(ok, "authres.result is not string")

	// startup sendup goroutine
	sendCtx, sendCancel := context.WithCancel(ctx)
	defer sendCancel()
	go self.sendUpWS(sendCtx, ws, disp)
	disp.TriggerChange()

	for {
		msg, err := msgutil.WSRecv(ws)
		if err != nil {
			var closeErr *websocket.CloseError
			if errors.Is(err, io.EOF) {
				log.Infof("websocket conn failed")
				return nil
			} else if errors.As(err, &closeErr) {
				log.Infof("websocket close error %d %s", closeErr.Code, closeErr.Text)
				return nil

			}
			return err
		}

		if msg.IsRequestOrNotify() {
			self.handleDownRequest(ctx, msg, disp, namespace)
		} else {
			self.handleLiveResult(msg)
		}
	}
	return nil
}

func (self *RPCClient) runGRPCLiveStream(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	if self.connected {
		return errors.New("live stream already exist")
	}

	stream, err := self.grpcClient.Live(ctx, grpc_retry.WithMax(500))
	if err != nil {
		return msgutil.GRPCHandleCodes(err, codes.Unavailable)
	}

	self.connected = true
	self.retry = 0
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
	if err != nil {
		return msgutil.GRPCHandleCodes(err, codes.Unavailable)
	}

	if authRes.IsError() {
		rpcError := authRes.MustError()
		return &RPCStatusError{
			Method: "_stream.authorize",
			Code:   rpcError.Code,
			Reason: rpcError.Message,
		}
	}

	misc.Assert(authRes.IsResult(), fmt.Sprintf("authres is not request %+v", authRes))

	namespace, ok := authRes.MustResult().(string)
	misc.Assert(ok, "authres.result is not string")

	if self.onAuthorized != nil {
		self.onAuthorized()
	}

	// startup sendup goroutine
	sendCtx, sendCancel := context.WithCancel(rootCtx)
	defer sendCancel()
	go self.sendUpGRPC(sendCtx, stream, disp)
	disp.TriggerChange()
	for {
		msg, err := msgutil.GRPCClientRecv(stream)
		if err != nil {
			return msgutil.GRPCHandleCodes(err, codes.Unavailable)
		}

		if msg.IsRequestOrNotify() {
			self.handleDownRequest(rootCtx, msg, disp, namespace)
		} else {
			self.handleLiveResult(msg)
		}
	}
	return nil
}
