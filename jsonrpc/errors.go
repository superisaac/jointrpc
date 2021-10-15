package jsonrpc

// https://www.jsonrpc.org/specification

var (
	ErrServerError = &RPCError{100, "server error", nil}
	ErrNilId       = &RPCError{102, "nil message id", nil}

	ErrMethodNotFound = &RPCError{-32601, "method not found", nil}
	ErrEmptyMethod    = &RPCError{-32601, "empty method", nil}

	ErrParseMessage = &RPCError{-32700, "parse error", nil}
	ErrMessageType  = &RPCError{105, "wrong message type", nil}

	ErrTimeout     = &RPCError{200, "request timeout", nil}
	ErrBadResource = &RPCError{201, "bad resource", nil}
	ErrWorkerExit  = &RPCError{202, "worker exit", nil}
)
