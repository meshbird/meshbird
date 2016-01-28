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
)

type State struct {
	Secret     *secure.NetworkSecret `json:"-"`
	ListenPort int                   `json:"port"`
	PrivateIP  net.IP                `json:"private_ip"`
}

func NewState(secret *secure.NetworkSecret) *State {
	s := &State{
		Secret: secret,
	}
	s.Load()

	var save bool

	if s.ListenPort < 1 {
		s.ListenPort = GetRandomPort()
		save = true
	}
	if s.PrivateIP == nil {
		var err error
		if s.PrivateIP, err = network.GenerateIPAddress(secret.Net); err == nil {
			save = true
		} else {
			log.Error("error on generate IP, %v", err)
		}
	}

	if save {
		s.Save()
	}

	return s
}

func (s *State) Load() {
	if data, err := ioutil.ReadFile(s.getConfigPath()); err == nil {
		if err = json.Unmarshal(data, s); err == nil {
			log.Info("state restored, %+v", s)
		}
	}
}

func (s *State) Save() {
	if data, err := json.Marshal(s); err == nil {
		if err = ioutil.WriteFile(s.getConfigPath(), data, os.ModePerm); err != nil {
			log.Error("error on write state, %v", err)
		}
	}
}

func (s *State) getConfigPath() string {
	return path.Join(os.Getenv("HOME"), fmt.Sprintf(".meshbird_%s.json", s.Secret.InfoHash()))
}
