package common

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nictuku/dht"
)

type LocalNode struct {
	Node
	Statusable

	config *Config
	state  *State

	mutex     sync.Mutex
	waitGroup sync.WaitGroup
}

func NewLocalNode(cfg *Config) *LocalNode {
	n := new(LocalNode)
	return n
}

func (n *LocalNode) Run() error {
	n.SetStatus(1)
	return nil
}

func (n *LocalNode) Stop() error {
	n.waitGroup.Wait()
	return nil
}

func (n *LocalNode) discovery(networkID string, port int) error {
	ih, err := dht.DecodeInfoHash(networkID)
	if err != nil {
		return fmt.Errorf("decode infohash err: %s", err)
	}

	config := dht.NewConfig()
	config.Port = port
	d, err := dht.New(config)
	if err != nil {
		return fmt.Errorf("new dht init err: %s", err)
	}

	if err = d.Start(); err != nil {
		return fmt.Errorf("dht start err: %s", err)
	}
	defer d.Stop()
	go n.discoveryAwait(d)

	log.Printf("peer request to DHT")
	for {
		log.Printf("dht request")
		d.PeersRequest(string(ih), true)
		time.Sleep(time.Second * 60)
	}
	return nil
}

func (n *LocalNode) discoveryAwait(dhtNetwork *dht.DHT) {
	defer log.Printf("discovery await exit")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case r := <-dhtNetwork.PeersRequestResults:
			for _, peers := range r {
				for _, x := range peers {
					log.Printf("peer: %v\n", dht.DecodePeerAddress(x))
				}
			}
		case <-ticker.C:
		}
		if n.Status() > 1 {
			break
		}
	}
}

func hashSecretKey(key string) string {
	hashBytes := sha1.Sum([]byte(key))
	return hex.EncodeToString(hashBytes[:])
}
