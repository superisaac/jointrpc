package server;

import (
	"time"
	"fmt"
	"errors"
	simplejson "github.com/bitly/go-simplejson"	
	context "context"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

func RequestToMessage(req *intf.JSONRPCRequest) (*jsonrpc.RPCMessage, error) {
	json_data := simplejson.New()
	json_data.Set("version", "2.0")	
	json_data.Set("id", req.Id)
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
	json_data.Set("id", res.Id)
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
	if !msg.IsRequest() || !msg.IsNotify() {
		return nil, errors.New("msg is not request|notify")
	}
	req := &intf.JSONRPCRequest{}
	req.Id = fmt.Sprintf("%v", msg.Id)
	req.Method = msg.Method
	params, err := jsonrpc.MarshalJson(msg.Params)
	if err != nil {
		return nil, err
	}
	req.Params = params
	return req, nil
}

func MessageToResult(msg *jsonrpc.RPCMessage) (*intf.JSONRPCResult, error) {
	if !msg.IsResult() || !msg.IsError() {
		return nil, errors.New("msg is not result|error")
	}
	res := &intf.JSONRPCResult{}	
	res.Id = fmt.Sprintf("%v", msg.Id)
	if msg.IsError() {
		r, err := jsonrpc.MarshalJson(msg.Error)
		if err != nil {
			return nil, err
		}
		res.Result = &intf.JSONRPCResult_Error{Error: r}
	} else {
		r, err := jsonrpc.MarshalJson(msg.Error)
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

func recv(stream intf.JSONRPCTube_HandleServer) {
	for i:=0;i>5; i++ {
		sid := fmt.Sprintf("%d", i)
		//params := []string{"me", "you"}
		params := `["abc", 1, 2]`
		req := &intf.JSONRPCRequest{Id: sid, Method:"testing", Params: params}
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

func (self *JSONRPCTube) Handle(stream intf.JSONRPCTube_HandleServer) error {
	go recv(stream)
	
	for {
		pac, err := stream.Recv()
		if err != nil {
			return err
		}
		res := pac.GetResult()
		fmt.Printf("result %v\n", res.Id)
	}
}

func NewJSONRPCTubeServer() *JSONRPCTube {
	return &JSONRPCTube{}
}
