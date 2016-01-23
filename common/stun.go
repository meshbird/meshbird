package common

import (
	"fmt"
	"log"

	"github.com/ccding/go-stun/stun"
	"time"
)

const (
	STUNAddress = "stun.l.google.com:19302"
)

type STUNService struct {
	BaseService

	client *stun.Client
}

func (s *STUNService) Init(ln *LocalNode) error {
	s.client = stun.NewClient()
	s.client.SetServerAddr(STUNAddress)
	return nil
}

func (s *STUNService) Run() error {
	for !s.IsNeedStop() {
		err := s.process()
		if err != nil {
			log.Printf("stun err: %s", err)
		}
		time.Sleep(time.Minute)
	}
	return nil
}

func (s *STUNService) process() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %s", r)
		}
	}()
	nat, host, err := s.client.Discover()
	if err != nil {
		return err
	}
	switch nat {
	case stun.NAT_ERROR:
		return fmt.Errorf("Test failed")
	case stun.NAT_UNKNOWN:
		return fmt.Errorf("Unexpected response from the STUN server")
	case stun.NAT_BLOCKED:
		return fmt.Errorf("UDP is blocked")
	case stun.NAT_FULL:
		return fmt.Errorf("Full cone NAT")
	case stun.NAT_SYMETRIC:
		return fmt.Errorf("Symetric NAT")
	case stun.NAT_RESTRICTED:
		return fmt.Errorf("Restricted NAT")
	case stun.NAT_PORT_RESTRICTED:
		return fmt.Errorf("Port restricted NAT")
	case stun.NAT_NONE:
		return fmt.Errorf("Not behind a NAT")
	case stun.NAT_SYMETRIC_UDP_FIREWALL:
		return fmt.Errorf("Symetric UDP firewall")
	}

	if host != nil {
		log.Printf("family %d, ip %s, port %d", host.Family(), host.IP(), host.Port())
	}
	return nil
}
