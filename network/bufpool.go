package network

import (
	"bytes"
	"sync"
)

type BufPool struct {
	pool sync.Pool
}

func newBytesSliceFunc(size int) func() interface{} {
	return func() interface{} {
		return make([]byte, size)
	}
}

func NewBufPool(defaultSize int) *BufPool {
	bp := &BufPool{}
	bp.pool.New = newBytesSliceFunc(defaultSize)
	return bp
}

func (bp *BufPool) Put(buf []byte) {
	bp.pool.Put(buf)
}

func (bp *BufPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

type BytesBufferPool struct {
	pool sync.Pool
}

func newBytesBuffer() interface{} {
	return &bytes.Buffer{}
}

func NewBytesBufferPool() *BytesBufferPool {
	bp := &BytesBufferPool{}
	bp.pool.New = newBytesBuffer
	return bp
}

func (bbp *BytesBufferPool) Put(buf *bytes.Buffer) {
	bbp.pool.Put(buf)
}

func (bbp *BytesBufferPool) Get() *bytes.Buffer {
	return bbp.pool.Get().(*bytes.Buffer)
}
