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
	serverUrl, err := url.Parse(serverEntry.ServerUrl)
	if err != nil {
		log.Panicf("parse url error %s %s", serverEntry.ServerUrl, err.Error())
	}

	scm := serverUrl.Scheme
	if !(scm == "h2" || scm == "h2c" || scm == "http" || scm == "https") {
		log.Panicf("urls scheme not allowed, %s", serverUrl)
	}

	c := &RPCClient{
		serverEntry:   serverEntry,
		serverUrl:     serverUrl,
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

func (self RPCClient) ServerEntry() ServerEntry {
	return self.serverEntry
}

func (self RPCClient) IsHttp() bool {
	return self.serverUrl.Scheme == "http" || self.serverUrl.Scheme == "https"
}

func (self RPCClient) IsH2() bool {
	return self.serverUrl.Scheme == "h2" || self.serverUrl.Scheme == "h2c"
}

func (self RPCClient) ConnPublicId() string {
	return self.connPublicId
}

func (self RPCClient) Connected() bool {
	return self.connected
}

func (self RPCClient) certFileFromFragment(serverUrl *url.URL) string {
	if serverUrl.Fragment != "" {
		v, err := url.ParseQuery(serverUrl.Fragment)
		if err != nil {
			log.Warnf("server url fragment parse error %s %+v", serverUrl.Fragment, err)
		} else {
			return v.Get("cert")
		}
	}
	return ""
}

func (self *RPCClient) Connect() error {
	var opts []grpc.DialOption

	if self.IsHttp() {
		// http method does nothing
		return nil
	} else if self.serverUrl.Scheme == "h2c" {
		opts = append(opts, grpc.WithInsecure())
	} else if self.serverUrl.Scheme == "h2" {
		certFile := self.certFileFromFragment(self.serverUrl)
		if certFile == "" {
			certFile = self.serverEntry.CertFile
		}
		if certFile != "" {
			creds, err := credentials.NewClientTLSFromFile(certFile, "")
			if err != nil {
				panic(err)
			}
			opts = append(opts, grpc.WithTransportCredentials(creds))
		}
	} else {
		log.Panicf("invalid server url scheme %s", self.serverUrl.Scheme)
	}
	conn, err := grpc.Dial(self.serverUrl.Host, opts...)
	if err != nil {
		return err
	}
	self.tubeClient = intf.NewJointRPCClient(conn)
	return nil
}

// Override Handler.OnHandlerChanged
func (self *RPCClient) OnHandlerChanged() {
	if self.tubeClient != nil && self.connPublicId != "" {
		self.declareMethods(context.Background())
	}
}

func (self *RPCClient) declareMethods(rootCtx context.Context) error {
	upMethods := make([](*intf.MethodInfo), 0)
	for m, info := range self.MethodHandlers {
		minfo := &intf.MethodInfo{Name: m, Help: info.Help, SchemaJson: info.SchemaJson}
		upMethods = append(upMethods, minfo)
	}
	log.Infof("declare methods %+v, %s", upMethods, self.connPublicId)
	return self.DeclareMethods(rootCtx, upMethods)
}

func (self *RPCClient) Handle(rootCtx context.Context) error {
	for {
		err := self.handleRPC(rootCtx)
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

func (self *RPCClient) sendUpResult(ctx context.Context, stream intf.JointRPC_HandleClient) {
	for {
		select {
		case <-ctx.Done():
			//stream.
			return
		case uppacket, ok := <-self.sendUpChannel:
			if !ok {
				log.Warnf("send up channel closed")
				return
			}
			stream.Send(uppacket)
		case resmsg, ok := <-self.ChResult:
			if !ok {
				log.Warnf("result msg closed")
				return
			}

			envo := encoding.MessageToEnvolope(resmsg)
			payload := &intf.JointRPCUpPacket_Envolope{Envolope: envo}
			uppac := &intf.JointRPCUpPacket{Payload: payload}
			stream.Send(uppac)
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

func (self *RPCClient) handleRPC(rootCtx context.Context) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	stream, err := self.tubeClient.Handle(ctx, grpc_retry.WithMax(500))

	if err != nil {
		log.Warnf("error on handle %v", err)
		return err
	}
	log.Debugf("connected")
	self.connected = true
	if self.onConnected != nil {
		self.onConnected()
	}

	sendCtx, sendCancel := context.WithCancel(rootCtx)
	defer sendCancel()

	go self.sendUpResult(sendCtx, stream)

	for {
		downpac, err := stream.Recv()
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
			pong := &intf.PONG{Text: ping.Text}
			payload := &intf.JointRPCUpPacket_Pong{Pong: pong}
			uppac := &intf.JointRPCUpPacket{Payload: payload}

			//stream.Send(uppac)
			self.sendUpChannel <- uppac
			continue
		}

		// subscribed state
		istate := downpac.GetState()
		if istate != nil {
			state := encoding.DecodeServerState(istate)
			if self.StateHandler != nil {
				self.StateHandler(state)
			}
			continue
		}

		// Set connPublicId
		greeting := downpac.GetGreeting()
		if greeting != nil {
			self.connPublicId = greeting.ConnPublicId
			log.Infof("Handle() got conn public id %s", self.connPublicId)
			self.TriggerChange()
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
			if !msg.IsRequestOrNotify() {
				log.Warnf("msg is none of reques|notify %+v ", msg)
				continue
			}
			self.handleDownRequest(msg, envo.TraceId)
			continue
		}
	}
	return nil
}

func (self *RPCClient) handleDownRequest(msg jsonrpc.IMessage, traceId string) {
	msgvec := rpcrouter.MsgVec{
		Msg:        msg,
		FromConnId: 0}
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
		value.ServerUrl = "h2c://localhost:50055"
	}

	if value.CertFile == "" {
		value.CertFile = os.Getenv("TUBE_CONNECT")
	}
	return value
}
