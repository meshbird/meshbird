package common_test

import (
	"github.com/gophergala2016/meshbird/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNodeSetAdd(t *testing.T) {
	var val int = 20

	nodeSet := common.NewNodeSet()
	nodeSet.Add("any_key", val)
	got := nodeSet.Select("any_key")

	if assert.NotNil(t, got) {
		gotInt, ok := got.(int)
		if assert.True(t, ok) {
			assert.Equal(t, val, gotInt)
		}
	}
}

func TestNodeSetAddThenRemove(t *testing.T) {
	var val int = 20

	nodeSet := common.NewNodeSet()
	nodeSet.Add("any_key", val)
	nodeSet.Remove("any_key")
	got := nodeSet.Select("any_key")

	assert.Nil(t, got)
}

func TestNodeSetAddThenRemoveThenAddAgain(t *testing.T) {
	var val int = 20

	nodeSet := common.NewNodeSet()
	nodeSet.Add("any_key", 10)
	nodeSet.Remove("any_key")
	nodeSet.Add("any_key", val)
	got := nodeSet.Select("any_key")

	if assert.NotNil(t, got) {
		gotInt, ok := got.(int)
		if assert.True(t, ok) {
			assert.Equal(t, val, gotInt)
		}
	}
}

func TestNodeSetMerge(t *testing.T) {
	var val int = 200

	nodeSet := common.NewNodeSet()
	nodeSet.Add("key1", 10)
	nodeSet.Add("key2", 20)
	nodeSet.Add("key3", 30)
	nodeSet.Remove("key2")

	newNodeSet := common.NewNodeSet()
	newNodeSet.Remove("key3")
	newNodeSet.Add("key2", val)

	nodeSet.Merge(newNodeSet.Data())

	got := nodeSet.Select("key2")

	// Test added
	if assert.NotNil(t, got) {
		gotInt, ok := got.(int)
		if assert.True(t, ok) {
			assert.Equal(t, val, gotInt)
		}
	}

	// Test removed
	assert.Nil(t, nodeSet.Select("key3"))
}
