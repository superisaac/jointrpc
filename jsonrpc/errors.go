package jsonrpc

var (
	ErrNoSuchMethod = &RPCError{404, "no such method", false}
)
