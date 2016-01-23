package common

import (
	"sync"
)

type LocalNode struct {
	Node

	config    *Config
	state     *State

	mutex     sync.Mutex
	waitGroup sync.WaitGroup
}

func NewLocalNode(cfg *Config) *LocalNode {
	n := new(LocalNode)
	return n
}

func (n *LocalNode) Run() error {

	return nil
}

func (n *LocalNode) Stop() error {
	n.waitGroup.Wait()
	return nil
}
