package client

import (
	"context"
	"errors"
	"flag"
	"fmt"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	"io"
	"net/url"
	"os"
	"time"
	//server "github.com/superisaac/jointrpc/server"
	"github.com/superisaac/jointrpc/dispatch"
	encoding "github.com/superisaac/jointrpc/encoding"
	"github.com/superisaac/jointrpc/rpcrouter"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	credentials "google.golang.org/grpc/credentials"
)

func (self RPCStatusError) Error() string {
	return fmt.Sprintf("RPC Response status %s %d %s", self.Method, self.Code, self.Reason)
}

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
	//c.disp = new(dispatch.Dispatcher)
	//c.disp.InitDispatcher()
	// c.disp.OnChange(func() {
	// 	c.OnHandlerChanged()
	// })
	return c
}

func (self RPCClient) ClientAuth() *intf.ClientAuth {
	if self.serverUrl.User == nil {
		return &intf.ClientAuth{}
	}
	pwd, ok := self.serverUrl.User.Password()
	if !ok {
		pwd = ""
	}

	auth := &intf.ClientAuth{
		Username: self.serverUrl.User.Username(),
		Password: pwd,
	}
	return auth
}

func (self RPCClient) String() string {
	return self.serverEntry.ServerUrl
}

func (self RPCClient) ServerEntry() ServerEntry {
	return self.serverEntry
}

//func (self RPCClient) AttachTo(disp Dispatcher) {
//	self.disp = disp
//}

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
	self.rpcClient = intf.NewJointRPCClient(conn)
	return nil
}

// Override Handler.OnHandlerChanged
func (self *RPCClient) OnHandlerChanged(disp *dispatch.Dispatcher) {
	if self.rpcClient != nil && self.workerStream != nil {
		self.declareMethods(context.Background(), disp)
	}
}

func (self *RPCClient) declareMethods(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	upMethods := make([](*intf.MethodInfo), 0)
	for m, info := range disp.MethodHandlers {
		minfo := &intf.MethodInfo{Name: m, Help: info.Help, SchemaJson: info.SchemaJson}
		upMethods = append(upMethods, minfo)
	}

	req := &intf.DeclareMethodsRequest{Methods: upMethods}
	payload := &intf.JointRPCUpPacket_MethodsRequest{MethodsRequest: req}
	uppac := &intf.JointRPCUpPacket{Payload: payload}
	self.DeliverUpPacket(uppac)
	return nil
}

func (self *RPCClient) Worker(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	disp.OnChange(func() {
		self.OnHandlerChanged(disp)
	})

	for {
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
		case resmsg, ok := <-disp.ChResult:
			if !ok {
				log.Warnf("result msg closed")
				return
			}

			envo := encoding.MessageToEnvolope(resmsg)
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
	payload := &intf.JointRPCUpPacket_Auth{Auth: self.ClientAuth()}
	uppac := &intf.JointRPCUpPacket{Payload: payload}
	return self.workerStream.Send(uppac)
}

func (self *RPCClient) runWorker(rootCtx context.Context, disp *dispatch.Dispatcher) error {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	if self.workerStream != nil {
		return errors.New("worker stream already exist")
	}

	stream, err := self.rpcClient.Worker(ctx, grpc_retry.WithMax(500))

	if err != nil {
		log.Warnf("error on handle %v", err)
		return err
	}
	self.connected = true
	self.workerStream = stream
	if self.onConnected != nil {
		self.onConnected()
	}

	err = self.requestAuth(rootCtx)
	if err != nil {
		return err
	}

	sendCtx, sendCancel := context.WithCancel(rootCtx)
	defer sendCancel()

	go self.sendUpResult(sendCtx, disp)

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
			pong := &intf.Pong{Text: ping.Text}
			payload := &intf.JointRPCUpPacket_Pong{Pong: pong}
			uppac := &intf.JointRPCUpPacket{Payload: payload}

			self.sendUpChannel <- uppac
			continue
		}

		// subscribed state
		istate := downpac.GetState()
		if istate != nil {
			state := encoding.DecodeServerState(istate)
			disp.TriggerStateChange(state)
			continue
		}

		// methods response
		methodsResp := downpac.GetMethodsResponse()
		if methodsResp != nil {
			err := self.CheckStatus(methodsResp.Status, "Worker.Methods")
			if err != nil {
				log.Warn(err.Error())
				return err
			}
			continue
		}

		// delegates response
		delegatesResp := downpac.GetDelegatesResponse()
		if delegatesResp != nil {
			err := self.CheckStatus(delegatesResp.Status, "Worker.Delegates")
			if err != nil {
				log.Warn(err.Error())
				return err
			}
			continue
		}

		// Set connPublicId
		echo := downpac.GetEcho()
		if echo != nil {
			err := self.CheckStatus(echo.Status, "Worker.Auth")
			if err != nil {
				log.Warn(err.Error())
				return err
			}
			self.connPublicId = echo.ConnPublicId
			log.Infof("Handle() got conn public id %s", self.connPublicId)
			disp.TriggerChange()
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
			self.handleDownRequest(msg, envo.TraceId, disp)
			continue
		}
	}
	return nil
}

func (self *RPCClient) handleDownRequest(msg jsonrpc.IMessage, traceId string, disp *dispatch.Dispatcher) {
	msgvec := rpcrouter.MsgVec{
		Msg:        msg,
		FromConnId: 0}
	disp.HandleRequestMessage(msgvec)
}

func (self *RPCClient) CheckStatus(status *intf.Status, methodName string) error {
	if status == nil || status.Code == 0 {
		return nil
	} else {
		return &RPCStatusError{
			Method: methodName,
			Code:   int(status.Code),
			Reason: status.Reason}
	}
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
