package server

import (
	json "encoding/json"
	"errors"
	simplejson "github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	intf "github.com/superisaac/rpctube/intf/tube"
	jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

func RequestToMessage(req *intf.JSONRPCRequest) (jsonrpc.IMessage, error) {

	var id interface{} = nil
	//if req.Id == 0 {
	//if req.Id == "" {
	if req.Id != "" {
		parsed, err := simplejson.NewJson([]byte(req.Id))
		if err != nil {
			return nil, err
		}
		id = parsed.Interface()
	}
	params := [](interface{}){}
	if len(req.Params) > 0 {
		paramsJson, err := simplejson.NewJson([]byte(req.Params))
		if err != nil {
			return nil, err
		}
		paramInterface := paramsJson.Interface()
		if paramInterface != nil {
			params = paramInterface.([]interface{})
		}
	}
	if id != nil {
		return jsonrpc.NewRequestMessage(id, req.Method, params, nil), nil
	} else {
		return jsonrpc.NewNotifyMessage(req.Method, params, nil), nil
	}
}

func NotifyToMessage(req *intf.JSONRPCNotifyRequest) (jsonrpc.IMessage, error) {
	params := [](interface{}){}
	if len(req.Params) > 0 {
		paramsJson, err := simplejson.NewJson([]byte(req.Params))
		if err != nil {
			return nil, err
		}
		paramInterface := paramsJson.Interface()
		if paramInterface != nil {
			params = paramInterface.([]interface{})
		}
	}
	msg := jsonrpc.NewNotifyMessage(req.Method, params, nil)
	return msg, nil
}

func ResultToMessage(res *intf.JSONRPCResult) (jsonrpc.IMessage, error) {
	json_data := simplejson.New()
	json_data.Set("version", "2.0")
	//if res.Id != 0 {
	if res.Id != "" {
		parsed, err := simplejson.NewJson([]byte(res.Id))
		if err != nil {
			return nil, err
		}
		json_data.Set("id", parsed.Interface())

		// idjson, err := simplejson.NewJson([]byte(res.Id))
		// if err != nil {
		// 	return nil, err
		// }
		//json_data.Set("id", res.Id) //idjson.Interface())

	}
	if res_ok := res.GetOk(); res_ok != "" {
		parsed, err := simplejson.NewJson([]byte(res_ok))
		if err != nil {
			return nil, err
		}
		json_data.Set("result", parsed.Interface())
	} else {
		res_error := res.GetError()
		parsed, err := simplejson.NewJson([]byte(res_error))
		if err != nil {
			return nil, err
		}
		json_data.Set("error", parsed.Interface())
	}
	return jsonrpc.Parse(json_data)
}

func MessageToRequest(msg jsonrpc.IMessage) (*intf.JSONRPCRequest, error) {
	if !msg.IsRequest() && !msg.IsNotify() {
		return nil, errors.New("msg is neither request nor notify")
	}
	req := &intf.JSONRPCRequest{}

	if msg.IsRequest() {
		idstr, err := json.Marshal(msg.MustId())
		if err != nil {
			return nil, err
		}
		req.Id = string(idstr)
	} else {
		req.Id = ""
	}
	req.Method = msg.MustMethod()
	params, err := jsonrpc.MarshalJson(msg.MustParams())
	if err != nil {
		return nil, err
	}
	req.Params = params
	return req, nil
}

func MessageToResult(msg jsonrpc.IMessage) (*intf.JSONRPCResult, error) {
	if !msg.IsResult() && !msg.IsError() {
		log.Debugf("msg is %+v", msg)
		return nil, errors.New("msg is neither result nor error")
	}
	res := &intf.JSONRPCResult{}
	if msg.IsResult() {
		iddata, err := json.Marshal(msg.MustId())
		if err != nil {
			return nil, err
		}
		res.Id = string(iddata)
	} else {
		res.Id = ""
	}
	if msg.IsError() {
		r, err := jsonrpc.MarshalJson(msg.MustError())
		if err != nil {
			return nil, err
		}
		res.Result = &intf.JSONRPCResult_Error{Error: r}
	} else {
		r, err := jsonrpc.MarshalJson(msg.MustResult())
		if err != nil {
			return nil, err
		}
		res.Result = &intf.JSONRPCResult_Ok{Ok: r}
	}
	return res, nil
}
