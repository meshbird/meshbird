package common

import (
	log "github.com/Sirupsen/logrus"
	"time"

	"fmt"
	"github.com/prestonTao/upnp"
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
	// TODO: FIX IT
	// s.logger = log.NewLogger(log.NewConcurrentWriter(os.Stderr), "[upnp] ")
	s.logger = log.New()
	s.mapping = new(upnp.Upnp)
	s.port = ln.State().ListenPort + 1
	return nil
}

func (s *UPnPService) Run() error {
	for !s.IsNeedStop() {
		err := s.process()
		if err != nil {
			log.Error("upnp err: %s", err)
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
	s.logger.Info(fmt.Sprintf("UPnP port mapping: %d", s.port))
	if err := s.mapping.AddPortMapping(s.port, s.port, "UDP"); err == nil {
		s.logger.Debug("port mapping passed")
	} else {
		s.logger.Debug("port mapping fail")
	}
	return nil
}
