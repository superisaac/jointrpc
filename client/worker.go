package client

import (
	"context"
	"errors"
	//"flag"
	//"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	//intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/jsonrpc"
	"io"
	//"net/url"
	//"os"
	"github.com/superisaac/jointrpc/dispatch"
	encoding "github.com/superisaac/jointrpc/encoding"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"time"
)

// Override Handler.OnHandlerChanged
func (self *RPCClient) OnHandlerChanged(disp *dispatch.Dispatcher) {
	if self.grpcClient != nil && self.workerStream != nil {
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

	return self.CallInWire(rootCtx, reqmsg, func(res jsonrpc.IMessage) {
		log.Debugf("declared methods")
	})
}

func (self *RPCClient) Worker(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	disp.OnChange(func() {
		self.OnHandlerChanged(disp)
	})

	for i := 0; i < self.WorkerRetryTimes; i++ {
		log.Debugf("Worker connect %d times", i)
		err := self.runWorker(rootCtx, disp)
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

func (self *RPCClient) sendUpResult(ctx context.Context, disp *dispatch.Dispatcher) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-self.chSendUp:
			if !ok {
				log.Warnf("send up channel closed")
				return
			}
			envo := encoding.MessageToEnvolope(msg)
			self.workerStream.Send(envo)
		case result, ok := <-self.chResult:
			if !ok {
				log.Warnf("result msg closed")
				return
			}

			envo := encoding.MessageToEnvolope(result.ResMsg)
			self.workerStream.Send(envo)
		}
	}
}

func (self *RPCClient) OnConnected(cb ConnectedCallback) {
	self.onConnected = cb
}
func (self *RPCClient) OnConnectionLost(cb ConnectionLostCallback) {
	self.onConnectionLost = cb
}

func (self *RPCClient) requestAuth(rootCtx context.Context) error {
	reqId := misc.NewUuid()
	auth := self.ClientAuth()
	params := [](interface{}){auth.Username, auth.Password}
	authmsg := jsonrpc.NewRequestMessage(reqId, "_stream.authorize", params)
	envo := encoding.MessageToEnvolope(authmsg)
	return self.workerStream.Send(envo)
}

func (self *RPCClient) runWorker(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	if self.workerStream != nil {
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

	self.workerStream = stream

	defer func() {
		self.workerStream = nil
	}()

	self.connected = true
	if self.onConnected != nil {
		self.onConnected()
	}

	err = self.requestAuth(rootCtx)
	if err != nil {
		return err
	}

	// wait for auth response
	authRespEnvo, err := self.workerStream.Recv()
	if err == io.EOF {
		log.Infof("client stream closed")
		return nil
	} else if grpc.Code(err) == codes.Unavailable {
		log.Debugf("connect closed retrying")
		return nil
	} else if err != nil {
		log.Debugf("down pack error %+v %d", err, grpc.Code(err))
		return err
	}

	authRes, err := encoding.MessageFromEnvolope(authRespEnvo)
	if err != nil {
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
	go self.sendUpResult(sendCtx, disp)
	disp.TriggerChange()
	for {
		envo, err := self.workerStream.Recv()
		if err == io.EOF {
			log.Infof("client stream closed")
			return nil
		} else if grpc.Code(err) == codes.Unavailable {
			log.Debugf("connect closed retrying")
			return nil
		} else if err != nil {
			log.Debugf("down pack error %+v %d", err, grpc.Code(err))
			return err
		}

		msg, err := encoding.MessageFromEnvolope(envo)
		if err != nil {
			return err
		}
		if msg.IsRequestOrNotify() {
			self.handleDownRequest(rootCtx, msg, envo.TraceId, disp, namespace)
		} else {
			self.handleWireResult(msg)
		}
	}
	return nil
}

func (self *RPCClient) handleDownRequest(ctx context.Context, msg jsonrpc.IMessage, traceId string, disp *dispatch.Dispatcher, namespace string) {
	msgvec := rpcrouter.MsgVec{
		Msg:        msg,
		Namespace:  namespace,
		FromConnId: 0}
	disp.Feed(ctx, msgvec, self.chResult)
}
