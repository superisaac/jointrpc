package server;

import (
	"time"
	"fmt"
	"errors"
	simplejson "github.com/bitly/go-simplejson"	
	context "context"
	tube "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

func RequestToMessage(req *tube.JSONRPCRequest) (*jsonrpc.RPCMessage, error) {
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

func ResultToMessage(res *tube.JSONRPCResult) (*jsonrpc.RPCMessage, error) {
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

func MessageToRequest(msg *jsonrpc.RPCMessage) (*tube.JSONRPCRequest, error) {
	if !msg.IsRequest() || !msg.IsNotify() {
		return nil, errors.New("msg is not request|notify")
	}
	req := &tube.JSONRPCRequest{}
	req.Id = fmt.Sprintf("%v", msg.Id)
	req.Method = msg.Method
	params, err := jsonrpc.MarshalJson(msg.Params)
	if err != nil {
		return nil, err
	}
	req.Params = params
	return req, nil
}

func MessageToResult(msg *jsonrpc.RPCMessage) (*tube.JSONRPCResult, error) {
	if !msg.IsResult() || !msg.IsError() {
		return nil, errors.New("msg is not result|error")
	}
	res := &tube.JSONRPCResult{}	
	res.Id = fmt.Sprintf("%v", msg.Id)
	if msg.IsError() {
		r, err := jsonrpc.MarshalJson(msg.Error)
		if err != nil {
			return nil, err
		}
		res.Result = &tube.JSONRPCResult_Error{Error: r}
	} else {
		r, err := jsonrpc.MarshalJson(msg.Error)
		if err != nil {
			return nil, err
		}
		res.Result = &tube.JSONRPCResult_Ok{Ok: r}
	}
	return res, nil
}


type JSONRPCTube struct {
	tube.UnimplementedJSONRPCTubeServer
}


func (self *JSONRPCTube) Call(context context.Context, req *tube.JSONRPCRequest) (*tube.JSONRPCResult, error) {
	req_msg, err := RequestToMessage(req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("sss %v %v\n", req.Method, req_msg.Id)		
	ok := &tube.JSONRPCResult_Ok{Ok: "okokook"}
	res := &tube.JSONRPCResult{Id: req.Id, Result: ok}
	return res, nil
}

func recv(stream tube.JSONRPCTube_HandleServer) {
	for i:=0;i>5; i++ {
		sid := fmt.Sprintf("%d", i)
		//params := []string{"me", "you"}
		params := `["abc", 1, 2]`
		req := &tube.JSONRPCRequest{Id:sid, Method:"testing", Params: params}
		err := stream.Send(req)
		if err != nil {
			//stream.Close()
			break
		}
		time.Sleep(3000 * time.Millisecond)
	}
}

func (self *JSONRPCTube) Handle(stream tube.JSONRPCTube_HandleServer) error {
	go recv(stream)
	
	for {
		res, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("result %v\n", res.Id)
	}
}

func NewJSONRPCTubeServer() *JSONRPCTube {
	return &JSONRPCTube{}
}
