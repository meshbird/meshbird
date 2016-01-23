package common

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nictuku/dht"
)

type DiscoveryDHT struct {
	BaseService

	node      *dht.DHT
	ih        dht.InfoHash
	localNode *LocalNode
	waitGroup sync.WaitGroup
	stopChan  chan bool
}

func (d DiscoveryDHT) Name() string {
	return "discovery-dht"
}

func (d *DiscoveryDHT) Init(ln *LocalNode) error {
	d.localNode = ln
	d.stopChan = make(chan bool)
	return nil
}

func (d *DiscoveryDHT) Run() error {
	var err error

	d.ih, err = dht.DecodeInfoHash(d.localNode.Config().NetworkID)
	if err != nil {
		return fmt.Errorf("decode infohash err: %s", err)
	}
	config := dht.NewConfig()
	config.Port = d.localNode.State().ListenPort
	d.node, err = dht.New(config)
	if err != nil {
		return fmt.Errorf("new dht init err: %s", err)
	}
	if err = d.node.Start(); err != nil {
		return fmt.Errorf("dht start err: %s", err)
	}
	d.waitGroup.Add(2)
	go d.process()
	go d.awaitPeers()
	d.waitGroup.Wait()
	return nil
}

func (d *DiscoveryDHT) Stop() {
	d.SetStatus(StatusStopping)
	d.stopChan <- true
}

func (d *DiscoveryDHT) process() {
	defer d.node.Stop()
	defer d.waitGroup.Done()
	t := time.NewTimer(60 * time.Second)
	defer t.Stop()
	log.Printf("dht request")
	d.node.PeersRequest(string(d.ih), true)
	for d.Status() != StatusStopping {
		select {
		case <-t.C:
			log.Printf("dht request")
			d.node.PeersRequest(string(d.ih), true)
		case <-d.stopChan:
			return
		}
	}
}

func (d *DiscoveryDHT) awaitPeers() {
	defer d.waitGroup.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for d.Status() != StatusStopping {
		select {
		case r := <-d.node.PeersRequestResults:
			for _, peers := range r {
				for _, x := range peers {
					log.Printf("peer: %v\n", dht.DecodePeerAddress(x))
				}
			}
		case <-ticker.C:
		}
	}
}
