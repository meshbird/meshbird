package common

import (
	"sync"
	"time"
)

type NetTable struct {
	BaseService

	localNode *LocalNode
	waitGroup sync.WaitGroup
	dhtInChan chan string

	lock      sync.RWMutex
	blackList map[string]time.Time
}

func (nt NetTable) Name() string {
	return "net-table"
}

func (nt *NetTable) Init(ln *LocalNode) error {
	nt.localNode = ln
	nt.dhtInChan = make(chan string, 10)
	nt.blackList = make(map[string]time.Time)
	return nil
}

func (nt *NetTable) Run() error {
	for i := 0; i < 10; i++ {
		go nt.processDHTIn()
	}
	return nil
}

func (nt *NetTable) Stop() {
	nt.SetStatus(StatusStopping)
}

func (nt *NetTable) GetDHTInChannel() chan<- string {
	return nt.dhtInChan
}

func (nt *NetTable) processDHTIn() {
	for nt.Status() != StatusStopping {
		select {
		case host, ok := <-nt.dhtInChan:
			if !ok {
				return
			}
			nt.lock.Lock()
			_, ok = nt.blackList[host]
			nt.lock.Unlock()

			if !ok {
				nt.tryConnect(host)
			}
		}
	}
}

func (nt *NetTable) tryConnect(h string) {
	rn, err := TryConnect(h, nt.localNode.NetworkSecret())
	if err != nil {
		nt.addToBlackList(h)
		return
	}
	if rn == nil {
		return
	}
}

func (nt *NetTable) addToBlackList(h string) {
	nt.lock.Lock()
	defer nt.lock.Unlock()
	nt.blackList[h] = time.Now()
}
