package common

import (
	"fmt"
	"github.com/meshbird/meshbird/log"
	"github.com/miolini/upnp"
	"time"
)

type UPnPService struct {
	BaseService

	mapping *upnp.Upnp
	port    int
	logger  log.Logger
}

func (s UPnPService) Name() string {
	return "UPnP"
}

func (s *UPnPService) Init(ln *LocalNode) error {
	s.logger = log.L(s.Name())
	s.mapping = new(upnp.Upnp)
	s.port = ln.State().ListenPort() + 1
	return nil
}

func (s *UPnPService) Run() error {
	for !s.IsNeedStop() {
		err := s.process()
		if err != nil {
			s.logger.Error("error on process, %v", err)
		}
		time.Sleep(time.Minute)
	}
	return nil
}

func (s *UPnPService) process() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %s", r)
		}
	}()
	s.logger.Debug("trying to map port %d...", s.port)
	if err := s.mapping.AddPortMapping(s.port, s.port, "UDP"); err == nil {
		s.logger.Debug("port mapping passed")
	} else {
		s.logger.Debug("port mapping fail, %v", err)
	}
	return nil
}
