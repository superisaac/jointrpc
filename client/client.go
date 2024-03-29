package client

import (
	//"context"
	"fmt"
	"github.com/pkg/errors"
	//grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jsonz"
	//"io"
	"net/url"
	//"time"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/misc"
	//"github.com/superisaac/jointrpc/rpcrouter"
	grpc "google.golang.org/grpc"
	//codes "google.golang.org/grpc/codes"
	credentials "google.golang.org/grpc/credentials"
)

func (self RPCStatusError) Error() string {
	return fmt.Sprintf("RPC Response status %s %d %s", self.Method, self.Code, self.Reason)
}

func NewRPCClient(serverEntry ServerEntry) *RPCClient {
	chSendUp := make(chan jsonz.Message, misc.DefaultChanSize())
	serverUrl, err := url.Parse(serverEntry.ServerUrl)
	if err != nil {
		log.Panicf("parse url error %s %+v", serverEntry.ServerUrl, errors.Wrap(err, "url.Parse()"))
	}

	scm := serverUrl.Scheme
	if !(scm == "h2" || scm == "h2c" || scm == "http" || scm == "https") {
		log.Panicf("urls scheme not allowed, %s", serverUrl)
	}

	chResult := make(chan dispatch.ResultT, misc.DefaultChanSize())

	c := &RPCClient{
		serverEntry:    serverEntry,
		serverUrl:      serverUrl,
		chSendUp:       chSendUp,
		LiveRetryTimes: 10,
		chResult:       chResult,
		//wirePendingRequests: make(map[interface{}]WireCallT),
	}
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

func (self RPCClient) WebsocketUrlString() string {
	var scheme = "ws"
	if self.serverUrl.Scheme == "https" {
		scheme = "wss"
	}
	return fmt.Sprintf("%s://%s/ws", scheme, self.serverUrl.Host)
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

func (self RPCClient) IsSecure() bool {
	return self.serverUrl.Scheme == "https" || self.serverUrl.Scheme == "h2"
}

func (self RPCClient) Connected() bool {
	return self.connected
}

func (self RPCClient) certFileFromFragment(serverUrl *url.URL) string {
	if serverUrl.Fragment != "" {
		v, err := url.ParseQuery(serverUrl.Fragment)
		if err != nil {
			log.Warnf("server url fragment parse error %s %+v", serverUrl.Fragment, errors.Wrap(err, "url.ParseQuery()"))
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
		return errors.Wrap(err, "grpc.Dial()")
	}
	self.grpcClient = intf.NewJointRPCClient(conn)
	return nil
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
