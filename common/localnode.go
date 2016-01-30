package common

import (
	"fmt"
	"github.com/meshbird/meshbird/log"
	"github.com/meshbird/meshbird/secure"
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

	logger log.Logger
}

func NewLocalNode(cfg *Config) (*LocalNode, error) {
	var err error
	n := new(LocalNode)
	n.logger = log.L("local")

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
	//n.AddService(&STUNService{})
	n.AddService(&UPnPService{})
	n.AddService(&HttpService{})
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
	for _, service := range n.services {
		n.logger.Info("initializing %s service", service.Name())
		if err := service.Init(n); err != nil {
			return fmt.Errorf("initialization of %q finished with error: %v", service.Name(), err)
		}
		n.waitGroup.Add(1)

		go func(srv Service) {
			defer n.waitGroup.Done()
			n.logger.Info("running %q service", srv.Name())
			if err := srv.Run(); err != nil {
				n.logger.Error("error on run %q service, %v", srv.Name(), err)
			}
		}(service)
	}
	return nil
}

func (n *LocalNode) Service(name string) Service {
	service, ok := n.services[name]
	if !ok {
		n.logger.Fatal("service %q not found", name)
	}
	return service
}

func (n *LocalNode) WaitStop() {
	n.waitGroup.Wait()
}

func (n *LocalNode) Stop() error {
	n.logger.Info("stopping...")
	for _, service := range n.services {
		service.Stop()
	}
	return nil
}

func (n *LocalNode) NetworkSecret() *secure.NetworkSecret {
	return n.secret
}

func (n *LocalNode) NetTable() *NetTable {
	service, ok := n.services["net-table"]
	if !ok {
		n.logger.Panic("net-table not found")
	}
	return service.(*NetTable)
}
