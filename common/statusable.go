package common

import (
	"sync/atomic"
)

const (
	StatusCreated = iota
	StatusRunned
	StatusStopping
	StatusStopped
)

type Statusable struct {
	status uint32
}

func (s *Statusable) Status() uint32 {
	return atomic.LoadUint32(&s.status)
}

func (s *Statusable) SetStatus(v uint32) {
	atomic.StoreUint32(&s.status, v)
}
