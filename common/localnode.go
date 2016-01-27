package common

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
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

	logger *log.Logger
}

func NewLocalNode(cfg *Config) (*LocalNode, error) {
	var err error
	n := new(LocalNode)

	// TODO: Add prefix
	n.logger = log.New()
	n.logger.Level = cfg.Loglevel

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
		n.logger.WithField("name", service.Name()).Info("Initializing service...")
		if err := service.Init(n); err != nil {
			return fmt.Errorf("Initialision of %s finished with error: %s", service.Name(), err)
		}
		n.waitGroup.Add(1)
		go func(srv Service) {
			defer n.waitGroup.Done()
			n.logger.WithField("name", srv.Name()).Info("Running service...")
			if err := srv.Run(); err != nil {
				n.logger.WithFields(log.Fields{"name": srv.Name(), "err": err}).Error()
			}
		}(service)
	}
	return nil
}

func (n *LocalNode) Service(name string) Service {
	service, ok := n.services[name]
	if !ok {
		n.logger.WithField("name", name).Fatal("Service not found")
	}
	return service
}

func (n *LocalNode) WaitStop() {
	n.waitGroup.Wait()
}

func (n *LocalNode) Stop() error {
	n.logger.Info("Closing up local node")
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
		panic("net-table not found")
	}
	return service.(*NetTable)
}
