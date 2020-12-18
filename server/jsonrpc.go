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
	//jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	tube "github.com/superisaac/rpctube/tube"
	peer "google.golang.org/grpc/peer"
)

type JSONRPCTube struct {
	intf.UnimplementedJSONRPCTubeServer
}

/*func leaveConn(conn *tube.ConnT) {
	//tube.Tube().Router.ChLeave <- tube.CmdLeave{ConnId: conn.ConnId}
	tube.Tube().Router.Leave(conn)
}
*/

func (self *JSONRPCTube) Call(context context.Context, req *intf.JSONRPCRequest) (*intf.JSONRPCResult, error) {
	log.Debugf("called method %s", req.Method)
	reqmsg, err := RequestToMessage(req)
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
		res := &intf.JSONRPCResult{Id: ""}
		res.Result = &intf.JSONRPCResult_Ok{Ok: ""}
		return res, nil
	}
	res, err := MessageToResult(recvmsg)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (self *JSONRPCTube) Notify(context context.Context, req *intf.JSONRPCNotifyRequest) (*intf.JSONRPCNotifyResponse, error) {
	log.Debugf("notify %s", req.Method)
	notifymsg, err := NotifyToMessage(req)
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
	resp := &intf.JSONRPCNotifyResponse{}
	return resp, nil
}

func encodeMethodInfo(minfo tube.MethodInfo) *intf.MethodInfo {
	return &intf.MethodInfo{
		Name:      minfo.Name,
		Help:      minfo.Help,
		Delegated: minfo.Delegated,
	}
}

func decodeMethodInfo(iminfo *intf.MethodInfo) tube.MethodInfo {
	return tube.MethodInfo{
		Name:      iminfo.Name,
		Help:      iminfo.Help,
		Delegated: iminfo.Delegated,
	}
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

func relayMessages(context context.Context, stream intf.JSONRPCTube_HandleServer, recv_ch tube.MsgChannel) {
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
		case msg, ok := <-recv_ch:
			if !ok {
				log.Debugf("recv channel closed")
				return
			}
			if msg.IsRequest() || msg.IsNotify() {
				req, err := MessageToRequest(msg)
				if err != nil {
					panic(err)
				}
				payload := &intf.JSONRPCDownPacket_Request{Request: req}
				pac := &intf.JSONRPCDownPacket{Payload: payload}
				err = stream.Send(pac)
				if err != nil {
					panic(err)
				}
			} else {
				// msg.IsResult() || msg.IsError()
				res, err := MessageToResult(msg)
				if err != nil {
					panic(err)
				}
				payload := &intf.JSONRPCDownPacket_Result{Result: res}
				pac := &intf.JSONRPCDownPacket{Payload: payload}
				err = stream.Send(pac)
				if err != nil {
					panic(err)
				}

			}

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
		//leaveConn(conn)
		router.Leave(conn)
		//time.Sleep(1 * time.Second)
	}()

	log.Debugf("Handler connected, conn %d from ip %s", conn.ConnId, conn.PeerAddr.String())

	go relayMessages(ctx, stream, conn.RecvChannel)

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
		req := up_pac.GetRequest()
		if req != nil {
			msg, err := RequestToMessage(req)
			if err != nil {
				log.Warnf("error on requesttomessage() %s", err.Error())
				return err
			}
			cmd_msg := tube.CmdMsg{Msg: msg, FromConnId: conn.ConnId}
			router.ChMsg <- cmd_msg
			continue
		}

		// Handle JSONRPC Result
		res := up_pac.GetResult()
		if res != nil {
			msg, err := ResultToMessage(res)
			if err != nil {
				return err
			}
			cmd_msg := tube.CmdMsg{Msg: msg, FromConnId: conn.ConnId}
			router.ChMsg <- cmd_msg
			continue
		}

		update := up_pac.GetUpdateMethods()
		if update != nil {
			//log.Debugf("update methods %+v", update)
			upMethods := make([]tube.MethodInfo, 0)

			for _, iminfo := range update.Methods {
				minfo := decodeMethodInfo(iminfo)
				upMethods = append(upMethods, minfo)
			}
			log.Debugf("conn %d, update methods %v", conn.ConnId, update.Methods)
			cmd_update := tube.CmdUpdate{
				ConnId:  conn.ConnId,
				Methods: upMethods,
			}
			router.ChUpdate <- cmd_update
			continue
		}

	}
}

func NewJSONRPCTubeServer() *JSONRPCTube {
	return &JSONRPCTube{}
}
