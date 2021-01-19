package server

import (
	"context"
	"errors"
	"io"
	"strings"
	//"time"
	//json "encoding/json"
	//"errors"
	"fmt"
	//"log"
	//simplejson "github.com/bitly/go-simplejson"
	uuid "github.com/google/uuid"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"

	log "github.com/sirupsen/logrus"
	encoding "github.com/superisaac/jointrpc/encoding"
	intf "github.com/superisaac/jointrpc/intf/jointrpc"
	jsonrpc "github.com/superisaac/jointrpc/jsonrpc"
	misc "github.com/superisaac/jointrpc/misc"
	//misc "github.com/superisaac/jointrpc/misc"
	schema "github.com/superisaac/jointrpc/jsonrpc/schema"
	"github.com/superisaac/jointrpc/rpcrouter"
	peer "google.golang.org/grpc/peer"
)

type JointRPC struct {
	intf.UnimplementedJointRPCServer
}

// Call
func (self *JointRPC) Call(context context.Context, req *intf.JSONRPCCallRequest) (*intf.JSONRPCCallResult, error) {
	remotePeer, ok := peer.FromContext(context)
	if !ok {
		return nil, errors.New("cannot get peer info from Call()")
	}

	reqmsg, err := encoding.MessageFromEnvolope(req.Envolope)
	if err != nil {
		return nil, err
	}

	if !reqmsg.IsRequest() {
		return nil, rpcrouter.ErrRequestNotifyRequired
	}
	if reqmsg.TraceId() == "" {
		reqmsg.SetTraceId(uuid.New().String())
	}
	reqmsg.Log().Infof("from ip %s", remotePeer.Addr)
	router := rpcrouter.RouterFromContext(context)
	msgvec := rpcrouter.MsgVec{Msg: reqmsg}
	recvmsg, err := router.CallOrNotify(msgvec, req.Broadcast)

	if err != nil {
		return nil, err
	}
	if recvmsg == nil {
		misc.AssertEqual(recvmsg.TraceId(), msgvec.Msg.TraceId(), "")
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

	notifymsg, err := encoding.MessageFromEnvolope(req.Envolope)
	if err != nil {
		return nil, err
	}

	if !notifymsg.IsNotify() {
		return nil, rpcrouter.ErrRequestNotifyRequired
	}
	if notifymsg.TraceId() == "" {
		notifymsg.SetTraceId(uuid.New().String())
	}

	notifymsg.Log().Infof("from ip %s", remotePeer.Addr)
	router := rpcrouter.RouterFromContext(context)
	msgvec := rpcrouter.MsgVec{Msg: notifymsg}

	_, err = router.CallOrNotify(msgvec, req.Broadcast)

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
	router := rpcrouter.RouterFromContext(context)
	minfos := router.GetMethods()
	intfMInfos := buildMethodInfos(minfos)
	resp := &intf.ListMethodsResponse{Methods: intfMInfos}
	return resp, nil
}

// DeclareMethods
func (self *JointRPC) DeclareMethods(context context.Context, req *intf.DeclareMethodsRequest) (*intf.DeclareMethodsResponse, error) {
	router := rpcrouter.RouterFromContext(context)
	conn, found := router.GetConnByPublicId(req.ConnPublicId)
	if !found {
		intfErr := &intf.Error{Code: 404, Reason: "conn not found"}
		return &intf.DeclareMethodsResponse{Error: intfErr}, nil
	}

	//log.Debugf("update methods %+v", update)
	upMethods := make([]rpcrouter.MethodInfo, 0)
	var methodNames []string
	for _, iminfo := range req.Methods {
		minfo := encoding.DecodeMethodInfo(iminfo)
		if strings.HasPrefix(minfo.Name, ".") {
			conn.Log().WithFields(log.Fields{
				"rpc": "DeclareMethods",
			}).Warnf("method %s cannot prefix with .", minfo.Name)
			intfErr := &intf.Error{
				Code:   11400,
				Reason: fmt.Sprintf("method %s cannot prefix with .", minfo.Name)}
			return &intf.DeclareMethodsResponse{Error: intfErr}, nil
		}
		methodNames = append(methodNames, minfo.Name)
		_, err := minfo.SchemaOrError()
		if err != nil {
			if buildError, ok := err.(*schema.SchemaBuildError); ok {
				// parse schema error
				conn.Log().WithFields(log.Fields{
					"rpc": "DeclareMethods",
				}).Warnf("error build schema %s, %+v", buildError.Error(), iminfo)
				intfErr := &intf.Error{
					Code:   11401,
					Reason: fmt.Sprintf("build schema error %s", buildError.Error())}
				return &intf.DeclareMethodsResponse{Error: intfErr}, nil
			}
			return &intf.DeclareMethodsResponse{}, err
		}
		upMethods = append(upMethods, *minfo)
	}

	conn.Log().Infof("declared methods %v", methodNames)
	cmdServe := rpcrouter.CmdServe{
		ConnId:  conn.ConnId,
		Methods: upMethods,
	}
	router.ChServe <- cmdServe
	return &intf.DeclareMethodsResponse{}, nil
}

// DeclareDelegates
func (self *JointRPC) DeclareDelegates(context context.Context, req *intf.DeclareDelegatesRequest) (*intf.DeclareDelegatesResponse, error) {
	router := rpcrouter.RouterFromContext(context)
	conn, found := router.GetConnByPublicId(req.ConnPublicId)
	if !found {
		intfErr := &intf.Error{Code: 404, Reason: "conn not found"}
		return &intf.DeclareDelegatesResponse{Error: intfErr}, nil
	}

	// TODO: validate delegate methods
	cmdDelegate := rpcrouter.CmdDelegate{
		ConnId:      conn.ConnId,
		MethodNames: req.Methods,
	}
	router.ChDelegate <- cmdDelegate

	return &intf.DeclareDelegatesResponse{}, nil
}

// ListDelegates
func (self *JointRPC) ListDelegates(context context.Context, req *intf.ListDelegatesRequest) (*intf.ListDelegatesResponse, error) {
	router := rpcrouter.RouterFromContext(context)
	delegates := router.GetDelegates()
	resp := &intf.ListDelegatesResponse{Delegates: delegates}
	return resp, nil
}

func sendState(state *rpcrouter.ServerState, stream intf.JointRPC_HandleServer) {
	iState := encoding.EncodeServerState(state)
	payload := &intf.JointRPCDownPacket_State{State: iState}
	downpac := &intf.JointRPCDownPacket{Payload: payload}
	stream.Send(downpac)
}

func sendServerGreeting(connPublicId string, stream intf.JointRPC_HandleServer) {
	greeting := &intf.ServerGreeting{ConnPublicId: connPublicId}
	payload := &intf.JointRPCDownPacket_Greeting{Greeting: greeting}
	downpac := &intf.JointRPCDownPacket{Payload: payload}
	stream.Send(downpac)
}

// Handler
func downMsgToDeliver(context context.Context, msgvec rpcrouter.MsgVec, stream intf.JointRPC_HandleServer, conn *rpcrouter.ConnT) {
	msg := msgvec.Msg

	router := rpcrouter.RouterFromContext(context)
	if msg.IsRequestOrNotify() {
		// validate params
		if validated, errmsg := conn.ValidateMsg(msg); !validated && errmsg != nil {
			msgvec := rpcrouter.MsgVec{
				Msg:        errmsg,
				FromConnId: conn.ConnId,
			}
			router.ChMsg <- rpcrouter.CmdMsg{MsgVec: msgvec}
		}
	}

	msg.Log().Infof("message down to client")
	env := encoding.MessageToEnvolope(msg)
	payload := &intf.JointRPCDownPacket_Envolope{Envolope: env}
	pac := &intf.JointRPCDownPacket{Payload: payload}
	err := stream.Send(pac)
	if err != nil {
		panic(err)
	}
}

func relayDownMessages(context context.Context, stream intf.JointRPC_HandleServer, conn *rpcrouter.ConnT) {
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
		case state, ok := <-conn.StateChannel():
			if !ok {
				log.Debugf("state channel closed")
				return
			}
			sendState(state, stream)
		case msgvec, ok := <-conn.RecvChannel:
			if !ok {
				log.Debugf("recv channel closed")
				return
			}
			downMsgToDeliver(context, msgvec, stream, conn)
		}
	} // and for loop
}

func (self *JointRPC) Handle(stream intf.JointRPC_HandleServer) error {
	remotePeer, ok := peer.FromContext(stream.Context())
	if !ok {
		return errors.New("cannot get peer info from stream")
	}

	router := rpcrouter.RouterFromContext(stream.Context())
	conn := router.Join(true)
	conn.PeerAddr = remotePeer.Addr
	conn.SetWatchState(true)
	log.Debugf("Joined conn %d", conn.ConnId)

	ctx, cancel := context.WithCancel(stream.Context())
	defer func() {
		cancel()
		router.Leave(conn)
	}()

	conn.Log().WithFields(log.Fields{
		"ip": conn.PeerAddr.String(),
	}).Debugf("handler connected")

	sendServerGreeting(conn.PublicId(), stream)
	state := router.GetState()
	sendState(state, stream)

	go relayDownMessages(ctx, stream, conn)

	for {
		uppac, err := stream.Recv()
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
		// Pong on Ping
		ping := uppac.GetPing()
		if ping != nil {
			pong := &intf.PONG{Text: ping.Text}
			payload := &intf.JointRPCDownPacket_Pong{Pong: pong}
			down_pac := &intf.JointRPCDownPacket{Payload: payload}

			stream.Send(down_pac)
			continue
		}

		// Handle JSONRPC Request
		envo := uppac.GetEnvolope()
		if envo != nil {
			msg, err := encoding.MessageFromEnvolope(envo)
			if err != nil {
				log.Warnf("error on requesttomessage() %s", err.Error())
				return err
			}
			msgvec := rpcrouter.MsgVec{
				Msg:        msg,
				FromConnId: conn.ConnId}
			router.ChMsg <- rpcrouter.CmdMsg{MsgVec: msgvec}
			continue
		}

	}
}

func NewJointRPCServer() *JointRPC {
	return &JointRPC{}
}
