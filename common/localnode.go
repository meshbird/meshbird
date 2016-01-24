package common

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"sync"
)

type LocalNode struct {
	Node

	secret       *NetworkSecret
	config    *Config
	state     *State

	mutex     sync.Mutex
	waitGroup sync.WaitGroup

	services  map[string]Service
}

func NewLocalNode(cfg *Config) (*LocalNode, error) {
	var err error
	n := new(LocalNode)

	n.secret, err = NetworkSecretUnmarshal(cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	n.config = cfg
	n.config.NetworkID = n.secret.InfoHash()
	n.state = NewState(n.secret)

	n.services = make(map[string]Service)

	n.services[NetTable{}.Name()] = &NetTable{}
	n.services[DiscoveryDHT{}.Name()] = &DiscoveryDHT{}
	//n.services[STUNService{}.Name()] = &STUNService{}
	//n.services[UPnPService{}.Name()] = &UPnPService{}
	return n, nil
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
