package network

import "sync"

type BufPool struct {
	pool sync.Pool
}

func NewBufPool(defaultSize int) *BufPool {
	bp := &BufPool{}
	bp.pool.New = func() interface{} {
		return make([]byte, defaultSize)
	}
	return bp
}

func (bp *BufPool) Put(buf []byte) {
	bp.pool.Put(buf)
}

func (bp *BufPool) Get() []byte {
	return bp.pool.Get().([]byte)
}
