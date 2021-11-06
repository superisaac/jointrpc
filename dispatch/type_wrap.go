package dispatch

import (
	"fmt"
	"reflect"
	"github.com/superisaac/jointrpc/jsonrpc"
	"github.com/superisaac/jointrpc/misc"	
)

func WrappedCall(req jsonrpc.IMessage, tfunc interface{}) {
	funcType := reflect.TypeOf(tfunc)
	fmt.Printf("func type %s\n", funcType.Kind)
	
}

