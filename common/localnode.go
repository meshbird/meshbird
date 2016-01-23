package common

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"sync"
)

type LocalNode struct {
	Node

	config *Config
	state  *State

	mutex     sync.Mutex
	waitGroup sync.WaitGroup

	services []Service
}

func NewLocalNode(cfg *Config) *LocalNode {
	n := new(LocalNode)
	n.config = cfg
	n.state = &State{}
	n.services = append(n.services, &DiscoveryDHT{})
	//n.services = append(n.services, &STUNService{})
	//n.services = append(n.services, &UPnPService{})
	return n
}

func (n *LocalNode) Config() Config {
	return *n.config
}

func (n *LocalNode) State() State {
	return *n.state
}

func (n *LocalNode) Start() error {
	serviceCounter := 0
	for _, service := range n.services {
		log.Printf("[%s] service init", service.Name())
		err := service.Init(n)
		if err != nil {
			log.Printf("[%s] init error: %s", service.Name(), err)
			continue
		}
		serviceCounter++
	}
	n.waitGroup.Add(serviceCounter)
	for _, service := range n.services {
		go func() {
			defer n.waitGroup.Done()
			log.Printf("[%s] service run", service.Name())
			err := service.Run()
			if err != nil {
				log.Printf("[%s] error: %s", service.Name(), err)
			}
		}()
	}
	return nil
}

func (n *LocalNode) WaitStop() {
	n.waitGroup.Wait()
}

func (n *LocalNode) Stop() error {
	log.Printf("closing up local node")
	for _, service := range n.services {
		service.Stop()
	}
	return nil
}

func hashSecretKey(key string) string {
	hashBytes := sha1.Sum([]byte(key))
	return hex.EncodeToString(hashBytes[:])
}
