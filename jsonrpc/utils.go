package jsonrpc

import (
	//json "encoding/json"
	"fmt"
	//"reflect"
	//"errors"
	simplejson "github.com/bitly/go-simplejson"
	"strconv"
)

func (self *RPCError) Error() string {
	return fmt.Sprintf("code=%d, reason=%s", self.Code, self.Reason)
}

func (self RPCError) ToMessage(reqmsg IMessage) *ErrorMessage {
	return RPCErrorMessage(reqmsg, self.Code, self.Reason, self.Retryable)
}

func MarshalJson(data interface{}) (string, error) {
	jsondata := simplejson.New()
	jsondata.SetPath(nil, data)
	bytes, err := jsondata.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func GuessJson(input string) (interface{}, error) {
	if len(input) == 0 {
		return "", nil
	}
	if input == "true" || input == "false" {
		bv, err := strconv.ParseBool(input)
		if err != nil {
			return nil, err
		}
		return bv, nil
	}
	iv, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		return iv, nil
	}
	fv, err := strconv.ParseFloat(input, 64)
	if err == nil {
		return fv, nil
	}

	fc := input[0]
	if fc == '[' {
		parsed, err := simplejson.NewJson([]byte(input))
		if err != nil {
			return nil, err
		}
		return parsed.MustArray(), nil
	} else if fc == '{' {
		parsed, err := simplejson.NewJson([]byte(input))
		if err != nil {
			return nil, err
		}
		return parsed.MustMap(), nil
	} else {
		return input, nil
	}
}

func GuessJsonArray(inputArr []string) ([]interface{}, error) {
	var arr []interface{}
	for _, input := range inputArr {
		v, err := GuessJson(input)
		if err != nil {
			return arr, err
		}
		arr = append(arr, v)
	}
	return arr, nil
}
