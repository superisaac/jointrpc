package jsonrpc

import (
	json "encoding/json"
)


func ValidateNumber(v interface{}) (json.Number, error) {
	if n, ok := v.(json.Number); ok {
		return n, nil
	} else {
		return json.Number("0"), &RPCError{10400, "require a number", false}
	}
}

func ValidateFloat64(v interface{}) (float64, error) {
	n, err := ValidateNumber(v)
	if err != nil {
		return 0, err
	}
	f, err := n.Float64()
	if err != nil {
		return 0, err
	}
	return f, nil
}
