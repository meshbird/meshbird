package common

import (
	"fmt"
	"github.com/gophergala2016/meshbird/secure"
	"log"
	"sync"
)

type LocalNode struct {
	Node

	secret *secure.NetworkSecret
	config *Config
	state  *State

	mutex     sync.Mutex
	waitGroup sync.WaitGroup

	services map[string]Service
}

func NewLocalNode(cfg *Config) (*LocalNode, error) {
	var err error
	n := new(LocalNode)

	n.secret, err = secure.NetworkSecretUnmarshal(cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	n.config = cfg
	n.config.NetworkID = n.secret.InfoHash()
	n.state = NewState(n.secret)

	n.services = make(map[string]Service)

	n.AddService(&NetTable{})
	n.AddService(&ListenerService{})
	n.AddService(&DiscoveryDHT{})
	n.AddService(&InterfaceService{})
	n.AddService(&STUNService{})
	n.AddService(&UPnPService{})

	return n, nil
}

func (n *LocalNode) Config() Config {
	return *n.config
}

func (n *LocalNode) State() State {
	return *n.state
}

func (n *LocalNode) AddService(srv Service) {
	n.services[srv.Name()] = srv
}

func (n *LocalNode) Start() error {
	for name, service := range n.services {
		log.Printf("[%s] service init", name)
		if err := service.Init(n); err != nil {
			return fmt.Errorf("[%s] service init err: %s", service.Name(), err)
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

func (n *LocalNode) Service(name string) Service {
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

func (n *LocalNode) NetworkSecret() *secure.NetworkSecret {
	return n.secret
}
