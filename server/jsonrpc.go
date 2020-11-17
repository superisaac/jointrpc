package server

import (
	context "context"
	//json "encoding/json"
	"errors"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	tube "github.com/superisaac/rpctube/tube"
	//"time"
)

func RequestToMessage(req *intf.JSONRPCRequest) (*jsonrpc.RPCMessage, error) {
	json_data := simplejson.New()
	json_data.Set("version", "2.0")
	if req.Id != 0 {
		// idjson, err := simplejson.NewJson([]byte(req.Id))
		// if err != nil {
		// 	return nil, err
		// }
		json_data.Set("id", req.Id) //idjson.Interface())
	}
	json_data.Set("method", req.Method)
	if len(req.Params) > 0 {
		params, err := simplejson.NewJson([]byte(req.Params))
		if err != nil {
			return nil, err
		}
		json_data.Set("params", params)
	} else {
		empty := [](interface{}){}
		json_data.Set("params", empty)
	}
	return jsonrpc.NewRPCMessage(json_data), nil
}

func ResultToMessage(res *intf.JSONRPCResult) (*jsonrpc.RPCMessage, error) {
	json_data := simplejson.New()
	json_data.Set("version", "2.0")
	if res.Id != 0 {
		// idjson, err := simplejson.NewJson([]byte(res.Id))
		// if err != nil {
		// 	return nil, err
		// }
		json_data.Set("id", res.Id) //idjson.Interface())
	}
	if res_ok := res.GetOk(); res_ok != "" {
		parsed, err := simplejson.NewJson([]byte(res_ok))
		if err != nil {
			return nil, err
		}
		json_data.Set("result", parsed)
	} else {
		res_error := res.GetError()
		parsed, err := simplejson.NewJson([]byte(res_error))
		if err != nil {
			return nil, err
		}
		json_data.Set("error", parsed)
	}
	return jsonrpc.NewRPCMessage(json_data), nil
}

func MessageToRequest(msg *jsonrpc.RPCMessage) (*intf.JSONRPCRequest, error) {
	if !msg.IsRequest() && !msg.IsNotify() {
		return nil, errors.New("msg is neither request nor notify")
	}
	req := &intf.JSONRPCRequest{}
	req.Id = int64(msg.Id)
	//	if msg.Id != 0 {
	// idstr, err := json.Marshal(msg.Id)
	// if err != nil {
	// 	return nil, err
	// }
	//req.Id = msg.Id //string(msg.)
	//}
	req.Method = msg.Method
	params, err := jsonrpc.MarshalJson(msg.Params)
	if err != nil {
		return nil, err
	}
	req.Params = params
	return req, nil
}

func MessageToResult(msg *jsonrpc.RPCMessage) (*intf.JSONRPCResult, error) {
	if !msg.IsResult() && !msg.IsError() {
		return nil, errors.New("msg is neither result nor error")
	}
	res := &intf.JSONRPCResult{}
	// idstr, err := json.Marshal(msg.Id)
	// if err != nil {
	// 	return nil, err
	// }
	// res.Id = string(idstr)
	res.Id = int64(msg.Id)
	//res.Id = fmt.Sprintf("%v", msg.Id)
	if msg.IsError() {
		r, err := jsonrpc.MarshalJson(msg.Error)
		if err != nil {
			return nil, err
		}
		res.Result = &intf.JSONRPCResult_Error{Error: r}
	} else {
		r, err := jsonrpc.MarshalJson(msg.Result)
		if err != nil {
			return nil, err
		}
		res.Result = &intf.JSONRPCResult_Ok{Ok: r}
	}
	return res, nil
}

type JSONRPCTube struct {
	intf.UnimplementedJSONRPCTubeServer
}

func leaveConn(conn *tube.ConnT) {
	tube.Tube().Router.ChLeave <- tube.CmdLeave{ConnId: conn.ConnId}
}

func (self *JSONRPCTube) Call(context context.Context, req *intf.JSONRPCRequest) (*intf.JSONRPCResult, error) {
	req_msg, err := RequestToMessage(req)
	if err != nil {
		return nil, err
	}

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
			req, err := MessageToRequest(msg)
			if err != nil {
				panic(err)
			}
			payload := &intf.JSONRPCRequestPacket_Request{Request: req}
			pac := &intf.JSONRPCRequestPacket{Payload: payload}
			err = stream.Send(pac)
			if err != nil {
				panic(err)
			}

		}
	}
}

func (self *JSONRPCTube) Handle(stream intf.JSONRPCTube_HandleServer) error {
	conn := tube.Tube().Router.Join()
	defer leaveConn(conn)
	//defer leaveConn(conn_id)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go relayMessages(ctx, stream, conn.RecvChannel)

	for {
		pac, err := stream.Recv()
		if err != nil {
			return err
		}
		res := pac.GetResult()
		fmt.Printf("result %v\n", res.Id)
		msg, err := ResultToMessage(res)
		if err != nil {
			return err
		}
		tube.Tube().Router.ChMsg <- tube.CmdMsg{Msg: msg, FromConnId: conn.ConnId}
	}
}

func NewJSONRPCTubeServer() *JSONRPCTube {
	return &JSONRPCTube{}
}
