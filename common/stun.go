package common

import (
	"fmt"
	log "github.com/mgutz/logxi/v1"

	"github.com/ccding/go-stun/stun"
	"os"
	"time"
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
	s.logger = log.NewLogger(log.NewConcurrentWriter(os.Stderr), "[stun] ")
	s.client = stun.NewClient()
	s.client.SetServerAddr(STUNAddress)
	return nil
}

func (s *STUNService) Run() error {
	for !s.IsNeedStop() {
		err := s.process()
		if err != nil {
			log.Error("stun err: %s", err)
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
		if s.logger.IsInfo() {
			s.logger.Info("family %d, ip %s, port %d", host.Family(), host.IP(), host.Port())

		}
	}
	return nil
}
