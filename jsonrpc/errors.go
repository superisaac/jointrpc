package jsonrpc

var (
	ErrNoSuchMethod = &RPCError{404, "no such method", false}
	ErrServerError  = &RPCError{500, "server error", false}
)
