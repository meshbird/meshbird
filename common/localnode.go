package common

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/gophergala2016/meshbird/ecdsa"
	"log"
	"sync"
)

type LocalNode struct {
	Node

	config *Config
	state  *State

	mutex     sync.Mutex
	waitGroup sync.WaitGroup

	services map[string]Service
}

func NewLocalNode(cfg *Config) *LocalNode {
	key := ecdsa.Unpack([]byte(cfg.SecretKey))

	n := new(LocalNode)
	n.config = cfg
	n.config.NetworkID = ecdsa.HashSecretKey(n.config.SecretKey)
	n.state = NewState(key.CIDR, n.config.NetworkID)

	n.services = make(map[string]Service)

	n.services[NetTable{}.Name()] = &NetTable{}
	n.services[DiscoveryDHT{}.Name()] = &DiscoveryDHT{}
	//n.services[STUNService{}.Name()] = &STUNService{}
	//n.services[UPnPService{}.Name()] = &UPnPService{}
	return n
}

func (n *LocalNode) Config() Config {
	return *n.config
}

func (n *LocalNode) State() State {
	return *n.state
}

func (n *LocalNode) Start() error {
	for name, service := range n.services {
		log.Printf("[%s] service init", name)
		if err := service.Init(n); err != nil {
			return err
		}
		n.waitGroup.Add(1)
		go func(srv Service) {
			defer n.waitGroup.Done()
			log.Printf("[%s] service run", srv.Name())
			if err := srv.Run(); err != nil {
				log.Printf("[%s] error: %s", srv.Name(), err)
			}
		}(service)
	}
	return nil
}

func (n *LocalNode) GetService(name string) Service {
	service, ok := n.services[name]
	if !ok {
		log.Panicf("service %s not found", name)
	}
	return service
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
