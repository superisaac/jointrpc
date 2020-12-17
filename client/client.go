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
	sendUpChannel := make(chan *intf.JSONRPCUpPacket)
	c := &RPCClient{ServerAddress: serverAddress, sendUpChannel: sendUpChannel}
	c.InitHandlerManager()
	c.OnChange(func() {
		c.OnHandlerChanged()
	})
	return c
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

// Override Handler.OnHandlerChanged
func (self *RPCClient) OnHandlerChanged() {
	if self.TubeClient != nil {
		self.updateMethods()
	}
}

func (self *RPCClient) updateMethods() {
	upMethods := make([](*intf.MethodInfo), 0)
	for m, info := range self.MethodHandlers {
		minfo := &intf.MethodInfo{Name: m, Help: info.Help}
		upMethods = append(upMethods, minfo)
	}
	up := &intf.UpdateMethodsRequest{Methods: upMethods}
	payload := &intf.JSONRPCUpPacket_UpdateMethods{UpdateMethods: up}
	uppac := &intf.JSONRPCUpPacket{Payload: payload}
	self.sendUpChannel <- uppac
}

func (self *RPCClient) RunHandlers() error {
	for {
		err := self.handleRPC()
		//log.Debugf("handle rpc %v", err)
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
		case resmsg, ok := <-self.ChResultMsg:
			if !ok {
				log.Warnf("result msg closed")
				return
			}
			rst, err := server.MessageToResult(resmsg)
			if err != nil {
				panic(err)
			}
			payload := &intf.JSONRPCUpPacket_Result{Result: rst}
			uppac := &intf.JSONRPCUpPacket{Payload: payload}
			stream.Send(uppac)
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
		log.Warnf("error on handle %v", err)
		return err
	}
	log.Debugf("connected")

	sendCtx, sendCancel := context.WithCancel(context.Background())
	defer sendCancel()

	go self.sendUpResult(sendCtx, stream)
	self.updateMethods()

	for {
		downpac, err := stream.Recv()
		if err == io.EOF {
			log.Infof("client stream closed")
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

func (self *RPCClient) handleDownRequest(req *intf.JSONRPCRequest) {
	msg, err := server.RequestToMessage(req)
	if err != nil {
		log.Warnf("parse request message error %+v", err)
		errmsg := jsonrpc.NewErrorMessage(req.Id, 10400, "parse message error", false)
		self.ReturnResultMessage(errmsg)
		return
	}
	self.HandleRequestMessage(msg)
}
