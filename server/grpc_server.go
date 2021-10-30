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
	"fmt"
	//"log"
	//simplejson "github.com/bitly/go-simplejson"
	"github.com/mitchellh/mapstructure"
	//grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"

	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	msgutil "github.com/superisaac/jointrpc/msgutil"
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
	fmt.Printf("router methods %s %+v\n", namespace, minfos)

	commonInfos := factory.CommonRouter().GetMethods()
	fmt.Printf("common router methods %+v\n", commonInfos)
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
func relayDownMessages(context context.Context, stream intf.JointRPC_LiveServer, conn *rpcrouter.ConnT, chResult chan dispatch.ResultT) {
	if false {
		defer func() {
			if r := recover(); r != nil {
				log.Warnf("recovered ERROR %+v", r)
			}
		}()
	}
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
			msgutil.GRPCServerSend(stream, rest.ResMsg)
		case msgvec, ok := <-conn.RecvChannel:
			if !ok {
				log.Debugf("recv channel closed")
				return
			}
			msgutil.GRPCServerSend(stream, msgvec.Msg)

		case cmdMsg, ok := <-conn.ChRouteMsg:
			if !ok {
				log.Debugf("ChRouteMsg closed")
				return
			}
			err := conn.HandleRouteMessage(context, cmdMsg)
			if err != nil {
				panic(err)
			}
		case state, ok := <-conn.StateChannel():
			if !ok {
				log.Debugf("state channel closed")
				return
			}
			stateJson := make(map[string]interface{})
			err := mapstructure.Decode(state, &stateJson)
			if err != nil {
				panic(err)
			}
			ntf := jsonrpc.NewNotifyMessage("_state.changed", []interface{}{stateJson})
			msgutil.GRPCServerSend(stream, ntf)
		}
	} // and for loop
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
	go relayDownMessages(ctx, stream, conn, chResult)
	streamDisp := NewStreamDispatcher()

	for {
		msg, err := msgutil.GRPCServerRecv(stream)
		if err != nil {
			return msgutil.GRPCHandleCodes(err, codes.Canceled)
		}
		msg.Log().Debugf("received from grpc stream")
		msgvec := rpcrouter.MsgVec{
			Msg:        msg,
			Namespace:  conn.Namespace,
			FromConnId: conn.ConnId}

		instRes := streamDisp.HandleMessage(ctx, msgvec, chResult, conn, true)
		if instRes != nil {
			msgutil.GRPCServerSend(stream, instRes)
			if instRes.IsError() {
				return nil
			}
		}
		// if handled {
		// 	continue
		// }
		// if conn.Joined() {
		// 	router := factory.Get(conn.Namespace)
		// 	router.DeliverMessage(rpcrouter.CmdMsg{MsgVec: msgvec})
		// }
		//continue
	}
}

func NewJointRPCServer() *JointRPC {
	return &JointRPC{}
}
