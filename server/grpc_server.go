package server

import (
	"context"
	"errors"
	"io"
	//"time"
	//json "encoding/json"
	//"errors"
	//"fmt"
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

func (self *JointRPC) Call(context context.Context, req *intf.JSONRPCCallRequest) (*intf.JSONRPCCallResult, error) {
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

func (self *JointRPC) Notify(context context.Context, req *intf.JSONRPCNotifyRequest) (*intf.JSONRPCNotifyResponse, error) {
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

func (self *JointRPC) ListMethods(context context.Context, req *intf.ListMethodsRequest) (*intf.ListMethodsResponse, error) {
	router := rpcrouter.RouterFromContext(context)
	minfos := router.GetMethods()
	intfMInfos := buildMethodInfos(minfos)
	resp := &intf.ListMethodsResponse{MethodInfos: intfMInfos}
	return resp, nil
}

func (self *JointRPC) ListDelegates(context context.Context, req *intf.ListDelegatesRequest) (*intf.ListDelegatesResponse, error) {
	router := rpcrouter.RouterFromContext(context)
	delegates := router.GetDelegates()
	resp := &intf.ListDelegatesResponse{Delegates: delegates}
	return resp, nil
}

func sendState(state *rpcrouter.TubeState, stream intf.JointRPC_HandleServer) {
	iState := encoding.EncodeTubeState(state)
	payload := &intf.JointRPCDownPacket_State{State: iState}
	downpac := &intf.JointRPCDownPacket{Payload: payload}
	stream.Send(downpac)
}

// Handler
func downMsgToDeliver(context context.Context, msgvec rpcrouter.MsgVec, stream intf.JointRPC_HandleServer, conn *rpcrouter.ConnT) {
	msg := msgvec.Msg

	router := rpcrouter.RouterFromContext(context)
	if msg.IsRequest() || msg.IsNotify() {
		// validate params
		validated, errmsg := conn.ValidateMsg(msg)
		if !validated {
			if errmsg != nil {
				msgvec := rpcrouter.MsgVec{
					Msg:        errmsg,
					FromConnId: conn.ConnId,
				}
				router.ChMsg <- rpcrouter.CmdMsg{MsgVec: msgvec}
			}
		}
	}

	env := encoding.MessageToEnvolope(msg)
	payload := &intf.JointRPCDownPacket_Envolope{Envolope: env}
	pac := &intf.JointRPCDownPacket{Payload: payload}
	err := stream.Send(pac)
	if err != nil {
		panic(err)
	}
}

func relayMessages(context context.Context, stream intf.JointRPC_HandleServer, conn *rpcrouter.ConnT) {
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
	conn := router.Join()
	conn.PeerAddr = remotePeer.Addr
	conn.SetWatchState(true)
	log.Debugf("Joined conn %d", conn.ConnId)

	ctx, cancel := context.WithCancel(stream.Context())
	defer func() {
		cancel()
		router.Leave(conn)
	}()

	log.Debugf("Handler connected, conn %d from ip %s", conn.ConnId, conn.PeerAddr.String())

	// send initial state
	state := router.GetState()
	sendState(state, stream)

	go relayMessages(ctx, stream, conn)

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

		update := uppac.GetCanServe()
		if update != nil {
			//log.Debugf("update methods %+v", update)
			upMethods := make([]rpcrouter.MethodInfo, 0)

			for _, iminfo := range update.Methods {
				minfo := encoding.DecodeMethodInfo(iminfo)
				_, err := minfo.SchemaOrError()
				if err != nil {
					if buildError, ok := err.(*schema.SchemaBuildError); ok {
						// parse schema error
						log.Warnf("error build schema %s, %+v", buildError.Error(), iminfo)
						resp := &intf.CanServeResponse{Text: buildError.Error()}
						payload := &intf.JointRPCDownPacket_CanServe{CanServe: resp}
						down_pac := &intf.JointRPCDownPacket{Payload: payload}
						stream.Send(down_pac)
						// close the handle
						return nil
					}
					return err
				}
				upMethods = append(upMethods, *minfo)
			}
			log.Debugf("conn %d, update methods %+v", conn.ConnId, update.Methods)
			cmdServe := rpcrouter.CmdServe{
				ConnId:  conn.ConnId,
				Methods: upMethods,
			}
			router.ChServe <- cmdServe
			continue
		}

		delegate := uppac.GetCanDelegate()
		if delegate != nil {
			log.Debugf("conn %d, delegate methods %+v", conn.ConnId, delegate.Methods)
			// TODO: validate delegate methods
			cmdDelegate := rpcrouter.CmdDelegate{
				ConnId:      conn.ConnId,
				MethodNames: delegate.Methods,
			}
			router.ChDelegate <- cmdDelegate

			resp := &intf.CanDelegateResponse{Text: "ok"}
			payload := &intf.JointRPCDownPacket_CanDelegate{CanDelegate: resp}
			down_pac := &intf.JointRPCDownPacket{Payload: payload}
			stream.Send(down_pac)

			continue
		}

	}
}

func NewJointRPCServer() *JointRPC {
	return &JointRPC{}
}
