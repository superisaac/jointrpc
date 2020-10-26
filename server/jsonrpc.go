package server

import (
	context "context"
	json "encoding/json"
	"errors"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
	tube "github.com/superisaac/rpctube/tube"
	"time"
)

func RequestToMessage(req *intf.JSONRPCRequest) (*jsonrpc.RPCMessage, error) {
	json_data := simplejson.New()
	json_data.Set("version", "2.0")
	if req.Id != "" {
		idjson, err := simplejson.NewJson([]byte(req.Id))
		if err != nil {
			return nil, err
		}
		json_data.Set("id", idjson.Interface())
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
	if res.Id != "" {
		idjson, err := simplejson.NewJson([]byte(res.Id))
		if err != nil {
			return nil, err
		}
		json_data.Set("id", idjson.Interface())
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
	if msg.Id != nil {
		idstr, err := json.Marshal(msg.Id)
		if err != nil {
			return nil, err
		}
		req.Id = string(idstr)
	}
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
	idstr, err := json.Marshal(msg.Id)
	if err != nil {
		return nil, err
	}
	res.Id = string(idstr)
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

func (self *JSONRPCTube) Call(context context.Context, req *intf.JSONRPCRequest) (*intf.JSONRPCResult, error) {
	req_msg, err := RequestToMessage(req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("sss %v %v\n", req.Method, req_msg.Id)
	ok := &intf.JSONRPCResult_Ok{Ok: "okokook"}
	res := &intf.JSONRPCResult{Id: req.Id, Result: ok}
	return res, nil
}

func relayResult(stream intf.JSONRPCTube_HandleServer) {
	for i := 0; i > 5; i++ {
		sid := fmt.Sprintf("%d", i)
		//params := []string{"me", "you"}
		params := `["abc", 1, 2]`
		req := &intf.JSONRPCRequest{Id: sid, Method: "testing", Params: params}
		payload := &intf.JSONRPCRequestPacket_Request{Request: req}
		pac := &intf.JSONRPCRequestPacket{Payload: payload}
		err := stream.Send(pac)
		if err != nil {
			//stream.Close()
			break
		}
		time.Sleep(3000 * time.Millisecond)
	}
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
	conn_id := tube.NextCID()
	recv_ch := make(tube.MsgChannel, 100)

	tube.Tube().Router.ChJoin <- tube.CmdJoin{RecvChannel: recv_ch, ConnId: conn_id}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go relayMessages(ctx, stream, recv_ch)
	//go relayResult(stream)

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
		tube.Tube().Router.ChMsg <- tube.CmdMsg{Msg: msg, FromConnId: conn_id}
	}
}

func NewJSONRPCTubeServer() *JSONRPCTube {
	return &JSONRPCTube{}
}
