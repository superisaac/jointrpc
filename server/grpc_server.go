package server

import (
	"context"
	"errors"
	"io"
	"net"
	//"strings"
	"time"
	//"time"
	//json "encoding/json"
	//"errors"
	"fmt"
	//"log"
	//simplejson "github.com/bitly/go-simplejson"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"

	log "github.com/sirupsen/logrus"
	encoding "github.com/superisaac/jointrpc/encoding"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	//misc "github.com/superisaac/jointrpc/misc"
	//datadir "github.com/superisaac/jointrpc/datadir"
	//schema "github.com/superisaac/jointrpc/jsonrpc/schema"
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

	reqmsg, err := encoding.MessageFromEnvolope(req.Envolope)
	if err != nil {
		return nil, err
	}

	if !reqmsg.IsRequest() {
		return nil, rpcrouter.ErrRequestNotifyRequired
	}
	if reqmsg.TraceId() == "" {
		reqmsg.SetTraceId(misc.NewUuid())
	}
	reqmsg.Log().Infof("from ip %s", remotePeer.Addr)

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
		misc.AssertEqual(recvmsg.TraceId(), reqmsg.TraceId(), "")
		recvmsg = jsonrpc.NewResultMessage(reqmsg, nil, nil)
	}
	if !recvmsg.IsResultOrError() {
		log.Warnf("bad recvmsg is neither result nor error %+v", recvmsg)
		return nil, errors.New("recv msg is neigher result or error")
	}
	res := &intf.JSONRPCCallResult{
		Envolope: encoding.MessageToEnvolope(recvmsg)}
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

	notifymsg, err := encoding.MessageFromEnvolope(req.Envolope)
	if err != nil {
		return nil, err
	}

	if !notifymsg.IsNotify() {
		return nil, rpcrouter.ErrRequestNotifyRequired
	}
	if notifymsg.TraceId() == "" {
		notifymsg.SetTraceId(misc.NewUuid())
	}

	notifymsg.Log().Infof("from ip %s", remotePeer.Addr)
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
		intfMInfos = append(intfMInfos, encoding.EncodeMethodInfo(minfo))
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

func sendState(state *rpcrouter.ServerState, stream intf.JointRPC_SubscribeStateServer) {
	iState := encoding.EncodeServerState(state)
	payload := &intf.SubscribeStateResponse_State{State: iState}
	resp := &intf.SubscribeStateResponse{Payload: payload}
	stream.Send(resp)
}

// Workers
func sendDownMessage(stream intf.JointRPC_WorkerServer, msg jsonrpc.IMessage) {
	//msg := msgvec.Msg
	msg.Log().Infof("message down to client, %+v", msg)
	envo := encoding.MessageToEnvolope(msg)
	err := stream.Send(envo)
	if err != nil {
		panic(err)
	}
}

func (self *JointRPC) requireAuthState(context context.Context, authReq *intf.AuthRequest, stream intf.JointRPC_SubscribeStateServer) (*rpcrouter.ConnT, error) {
	remotePeer, ok := peer.FromContext(context)
	if !ok {
		return nil, errors.New("cannot get peer info from stream")
	}
	logger := log.WithFields(log.Fields{"ip": remotePeer.Addr})

	auth := authReq.ClientAuth
	status, namespace := self.Authorize(stream.Context(), auth, remotePeer.Addr)
	if status != nil {
		logger.Warnf("client auth failed")
		authResp := &intf.AuthResponse{Status: status, RequestId: authReq.RequestId}
		payload := &intf.SubscribeStateResponse_AuthResponse{AuthResponse: authResp}
		resp := &intf.SubscribeStateResponse{Payload: payload}
		stream.Send(resp)
		return nil, errors.New("client auth failed")
	}

	factory := rpcrouter.RouterFactoryFromContext(context)
	conn := factory.Get(namespace).Join()
	conn.PeerAddr = remotePeer.Addr

	authResp := &intf.AuthResponse{RequestId: authReq.RequestId}
	payload := &intf.SubscribeStateResponse_AuthResponse{AuthResponse: authResp}
	resppac := &intf.SubscribeStateResponse{Payload: payload}
	stream.Send(resppac)
	return conn, nil
}

func (self *JointRPC) requireAuth(stream intf.JointRPC_WorkerServer) (*rpcrouter.ConnT, error) {
	remotePeer, ok := peer.FromContext(stream.Context())
	if !ok {
		return nil, errors.New("cannot get peer info from stream")
	}

	logger := log.WithFields(log.Fields{"ip": remotePeer.Addr})
	envo, err := stream.Recv()
	if err != nil {
		if err == io.EOF {
			log.Debugf("eof met")
			return nil, nil
		} else if grpc.Code(err) == codes.Canceled {
			log.Debugf("stream canceled")
			return nil, nil
		}
		logger.Warnf("error on stream Recv() %s", err.Error())
		return nil, err
	}

	msg, err := encoding.MessageFromEnvolope(envo)
	if err != nil {
		return nil, err
	}
	if !msg.IsRequest() || msg.MustMethod() != "_conn.Authorize" {
		return nil, errors.New("expect _conn.Authorize()")
	}
	params := msg.MustParams()
	if len(params) != 2 {
		return nil, errors.New("len(params) != 2")
	}

	username, ok := params[0].(string)
	if !ok {
		return nil, errors.New("username is not string")
	}

	password, ok := params[1].(string)
	if !ok {
		return nil, errors.New("password is not string")
	}

	auth := &intf.ClientAuth{Username: username, Password: password}
	status, namespace := self.Authorize(stream.Context(), auth, remotePeer.Addr)

	if status != nil {
		return nil, errors.New(fmt.Sprintf("fail to authorize %s", status))
	}
	resmsg := jsonrpc.NewResultMessage(msg, namespace, nil)

	factory := rpcrouter.RouterFactoryFromContext(stream.Context())
	conn := factory.Get(namespace).Join()
	conn.PeerAddr = remotePeer.Addr

	sendDownMessage(stream, resmsg)
	return conn, nil
}

func relayDownMessages(context context.Context, stream intf.JointRPC_WorkerServer, conn *rpcrouter.ConnT, chResult chan dispatch.ResultT) {
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("recovered ERROR %+v", r)
			//stream.Close()
		}
	}()
	for {
		select {
		case <-context.Done():
			log.Debugf("context done")
			return
		case rest, ok := <-chResult:
			if !ok {
				log.Debugf("conn handler channel closed")
				return
			}
			sendDownMessage(stream, rest.ResMsg)
		case msgvec, ok := <-conn.RecvChannel:
			if !ok {
				log.Debugf("recv channel closed")
				return
			}
			sendDownMessage(stream, msgvec.Msg)
		}
	} // and for loop
}

func (self *JointRPC) SubscribeState(authReq *intf.AuthRequest, stream intf.JointRPC_SubscribeStateServer) error {
	conn, err := self.requireAuthState(stream.Context(), authReq, stream)
	if err != nil {
		return err
	}
	if conn == nil {
		return nil
	}

	factory := rpcrouter.RouterFactoryFromContext(stream.Context())
	router := factory.Get(conn.Namespace)
	conn.SetWatchState(true)
	ctx, cancel := context.WithCancel(stream.Context())
	defer func() {
		cancel()
		router.Leave(conn)
	}()

	state := router.GetState()
	sendState(state, stream)

	for {
		select {
		case <-ctx.Done():
			return nil
		case state, ok := <-conn.StateChannel():
			if !ok {
				log.Infof("state channel closed")
				return errors.New("state channnel closed")
			}
			sendState(state, stream)
		}
	}
}

func (self *JointRPC) Worker(stream intf.JointRPC_WorkerServer) error {
	conn, err := self.requireAuth(stream)
	if err != nil {
		return err
	}
	if conn == nil {
		return nil
	}

	factory := rpcrouter.RouterFactoryFromContext(stream.Context())
	router := factory.Get(conn.Namespace)

	conn.SetWatchState(true)
	ctx, cancel := context.WithCancel(stream.Context())
	defer func() {
		cancel()
		router.Leave(conn)
	}()

	chResult := make(chan dispatch.ResultT, 100)
	go relayDownMessages(ctx, stream, conn, chResult)

	for {
		envo, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Debugf("eof met")
				return nil
			} else if grpc.Code(err) == codes.Canceled {
				log.Debugf("stream canceled")
				return nil
			}
			log.Warnf("error on stream Recv() %s", err.Error())
			return err
		}

		msg, err := encoding.MessageFromEnvolope(envo)
		if err != nil {
			conn.Log().Warnf("error on recover message from envo %s", err.Error())
			return err
		}
		// deliver to routers
		msgvec := rpcrouter.MsgVec{
			Msg:        msg,
			Namespace:  conn.Namespace,
			FromConnId: conn.ConnId}

		connDisp := GetConnDispatcher()
		handled := connDisp.HandleRequest(ctx, msgvec, chResult)
		if handled {
			continue
		}
		router.DeliverMessage(rpcrouter.CmdMsg{MsgVec: msgvec})
		continue
	}
}

func NewJointRPCServer() *JointRPC {
	return &JointRPC{}
}
