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
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"

	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	schema "github.com/superisaac/rpctube/jsonrpc/schema"
	tube "github.com/superisaac/rpctube/tube"
	peer "google.golang.org/grpc/peer"
)

type JSONRPCTube struct {
	intf.UnimplementedJSONRPCTubeServer
}

func (self *JSONRPCTube) Call(context context.Context, req *intf.JSONRPCCallRequest) (*intf.JSONRPCCallResult, error) {
	reqmsg, err := jsonrpc.ParseBytes([]byte(req.Envolope.Body))
	if err != nil {
		return nil, err
	}

	if !reqmsg.IsRequest() {
		return nil, tube.ErrRequestNotifyRequired
	}

	router := tube.Tube().Router
	recvmsg, err := router.SingleCall(reqmsg, false)
	if err != nil {
		return nil, err
	}
	if recvmsg == nil {
		recvmsg = jsonrpc.NewResultMessage(reqmsg.MustId(), nil, nil)
	}
	if !recvmsg.IsResultOrError() {
		log.Warnf("bad recvmsg is neither result nor error %+v", recvmsg)
		return nil, errors.New("recv msg is neigher result or error")
	}
	res := &intf.JSONRPCCallResult{
		Envolope: &intf.JSONRPCEnvolope{
			Body: recvmsg.MustString()}}
	return res, nil
}

func (self *JSONRPCTube) Notify(context context.Context, req *intf.JSONRPCNotifyRequest) (*intf.JSONRPCNotifyResponse, error) {
	notifymsg, err := jsonrpc.ParseBytes([]byte(req.Envolope.Body))
	if err != nil {
		return nil, err
	}

	if !notifymsg.IsNotify() {
		return nil, tube.ErrRequestNotifyRequired
	}

	router := tube.Tube().Router
	_, err = router.SingleCall(notifymsg, req.Broadcast)
	if err != nil {
		return nil, err
	}
	resp := &intf.JSONRPCNotifyResponse{Text: "ok"}
	return resp, nil
}

// turn from tube's struct to protobuf message
func encodeMethodInfo(minfo tube.MethodInfo) *intf.MethodInfo {
	return &intf.MethodInfo{
		Name:      minfo.Name,
		Help:      minfo.Help,
		Delegated: minfo.Delegated,
	}
}

// turn from protobuf to tube's struct
func decodeMethodInfo(iminfo *intf.MethodInfo) (*tube.MethodInfo, error) {
	var s schema.Schema
	var err error
	if iminfo.SchemaJson != "" {
		builder := schema.NewSchemaBuilder()
		s, err = builder.BuildBytes([]byte(iminfo.SchemaJson))
		if err != nil {
			return nil, err
		}
	}
	return &tube.MethodInfo{
		Name:      iminfo.Name,
		Help:      iminfo.Help,
		Schema:    s,
		Delegated: iminfo.Delegated,
	}, nil
}

func (self *JSONRPCTube) ListMethods(context context.Context, req *intf.ListMethodsRequest) (*intf.ListMethodsResponse, error) {
	minfos := tube.Tube().Router.GetLocalMethods()
	intfMInfos := make([]*intf.MethodInfo, 0)
	for _, minfo := range minfos {
		intfMInfos = append(intfMInfos, encodeMethodInfo(minfo))
	}
	resp := &intf.ListMethodsResponse{MethodInfos: intfMInfos}
	log.Debugf("list methods resp %v", resp)
	return resp, nil
}

func downMsgReceived(msgvec tube.MsgVec, stream intf.JSONRPCTube_HandleServer, conn *tube.ConnT) {
	msg := msgvec.Msg

	if msg.IsRequest() || msg.IsNotify() {
		// validate params
		validated, errmsg := conn.ValidateMsg(msg)
		if !validated {
			if errmsg != nil {
				msgvec := tube.MsgVec{
					Msg:        errmsg,
					FromConnId: conn.ConnId,
				}
				router := tube.Tube().Router
				router.ChMsg <- tube.CmdMsg{MsgVec: msgvec}
			}
		}
	}

	env := &intf.JSONRPCEnvolope{Body: msg.MustString()}
	payload := &intf.JSONRPCDownPacket_Envolope{Envolope: env}
	pac := &intf.JSONRPCDownPacket{Payload: payload}
	err := stream.Send(pac)
	if err != nil {
		panic(err)
	}
}

func relayMessages(context context.Context, stream intf.JSONRPCTube_HandleServer, conn *tube.ConnT) {
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
		case msgvec, ok := <-conn.RecvChannel:
			if !ok {
				log.Debugf("recv channel closed")
				return
			}
			downMsgReceived(msgvec, stream, conn)
		}
	} // and for loop
}

func (self *JSONRPCTube) Handle(stream intf.JSONRPCTube_HandleServer) error {
	remotePeer, ok := peer.FromContext(stream.Context())
	if !ok {
		return errors.New("cannot get peer info from stream")
	}

	router := tube.Tube().Router
	conn := router.Join()
	conn.PeerAddr = remotePeer.Addr
	log.Debugf("Joined conn %d", conn.ConnId)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		router.Leave(conn)
	}()

	log.Debugf("Handler connected, conn %d from ip %s", conn.ConnId, conn.PeerAddr.String())

	go relayMessages(ctx, stream, conn)

	for {
		up_pac, err := stream.Recv()
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
		ping := up_pac.GetPing()
		if ping != nil {
			pong := &intf.PONG{Text: ping.Text}
			payload := &intf.JSONRPCDownPacket_Pong{Pong: pong}
			down_pac := &intf.JSONRPCDownPacket{Payload: payload}

			stream.Send(down_pac)
			continue
		}

		// Handle JSONRPC Request
		env := up_pac.GetEnvolope()
		if env != nil {
			msg, err := jsonrpc.ParseBytes([]byte(env.Body))
			if err != nil {
				log.Warnf("error on requesttomessage() %s", err.Error())
				return err
			}
			msgvec := tube.MsgVec{Msg: msg, FromConnId: conn.ConnId}
			router.ChMsg <- tube.CmdMsg{MsgVec: msgvec}
			continue
		}

		update := up_pac.GetUpdateMethods()
		if update != nil {
			//log.Debugf("update methods %+v", update)
			upMethods := make([]tube.MethodInfo, 0)

			for _, iminfo := range update.Methods {
				minfo, err := decodeMethodInfo(iminfo)
				if err != nil {
					if buildError, ok := err.(*schema.SchemaBuildError); ok {
						// parse schema error
						log.Warnf("error build schema %s, %+v", buildError.Error(), iminfo)
						resp := &intf.UpdateMethodsResponse{Text: buildError.Error()}
						payload := &intf.JSONRPCDownPacket_UpdateMethods{UpdateMethods: resp}
						down_pac := &intf.JSONRPCDownPacket{Payload: payload}
						stream.Send(down_pac)
						// close the handle
						return nil
					}
					return err
				}
				upMethods = append(upMethods, *minfo)
			}
			log.Debugf("conn %d, update methods %v", conn.ConnId, update.Methods)
			cmdUpdate := tube.CmdUpdate{
				ConnId:  conn.ConnId,
				Methods: upMethods,
			}
			router.ChUpdate <- cmdUpdate
			continue
		}

	}
}

func NewJSONRPCTubeServer() *JSONRPCTube {
	return &JSONRPCTube{}
}
