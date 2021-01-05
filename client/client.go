package client

import (
	//"fmt"
	"context"
	"flag"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"io"
	"net/url"
	"os"
	"time"
	//server "github.com/superisaac/jointrpc/server"
	encoding "github.com/superisaac/jointrpc/encoding"
	"github.com/superisaac/jointrpc/rpcrouter"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	credentials "google.golang.org/grpc/credentials"
)

func NewRPCClient(serverEntry ServerEntry) *RPCClient {
	sendUpChannel := make(chan *intf.JointRPCUpPacket)
	c := &RPCClient{
		serverEntry:   serverEntry,
		sendUpChannel: sendUpChannel,
	}
	c.InitHandlerManager()
	c.OnChange(func() {
		c.OnHandlerChanged()
	})
	return c
}

func (self RPCClient) String() string {
	return self.serverEntry.ServerUrl
}

func (self *RPCClient) Connect() error {
	var opts []grpc.DialOption

	serverUrl, err := url.Parse(self.serverEntry.ServerUrl)
	if err != nil {
		log.Panicf("error parse server entry url %s, %+v", self.serverEntry.ServerUrl, err)
	}

	if serverUrl.Scheme == "h2c" {
		opts = append(opts, grpc.WithInsecure())
	} else if serverUrl.Scheme == "h2" {
		if self.serverEntry.CertFile != "" {
			creds, err := credentials.NewClientTLSFromFile(self.serverEntry.CertFile, "")
			if err != nil {
				panic(err)
			}
			opts = append(opts, grpc.WithTransportCredentials(creds))
		}
	} else {
		log.Panicf("invalid server url scheme %s", serverUrl.Scheme)
	}
	conn, err := grpc.Dial(serverUrl.Host, opts...)
	if err != nil {
		return err
	}
	self.tubeClient = intf.NewJointRPCClient(conn)
	return nil
}

// Override Handler.OnHandlerChanged
func (self *RPCClient) OnHandlerChanged() {
	if self.tubeClient != nil {
		self.updateMethods()
	}
}

func (self *RPCClient) updateMethods() {
	upMethods := make([](*intf.MethodInfo), 0)
	for m, info := range self.MethodHandlers {
		minfo := &intf.MethodInfo{Name: m, Help: info.Help, SchemaJson: info.SchemaJson}
		upMethods = append(upMethods, minfo)
	}
	up := &intf.CanServeRequest{Methods: upMethods}
	payload := &intf.JointRPCUpPacket_CanServe{CanServe: up}
	uppac := &intf.JointRPCUpPacket{Payload: payload}
	self.sendUpChannel <- uppac
}

func (self *RPCClient) UpdateDelegateMethods(methods []string) {
	up := &intf.CanDelegateRequest{Methods: methods}
	payload := &intf.JointRPCUpPacket_CanDelegate{CanDelegate: up}
	uppac := &intf.JointRPCUpPacket{Payload: payload}
	self.sendUpChannel <- uppac
}

func (self *RPCClient) Handle(rootCtx context.Context) error {
	for {
		err := self.handleRPC(rootCtx)
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

func (self *RPCClient) sendUpResult(ctx context.Context, stream intf.JointRPC_HandleClient) {
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

			envo := &intf.JSONRPCEnvolope{Body: resmsg.MustString()}
			payload := &intf.JointRPCUpPacket_Envolope{Envolope: envo}
			uppac := &intf.JointRPCUpPacket{Payload: payload}
			stream.Send(uppac)
		}
	}
}

func (self *RPCClient) DeliverUpPacket(uppack *intf.JointRPCUpPacket) {
	self.sendUpChannel <- uppack
}

func (self *RPCClient) handleRPC(rootCtx context.Context) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	stream, err := self.tubeClient.Handle(ctx, grpc_retry.WithMax(500))

	if err != nil {
		log.Warnf("error on handle %v", err)
		return err
	}
	log.Debugf("connected")

	sendCtx, sendCancel := context.WithCancel(rootCtx)
	defer sendCancel()

	go self.sendUpResult(sendCtx, stream)
	if len(self.MethodHandlers) > 0 {
		self.updateMethods()
	}

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
			payload := &intf.JointRPCUpPacket_Pong{Pong: pong}
			uppac := &intf.JointRPCUpPacket{Payload: payload}

			//stream.Send(uppac)
			self.sendUpChannel <- uppac
			continue
		}

		istate := downpac.GetState()
		if istate != nil {
			state := encoding.DecodeTubeState(istate)
			if self.StateHandler != nil {
				self.StateHandler(state)
			}
		}

		// Handle JSONRPC Request
		//req := downpac.GetRequest()
		envo := downpac.GetEnvolope()
		if envo != nil {
			msg, err := jsonrpc.ParseBytes([]byte(envo.Body))
			if err != nil {
				return err
			}
			if !msg.IsRequestOrNotify() {
				log.Warnf("msg is none of reques|notify %+v ", msg)
				continue
			}
			if self.CanRunConcurrent(msg.MustMethod()) {
				go self.handleDownRequest(msg)
			} else {
				self.handleDownRequest(msg)
			}
			continue
		}
	}
	return nil
}

func (self *RPCClient) handleDownRequest(msg jsonrpc.IMessage) {
	msgvec := rpcrouter.MsgVec{Msg: msg, FromConnId: 0}
	self.HandleRequestMessage(msgvec)
}

// misc functions

func NewServerFlag(flagSet *flag.FlagSet) *ServerFlag {
	seflag := new(ServerFlag)
	seflag.pAddress = flagSet.String("c", "", "the tube server address")
	seflag.pCertFile = flagSet.String("cert", "", "the cert file, default empty")
	return seflag
}

func (self *ServerFlag) ptrValue() ServerEntry {
	return ServerEntry{
		ServerUrl: *self.pAddress,
		CertFile:  *self.pCertFile,
	}
}

func (self *ServerFlag) Get() ServerEntry {
	value := self.ptrValue()
	if value.ServerUrl == "" {
		value.ServerUrl = os.Getenv("TUBE_CONNECT")
	}

	if value.ServerUrl == "" {
		value.ServerUrl = "localhost:50055"
	}

	if value.CertFile == "" {
		value.CertFile = os.Getenv("TUBE_CONNECT")
	}
	return value
}
