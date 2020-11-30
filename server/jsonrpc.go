package server
import (
	context "context"
	//json "encoding/json"
	//"errors"
	//"fmt"
	"log"
	//simplejson "github.com/bitly/go-simplejson"
	intf "github.com/superisaac/rpctube/intf/tube"
	//jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	tube "github.com/superisaac/rpctube/tube"
	//"time"
)

type JSONRPCTube struct {
	intf.UnimplementedJSONRPCTubeServer
}

func leaveConn(conn *tube.ConnT) {
	//tube.Tube().Router.ChLeave <- tube.CmdLeave{ConnId: conn.ConnId}
	log.Printf("leave connection %s", conn.ConnId)
	tube.Tube().Router.Leave(conn)
}

func (self *JSONRPCTube) Call(context context.Context, req *intf.JSONRPCRequest) (*intf.JSONRPCResult, error) {
	log.Printf("called method %s", req.Method)
	req_msg, err := RequestToMessage(req)
	if err != nil {
		return nil, err
	}
	s, _ := req_msg.EncodePretty()
	log.Printf("ddd %v", s)
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
	router := tube.Tube().Router
	conn := router.Join()
	defer leaveConn(conn)
	//defer leaveConn(conn_id)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go relayMessages(ctx, stream, conn.RecvChannel)

	for {
		up_pac, err := stream.Recv()
		if err != nil {
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
			log.Printf("reg methods %v", reg.Methods)
			for _, method := range reg.Methods {
				cmd_reg := tube.CmdReg{
					ConnId:   conn.ConnId,
					Method:   method,
					Location: loc,
				}
				router.ChReg <- cmd_reg
			}
			continue
		}

		unreg := up_pac.GetUnregisterMethods()
		if unreg != nil {
			for _, method := range unreg.Methods {
				cmd_unreg := tube.CmdUnreg{
					ConnId: conn.ConnId,
					Method: method,
				}
				router.ChUnreg <- cmd_unreg
			}
			continue
		}

	}
}

func NewJSONRPCTubeServer() *JSONRPCTube {
	return &JSONRPCTube{}
}
