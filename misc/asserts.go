package misc

import (
	"errors"
)

func Assert(condition bool, hint string) {
	if !condition {
		panic(errors.New(hint))
	}
}
