package common

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/nictuku/dht"
	"os"
	"sync"
	"time"
)

type DiscoveryDHT struct {
	BaseService

	node      *dht.DHT
	ih        dht.InfoHash
	localNode *LocalNode
	waitGroup sync.WaitGroup
	stopChan  chan bool

	lastPeers []string
	mutex     sync.Mutex

	logger *log.Logger
}

func (d DiscoveryDHT) Name() string {
	return "discovery-dht"
}

func (d *DiscoveryDHT) Init(ln *LocalNode) error {
	// TODO: Add prefix
	d.logger = log.New()
	d.logger.Level = ln.config.Loglevel
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

	d.request()

	for d.Status() != StatusStopping {
		select {
		case <-t.C:
			d.request()
		case <-d.stopChan:
			return
		}
	}
}

func (d *DiscoveryDHT) request() {
	d.logger.Debug("Request...")
	d.node.PeersRequest(string(d.ih), true)
}

func (d *DiscoveryDHT) addPeer(peer string) {
	d.mutex.Lock()
	exists := false
	for _, lastPeer := range d.lastPeers {
		if lastPeer == peer {
			exists = true
			break
		}
	}
	if !exists {
		d.lastPeers = append(d.lastPeers, peer)
		if len(d.lastPeers) > 1000 {
			d.lastPeers = d.lastPeers[1:]
		}
	}
	d.mutex.Unlock()
	if exists {
		return
	}

	d.logger.WithField("peer", peer).Debug("New peer received")
	d.localNode.NetTable().GetDHTInChannel() <- peer
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
					host := dht.DecodePeerAddress(x)
					d.addPeer(host)
				}
			}
		case <-ticker.C:
		}
	}
}
