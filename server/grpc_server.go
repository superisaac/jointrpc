package server

import (
	"context"
	"github.com/pkg/errors"
	//"io"
	"net"
	//"strings"
	"time"
	//"time"
	//json "encoding/json"
	//"errors"
	//"fmt"
	//"log"
	//simplejson "github.com/bitly/go-simplejson"
	//grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"

	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	misc "github.com/superisaac/jointrpc/misc"
	msgutil "github.com/superisaac/jointrpc/msgutil"
	jsonrpc "github.com/superisaac/jsonrpc"
	//misc "github.com/superisaac/jointrpc/misc"
	//datadir "github.com/superisaac/jointrpc/datadir"
	//schema "github.com/superisaac/jsonrpc/schema"
	"github.com/superisaac/jointrpc/dispatch"
	"github.com/superisaac/jointrpc/rpcrouter"
	peer "google.golang.org/grpc/peer"
)

type JointRPC struct {
	intf.UnimplementedJointRPCServer
}

func (self JointRPC) Authorize(context context.Context, auth *intf.ClientAuth, ipAddr net.Addr) (*intf.Status, string) {
	factory := rpcrouter.RouterFactoryFromContext(context)
	cfg := factory.Config
	if len(cfg.Authorizations) == 0 {
		return nil, "default"
	}
	for _, bauth := range cfg.Authorizations {
		if bauth.Authorize(auth.Username, auth.Password, ipAddr.String()) {
			ns := bauth.Namespace
			if ns == "" {
				ns = "default"
			}
			return nil, ns
		}
	}
	return &intf.Status{Code: 401, Reason: "auth failed"}, ""
}

// Call
func (self *JointRPC) Call(context context.Context, req *intf.JSONRPCCallRequest) (*intf.JSONRPCCallResult, error) {
	remotePeer, ok := peer.FromContext(context)
	if !ok {
		return nil, errors.New("cannot get peer info from Call()")
	}

	status, namespace := self.Authorize(context, req.Auth, remotePeer.Addr)
	if status != nil {
		log.WithFields(log.Fields{"ip": remotePeer.Addr}).Warnf("client auth failed")
		return &intf.JSONRPCCallResult{Status: status}, nil
	}

	factory := rpcrouter.RouterFactoryFromContext(context)

	reqmsg, err := msgutil.MessageFromEnvolope(req.Envolope)
	if err != nil {
		return nil, err
	}

	if !reqmsg.IsRequest() {
		return nil, rpcrouter.ErrRequestNotifyRequired
	}
	if reqmsg.TraceId() == "" {
		reqmsg.SetTraceId(misc.NewUuid())
	}
	reqmsg.Log().Debugf("from ip %s", remotePeer.Addr)

	router := factory.CommonRouter()
	if !router.HasMethod(reqmsg.MustMethod()) {
		router = factory.Get(namespace)
	}

	recvmsg, err := router.CallOrNotify(reqmsg,
		namespace,
		rpcrouter.WithBroadcast(req.Broadcast),
		rpcrouter.WithTimeout(time.Second*time.Duration(req.Timeout)))

	if err != nil {
		return nil, err
	}
	if recvmsg == nil {
		recvmsg = jsonrpc.NewResultMessage(reqmsg, nil)
	}
	if !recvmsg.IsResultOrError() {
		log.Warnf("bad recvmsg is neither result nor error %+v", recvmsg)
		return nil, errors.New("recv msg is neigher result or error")
	}
	misc.AssertEqual(recvmsg.TraceId(), reqmsg.TraceId(), "res.traceId != req.traceId")
	res := &intf.JSONRPCCallResult{
		Envolope: msgutil.MessageToEnvolope(recvmsg)}
	return res, nil
}

// Notify
func (self *JointRPC) Notify(context context.Context, req *intf.JSONRPCNotifyRequest) (*intf.JSONRPCNotifyResponse, error) {

	remotePeer, ok := peer.FromContext(context)
	if !ok {
		return nil, errors.New("cannot get peer info from Notify()")
	}

	status, namespace := self.Authorize(context, req.Auth, remotePeer.Addr)
	if status != nil {
		log.WithFields(log.Fields{"ip": remotePeer.Addr}).Warnf("client auth failed")
		return &intf.JSONRPCNotifyResponse{Status: status}, nil
	}

	notifymsg, err := msgutil.MessageFromEnvolope(req.Envolope)
	if err != nil {
		return nil, err
	}

	if !notifymsg.IsNotify() {
		return nil, rpcrouter.ErrRequestNotifyRequired
	}
	if notifymsg.TraceId() == "" {
		notifymsg.SetTraceId(misc.NewUuid())
	}

	notifymsg.Log().Debugf("from ip %s", remotePeer.Addr)
	factory := rpcrouter.RouterFactoryFromContext(context)
	router := factory.CommonRouter()
	if !router.HasMethod(notifymsg.MustMethod()) {
		router = factory.Get(namespace)
	}

	_, err = router.CallOrNotify(notifymsg,
		namespace,
		rpcrouter.WithBroadcast(req.Broadcast))

	if err != nil {
		return nil, err
	}
	resp := &intf.JSONRPCNotifyResponse{Text: "ok"}
	return resp, nil
}

func buildMethodInfos(minfos []rpcrouter.MethodInfo) []*intf.MethodInfo {
	intfMInfos := make([]*intf.MethodInfo, 0)
	for _, minfo := range minfos {
		intfMInfos = append(intfMInfos, msgutil.EncodeMethodInfo(minfo))
	}
	return intfMInfos
}

// ListMethods
func (self *JointRPC) ListMethods(context context.Context, req *intf.ListMethodsRequest) (*intf.ListMethodsResponse, error) {
	remotePeer, ok := peer.FromContext(context)
	if !ok {
		return nil, errors.New("cannot get peer info from ListMethods()")
	}

	status, namespace := self.Authorize(context, req.Auth, remotePeer.Addr)
	if status != nil {
		log.WithFields(log.Fields{"ip": remotePeer.Addr}).Warnf("client auth failed")
		return &intf.ListMethodsResponse{Status: status}, nil
	}

	factory := rpcrouter.RouterFactoryFromContext(context)
	router := factory.Get(namespace)

	minfos := router.GetMethods()

	commonInfos := factory.CommonRouter().GetMethods()
	minfos = append(minfos, commonInfos...)
	intfMInfos := buildMethodInfos(minfos)
	resp := &intf.ListMethodsResponse{Methods: intfMInfos}
	return resp, nil
}

// ListDelegates
func (self *JointRPC) ListDelegates(context context.Context, req *intf.ListDelegatesRequest) (*intf.ListDelegatesResponse, error) {
	remotePeer, ok := peer.FromContext(context)
	if !ok {
		return nil, errors.New("cannot get peer info from ListDelegates()")
	}

	status, namespace := self.Authorize(context, req.Auth, remotePeer.Addr)
	if status != nil {
		log.WithFields(log.Fields{"ip": remotePeer.Addr}).Warnf("client auth failed")
		return &intf.ListDelegatesResponse{Status: status}, nil
	}

	factory := rpcrouter.RouterFactoryFromContext(context)
	router := factory.Get(namespace)
	delegates := router.GetDelegates()
	resp := &intf.ListDelegatesResponse{Delegates: delegates}
	return resp, nil
}

// Lives
// GRPCServer implements dispatch.ISender
type GRPCSender struct {
	stream intf.JointRPC_LiveServer
	done   chan error
}

func NewGRPCSender(stream intf.JointRPC_LiveServer) *GRPCSender {

	sender := &GRPCSender{
		stream: stream,
		done:   make(chan error, 10),
	}
	return sender
}

func (self GRPCSender) SendMessage(context context.Context, msg jsonrpc.IMessage) error {
	return msgutil.GRPCServerSend(self.stream, msg)
}

func (self GRPCSender) SendCmdMsg(context context.Context, cmdMsg rpcrouter.CmdMsg) error {
	return msgutil.GRPCServerSend(self.stream, cmdMsg.Msg)
}

func (self GRPCSender) Done() chan error {
	return self.done
}

func (self *JointRPC) Live(stream intf.JointRPC_LiveServer) error {
	conn := rpcrouter.NewConn()

	factory := rpcrouter.RouterFactoryFromContext(stream.Context())
	//router := factory.Get(conn.Namespace)

	//conn.SetWatchState(true)
	ctx, cancel := context.WithCancel(stream.Context())
	defer func() {
		cancel()
		if conn.Joined() {
			router := factory.Get(conn.Namespace)
			//router.Leave(conn)
			router.ChLeave <- rpcrouter.CmdLeave{Conn: conn}
		}
	}()

	chResult := make(chan dispatch.ResultT, misc.DefaultChanSize())
	sender := NewGRPCSender(stream)
	go dispatch.SenderLoop(ctx, sender, conn, chResult)
	go self.serverReceive(ctx, sender, conn, chResult)

	for {
		select {
		case err, ok := <-sender.Done():
			if !ok {
				log.Debugf("done received not ok")
			} else if err != nil {
				log.Errorf("stream err %+v", err)
			}
			return nil
		}
	}
}

func (self *JointRPC) serverReceive(ctx context.Context, sender *GRPCSender, conn *rpcrouter.ConnT, chResult chan dispatch.ResultT) {
	streamDisp := GetStreamDispatcher()

	for {
		msg, err := msgutil.GRPCServerRecv(sender.stream)
		if err != nil {
			err1 := msgutil.GRPCHandleCodes(err, codes.Canceled)
			sender.Done() <- err1
			return
		}
		msg.Log().Debugf("received from grpc stream")
		instRes := streamDisp.HandleMessage(ctx,
			msg,
			conn.Namespace,
			chResult,
			conn, true)
		if instRes != nil {
			msgutil.GRPCServerSend(sender.stream, instRes)
			if instRes.IsError() {
				sender.Done() <- nil
				return
			}
		}
	}
}

func NewJointRPCServer() *JointRPC {
	return &JointRPC{}
}
