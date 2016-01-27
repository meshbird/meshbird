package common

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/prestonTao/upnp"
	"time"
)

type UPnPService struct {
	BaseService

	mapping *upnp.Upnp
	port    int
	logger  *log.Logger
}

func (d UPnPService) Name() string {
	return "UPnP"
}

func (s *UPnPService) Init(ln *LocalNode) error {
	// TODO: Add prefix
	s.logger = log.New()
	s.logger = ln.config.Loglevel
	s.mapping = new(upnp.Upnp)
	s.port = ln.State().ListenPort + 1
	return nil
}

func (s *UPnPService) Run() error {
	for !s.IsNeedStop() {
		err := s.process()
		if err != nil {
			s.logger.WithError(err).Error()
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
	s.logger.WithField("port", s.port).Info("UPnP port mapping")
	if err := s.mapping.AddPortMapping(s.port, s.port, "UDP"); err == nil {
		s.logger.Debug("port mapping passed")
	} else {
		s.logger.WithError(err).Error("port mapping fail")
	}
	return nil
}
