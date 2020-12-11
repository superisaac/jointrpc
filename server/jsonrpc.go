package server

import (
	context "context"
	"errors"
	//json "encoding/json"
	//"errors"
	//"fmt"
	//"log"
	//simplejson "github.com/bitly/go-simplejson"
	intf "github.com/superisaac/rpctube/intf/tube"
	utils "github.com/superisaac/rpctube/utils"
	//jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	tube "github.com/superisaac/rpctube/tube"
	peer "google.golang.org/grpc/peer"
)

type JSONRPCTube struct {
	intf.UnimplementedJSONRPCTubeServer
}

func leaveConn(conn *tube.ConnT) {
	//tube.Tube().Router.ChLeave <- tube.CmdLeave{ConnId: conn.ConnId}
	tube.Tube().Router.Leave(conn)
}

func (self *JSONRPCTube) Call(context context.Context, req *intf.JSONRPCRequest) (*intf.JSONRPCResult, error) {
	utils.DebugLogger.Printf("called method %s", req.Method)
	req_msg, err := RequestToMessage(req)
	if err != nil {
		return nil, err
	}
	//s, _ := req_msg.EncodePretty()
	router := tube.Tube().Router
	recvmsg, err := router.SingleCall(req_msg)
	if err != nil {
		return nil, err
	}
	res, err := MessageToResult(recvmsg)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (self *JSONRPCTube) ListMethods(context context.Context, req *intf.ListMethodsRequest) (*intf.ListMethodsResponse, error) {
	methods := tube.Tube().Router.GetLocalMethods()
	resp := &intf.ListMethodsResponse{Methods: methods}
	return resp, nil
}

func relayMessages(context context.Context, stream intf.JSONRPCTube_HandleServer, recv_ch tube.MsgChannel) {
	for {
		select {
		case <-context.Done():
			return
		case msg := <-recv_ch:
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
	}
}

func (self *JSONRPCTube) Handle(stream intf.JSONRPCTube_HandleServer) error {
	remotePeer, ok := peer.FromContext(stream.Context())
	if !ok {
		return errors.New("cannot get peer info from stream")
	}

	router := tube.Tube().Router
	conn := router.Join()
	conn.PeerAddr = remotePeer.Addr
	utils.DebugLogger.Printf("Joined conn %d", conn.ConnId)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		leaveConn(conn)
	}()

	utils.DebugLogger.Printf("Handler connected, conn %d from ip %s", conn.ConnId, conn.PeerAddr.String())

	go relayMessages(ctx, stream, conn.RecvChannel)

	for {
		up_pac, err := stream.Recv()
		if err != nil {
			utils.InfoLogger.Printf("error on stream Recv() %s", err.Error())
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
				utils.WarningLogger.Printf("error on requesttomessage() %s", err.Error())
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

		reg := up_pac.GetRegisterMethods()
		if reg != nil {
			var loc tube.MethodLocation = tube.Location_Local
			if reg.Location == intf.MethodLocation_REMOTE {
				loc = tube.Location_Remote
			}
			utils.DebugLogger.Printf("conn %d, registered methods %v", conn.ConnId, reg.Methods)
			cmd_reg := tube.CmdReg{
				ConnId:   conn.ConnId,
				Methods:  reg.Methods,
				Location: loc,
			}
			router.ChReg <- cmd_reg
			continue
		}

		unreg := up_pac.GetUnregisterMethods()
		if unreg != nil {
			cmd_unreg := tube.CmdUnreg{
				ConnId:  conn.ConnId,
				Methods: unreg.Methods,
			}
			router.ChUnreg <- cmd_unreg
			continue
		}

	}
}

func NewJSONRPCTubeServer() *JSONRPCTube {
	return &JSONRPCTube{}
}
