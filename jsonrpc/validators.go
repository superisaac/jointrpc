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

func ValidateFloat(v interface{}) (float64, error) {
	n, err := ValidateNumber(v)
	if err != nil {
		return 0, err
	}
	f, err := n.Float64()
	if err != nil {
		return 0, &RPCError{10401, "require a float number", false}
	}
	return f, nil
}

func ValidateInt(v interface{}) (int64, error) {
	n, err := ValidateNumber(v)
	if err != nil {
		return 0, err
	}
	i, err := n.Int64()
	if err != nil {
		return 0, &RPCError{11402, "require a int number", false}
	}
	return i, nil
}
