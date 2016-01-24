package common

import (
	"encoding/json"
	"fmt"
	"github.com/gophergala2016/meshbird/network"
	"github.com/gophergala2016/meshbird/secure"
	"io/ioutil"
	"log"
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
			log.Printf("Error on generate IP: %s", err)
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
			log.Printf("State restored: %+v, private IP: %x", s, s.PrivateIP)
		}
	}
}

func (s *State) Save() {
	if data, err := json.Marshal(s); err == nil {
		if err = ioutil.WriteFile(s.getConfigPath(), data, os.ModePerm); err != nil {
			log.Printf("Error on write state: %s", err)
		}
	}
}

func (s *State) getConfigPath() string {
	return path.Join(os.Getenv("HOME"), fmt.Sprintf(".meshbird_%s.json", s.Secret.InfoHash()))
}
