package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type State struct {
	ListenPort int    `json:"port"`
	ListenHost string `json:"host"`
}

func NewState() *State {
	s := &State{}
	s.Load()
	return s
}

func (s *State) Load() {
	if data, err := ioutil.ReadFile(s.getConfigPath()); err == nil {
		json.Unmarshal(data, s)
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
	return path.Join(os.Getenv("HOME"), ".meshbird.json")
}
