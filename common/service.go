package common

import (
	"sync"
)

type Service interface {
	Init(*LocalNode, *sync.WaitGroup) error
	Name() string
	Run() error
	Stop()
}