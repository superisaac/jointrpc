package server
import (
	"errors"
	simplejson "github.com/bitly/go-simplejson"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

func RequestToMessage(req *intf.JSONRPCRequest) (*jsonrpc.RPCMessage, error) {

	var id interface{} = req.Id
	if req.Id == 0 {
		id = nil
	}
	params := [](interface{}){}
	if len(req.Params) > 0 {
		paramsJson, err := simplejson.NewJson([]byte(req.Params))
		if err != nil {
			return nil, err
		}
		params = paramsJson.Interface().([]interface{})
	}
	msg := jsonrpc.NewRequestMessage(id, req.Method, params)
	return msg, nil
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
