package common

import (
	"sync"
"log"
)

type NetTable struct {
	BaseService
	localNode *LocalNode
	waitGroup sync.WaitGroup
	dhtInChan chan string
}

func (nt NetTable) Name() string {
	return "net-table"
}

func (nt *NetTable) Init(ln *LocalNode) error {
	nt.localNode = ln
	nt.dhtInChan = make(chan string)
	return nil
}

func (nt *NetTable) Run() error {
	nt.waitGroup.Add(1)
	go nt.processDHTIn()
	nt.waitGroup.Wait()
	return nil
}

func (nt *NetTable) Stop() {
	nt.SetStatus(StatusStopping)
}

func (nt *NetTable) GetDHTInChannel() chan<- string {
	return nt.dhtInChan
}

func (nt *NetTable) processDHTIn() {
	defer nt.waitGroup.Done()

	for nt.Status() != StatusStopping {
		select {
		case host, ok := <-nt.dhtInChan:
			if !ok {
				return
			}
			log.Printf("Got %s host from DHT", host)
		}
	}
}
