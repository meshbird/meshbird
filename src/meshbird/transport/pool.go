package transport

import (
	"bytes"
	"sync"
)

var noncePool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 12)
	},
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}
