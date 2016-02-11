package common

import (
	"encoding/json"
	"fmt"
	"github.com/meshbird/meshbird/log"
	"github.com/meshbird/meshbird/network"
	"github.com/meshbird/meshbird/secure"
	"io/ioutil"
	"net"
	"os"
	"path"
	"sync"
)

type State struct {
	secret     *secure.NetworkSecret `json:"-"`
	listenPort int                   `json:"port"`
	privateIP  net.IP                `json:"private_ip"`
	goodPeers  []string              `json:"good_peers"`
	mutex      sync.Mutex            `json:"-"`
}

func NewState(secret *secure.NetworkSecret) *State {
	s := &State{
		secret: secret,
	}

	if err := s.Load(); err != nil {
		log.Debug("state load err: %s", err) // probably no file
	}

	var save bool

	if s.listenPort < 1 {
		s.listenPort = GetRandomPort()
		save = true
	}
	if s.privateIP == nil {
		var err error
		if s.privateIP, err = network.GenerateIPAddress(secret.Net); err == nil {
			save = true
		} else {
			log.Error("error on generate IP, %v", err)
		}
	}

	if save {
		if err := s.Save(); err != nil {
			log.Error("state save err: %s", err)
		}
	}

	return s
}

func (s *State) Secret() *secure.NetworkSecret {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.secret
}

func (s *State) ListenPort() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.listenPort
}

func (s *State) PrivateIP() net.IP {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.privateIP
}

func (s *State) GoodPeers() []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.goodPeers
}

func (s *State) AddGoodPeer(peer string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.goodPeers = append(s.goodPeers, peer)
	go func() {
		if err := s.Save(); err != nil {
			log.Error("state save err: %s", err)
		}
	}()
}

func (s *State) Load() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	data, err := ioutil.ReadFile(s.getConfigPath())
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, s); err == nil {
		return err
	}
	return nil
}

func (s *State) Save() error {
	log.Debug("state saving")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	data, err := json.Marshal(s)
	if err != nil {
		log.Error("%s", err)
		return err
	}
	log.Debug("state content: %#v", s)
	log.Debug("state json content: %s", string(data))
	if err = ioutil.WriteFile(s.getConfigPath(), data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (s *State) getConfigPath() string {
	return path.Join(os.Getenv("HOME"), fmt.Sprintf(".meshbird_%s.json", s.secret.InfoHash()))
}
