package client

import (
	"context"
	"errors"
	//"flag"
	//"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/jsonrpc"
	"io"
	//"net/url"
	//"os"
	"time"
	//server "github.com/superisaac/jointrpc/server"
	"github.com/superisaac/jointrpc/dispatch"
	encoding "github.com/superisaac/jointrpc/encoding"
	"github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	//credentials "google.golang.org/grpc/credentials"
)

// Override Handler.OnHandlerChanged
func (self *RPCClient) OnHandlerChanged(disp *dispatch.Dispatcher) {
	if self.grpcClient != nil && self.workerStream != nil {
		self.declareMethods(context.Background(), disp)
	}
}

func (self *RPCClient) declareMethods(rootCtx context.Context, disp *dispatch.Dispatcher) {
	upMethods := make([](map[string](interface{})), 0)
	for m, info := range disp.MethodHandlers {
		minfo := map[string](interface{}){
			"name":   m,
			"help":   info.Help,
			"schema": info.SchemaJson,
		}
		upMethods = append(upMethods, minfo)
	}

	reqId := misc.NewUuid()
	params := make([]interface{}, 0)
	params = append(params, upMethods)
	reqmsg := jsonrpc.NewRequestMessage(reqId, "_conn.declareMethods", params, nil)

	self.CallInWire(rootCtx, reqmsg, func(res jsonrpc.IMessage) {
		log.Debugf("declared methods")
	})
}

func (self *RPCClient) Worker(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	disp.OnChange(func() {
		self.OnHandlerChanged(disp)
	})

	//for {
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
		case uppacket, ok := <-self.sendUpChannel:
			if !ok {
				log.Warnf("send up channel closed")
				return
			}
			self.workerStream.Send(uppacket)
		case result, ok := <-disp.ChResult:
			if !ok {
				log.Warnf("result msg closed")
				return
			}

			envo := encoding.MessageToEnvolope(result.ResMsg)
			payload := &intf.JointRPCUpPacket_Envolope{Envolope: envo}
			uppac := &intf.JointRPCUpPacket{Payload: payload}
			self.workerStream.Send(uppac)
		}
	}
}

func (self *RPCClient) OnConnected(cb ConnectedCallback) {
	self.onConnected = cb
}
func (self *RPCClient) OnConnectionLost(cb ConnectionLostCallback) {
	self.onConnectionLost = cb
}

func (self *RPCClient) DeliverUpPacket(uppack *intf.JointRPCUpPacket) {
	self.sendUpChannel <- uppack
}

func (self *RPCClient) requestAuth(rootCtx context.Context) error {
	authReq := &intf.AuthRequest{
		RequestId:  misc.NewUuid(),
		ClientAuth: self.ClientAuth(),
	}
	payload := &intf.JointRPCUpPacket_AuthRequest{AuthRequest: authReq}
	uppac := &intf.JointRPCUpPacket{Payload: payload}
	return self.workerStream.Send(uppac)
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
	downpac, err := self.workerStream.Recv()
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
	authResp := downpac.GetAuthResponse()
	if authResp == nil {
		log.Warnf("cannot wait for auth response")
		return err
	}

	err = self.CheckStatus(authResp.Status, "Worker.Auth")
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	namespace := authResp.Namespace

	// startup sendup goroutine
	sendCtx, sendCancel := context.WithCancel(rootCtx)
	defer sendCancel()
	go self.sendUpResult(sendCtx, disp)
	disp.TriggerChange()
	for {
		downpac, err := self.workerStream.Recv()
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

		// On Ping
		ping := downpac.GetPing()
		if ping != nil {
			// Send Pong
			pong := &intf.Pong{RequestId: ping.RequestId}
			payload := &intf.JointRPCUpPacket_Pong{Pong: pong}
			uppac := &intf.JointRPCUpPacket{Payload: payload}

			self.sendUpChannel <- uppac
			continue
		}

		// Handle JSONRPC Request
		//req := downpac.GetRequest()
		envo := downpac.GetEnvolope()
		if envo != nil {
			msg, err := encoding.MessageFromEnvolope(envo)
			if err != nil {
				return err
			}
			if msg.IsRequestOrNotify() {
				self.handleDownRequest(msg, envo.TraceId, disp, namespace)
			} else {
				self.handleWireResult(msg)
			}

			continue
		}
	}
	return nil
}

func (self *RPCClient) handleDownRequest(msg jsonrpc.IMessage, traceId string, disp *dispatch.Dispatcher, namespace string) {
	msgvec := rpcrouter.MsgVec{
		Msg:        msg,
		Namespace:  namespace,
		FromConnId: 0}
	disp.HandleRequestMessage(msgvec)
}
