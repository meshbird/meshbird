package common

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
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
	logger     *log.Logger           `json:"-"`
}

func NewState(secret *secure.NetworkSecret) *State {
	// TODO Add prefix
	s := &State{
		Secret: secret,
		logger: log.New(),
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
			s.logger.WithError(err).Error("Error on generate IP")
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
			s.logger.WithField("state", s).Info("State restored")
		}
	}
}

func (s *State) Save() {
	if data, err := json.Marshal(s); err == nil {
		if err = ioutil.WriteFile(s.getConfigPath(), data, os.ModePerm); err != nil {
			s.logger.WithError(err).Error("Error on write state")
		}
	}
}

func (s *State) getConfigPath() string {
	return path.Join(os.Getenv("HOME"), fmt.Sprintf(".meshbird_%s.json", s.Secret.InfoHash()))
}
