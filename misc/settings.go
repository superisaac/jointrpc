package misc

import (
	"os"
	"strconv"
)

var chanSize uint = 0

func DefaultChanSize() uint {
	if chanSize == 0 {
		vs := os.Getenv("SERVER_CHAN_SIZE")
		if vs != "" {
			iv, err := strconv.ParseUint(vs, 10, 64)
			if err != nil {
				panic(err)
			}
			chanSize = uint(iv)
		} else {
			chanSize = 1000
		}
	}
	return chanSize
}
