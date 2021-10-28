package misc

import (
	"fmt"
	"github.com/pkg/errors"
)

func Assert(condition bool, hint string) {
	if !condition {
		panic(errors.New(hint))
	}
}

func AssertEqual(a interface{}, b interface{}, hint string) {
	if a != b {
		panic(errors.New(fmt.Sprintf("assert failed, %v == %v, %s", a, b, hint)))
	}
}
