package common

import (
	"fmt"
	"time"

	"github.com/ccding/go-stun/stun"
	"github.com/meshbird/meshbird/log"
)

const (
	STUNAddress = "stun.l.google.com:19302"
)

type STUNService struct {
	BaseService

	client *stun.Client
	logger log.Logger
}

func (d STUNService) Name() string {
	return "STUN"
}

func (s *STUNService) Init(ln *LocalNode) error {
	s.logger = log.L(s.Name())
	s.client = stun.NewClient()
	s.client.SetServerAddr(STUNAddress)
	return nil
}

func (s *STUNService) Run() error {
	for !s.IsNeedStop() {
		err := s.process()
		if err != nil {
			s.logger.Error("error on process, %v", err)
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
	case stun.NATError:
		return fmt.Errorf("test failed")
	case stun.NATUnknown:
		return fmt.Errorf("unexpected response from the STUN server")
	case stun.NATBlocked:
		return fmt.Errorf("UDP is blocked")
	case stun.NATFull:
		return fmt.Errorf("full cone NAT")
	case stun.NATSymetric:
		return fmt.Errorf("symetric NAT")
	case stun.NATRestricted:
		return fmt.Errorf("restricted NAT")
	case stun.NATPortRestricted:
		return fmt.Errorf("port restricted NAT")
	case stun.NATNone:
		return fmt.Errorf("not behind a NAT")
	case stun.NATSymetricUDPFirewall:
		return fmt.Errorf("symetric UDP firewall")
	}

	if host != nil {
		s.logger.Info("processed, family %d, host %q, port %d", host.Family(), host.IP(), host.Port())
	}
	return nil
}
