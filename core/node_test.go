package core

import (
	"testing"
)

func TestNodeRun(t *testing.T) {
	node := NewNode(nil)
	node.Run()
}