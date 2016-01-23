package common

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

type (
	NodeSetElement struct {
		Add, Remove time.Time
		Value       interface{}
	}

	NodeSet struct {
		data map[string]NodeSetElement
		lock sync.RWMutex
	}
)

func (ne NodeSetElement) String() string {
	return fmt.Sprintf("Add: %s; Remove: %s; Type: %T; Value: %v", ne.Add, ne.Remove, ne.Value, ne.Value)
}

func NewNodeSet() *NodeSet {
	return &NodeSet{
		data: make(map[string]NodeSetElement),
	}
}

func (s *NodeSet) Merge(values map[string]NodeSetElement) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for k, v := range values {
		el, exist := s.data[k]
		if exist {
			el.Value = v.Value

			if el.Add.Before(v.Add) {
				el.Add = v.Add
			}
			if el.Remove.Before(v.Remove) {
				el.Remove = v.Remove
			}
		} else {
			el = v
		}
		s.data[k] = el
	}
}

func (s *NodeSet) Add(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	el, exist := s.data[key]
	if !exist {
		el = NodeSetElement{}
	}

	el.Value = value
	el.Add = time.Now()

	s.data[key] = el
}

func (s *NodeSet) Remove(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	el, exist := s.data[key]
	if !exist {
		el = NodeSetElement{}
	}

	el.Remove = time.Now()

	s.data[key] = el
}

func (s *NodeSet) Select(key string) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	el, exist := s.data[key]

	if exist && el.Add.After(el.Remove) {
		return el.Value
	}

	return nil
}

func (s *NodeSet) Data() map[string]NodeSetElement {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.data
}

func (s *NodeSet) String() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	buf := &bytes.Buffer{}
	for k, v := range s.data {
		buf.WriteString(fmt.Sprintf("Key: %s; ", k))
		buf.WriteString(v.String())
		buf.WriteByte('\n')
	}

	return buf.String()
}
