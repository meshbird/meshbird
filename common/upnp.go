package common

import (
	"log"
	"time"

	"fmt"
	"github.com/prestonTao/upnp"
)

type UPnPService struct {
	BaseService

	mapping *upnp.Upnp
	port    int
}

func (d UPnPService) Name() string {
	return "UPnP"
}

func (s *UPnPService) Init(ln *LocalNode) error {
	s.mapping = new(upnp.Upnp)
	s.port = ln.State().ListenPort + 1
	return nil
}

func (s *UPnPService) Run() error {
	for !s.IsNeedStop() {
		err := s.process()
		if err != nil {
			log.Printf("upnp err: %s", err)
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
	log.Printf("UPnP port mapping: %d", s.port)
	if err := s.mapping.AddPortMapping(s.port, s.port, "UDP"); err == nil {
		log.Printf("port mapping passed")
	} else {
		log.Printf("port mapping fail")
	}
	return nil
}
