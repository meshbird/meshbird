package common

import (
	"encoding/json"
	"fmt"
	"github.com/gophergala2016/meshbird/network"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type State struct {
	ListenPort int    `json:"port"`
	ListenHost string `json:"host"`
	netId      string `json:"-"`
}

func NewState(net, netId string) *State {
	s := &State{
		netId: netId,
	}
	s.Load()

	var save bool

	if s.ListenPort < 1 {
		s.ListenPort = GetRandomPort()
		save = true
	}
	if s.ListenHost == "" {
		var err error
		if s.ListenHost, err = network.GenerateIPAddress(net); err == nil {
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
			log.Printf("State restored: %v", s)
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
	return path.Join(os.Getenv("HOME"), fmt.Sprintf(".meshbird_%s.json", s.netId))
}
