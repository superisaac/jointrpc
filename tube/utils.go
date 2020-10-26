package tube

import (
	"sync/atomic"
	//jsonrpc "github.com/superisaac/rpctube/jsonrpc"
)

var counter uint64 = 10000

func NextUID() uint64 {
	return atomic.AddUint64(&counter, 1)
}

func NextCID() CID {
	return CID(NextUID())
}
