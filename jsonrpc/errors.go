package jsonrpc

var (
	ErrServerError  = &RPCError{100, "server error", false}
	ErrNoSuchMethod = &RPCError{101, "no such method", false}

	ErrNilId       = &RPCError{102, "nil message id", false}
	ErrEmptyMethod = &RPCError{103, "empty method", false}

	ErrParseMessage = &RPCError{104, "parse message error", false}
	ErrMessageType  = &RPCError{105, "wrong message type", false}

	ErrTimeout     = &RPCError{200, "request timeout", true}
	ErrBadResource = &RPCError{201, "bad resource", false}
	ErrWorkerExit  = &RPCError{202, "worker exit", false}
)
