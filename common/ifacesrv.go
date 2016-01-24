package common

import (
	"log"

	"github.com/gophergala2016/meshbird/network"
	"github.com/miolini/water"
)

type InterfaceService struct {
	BaseService

	localnode *LocalNode
	instance  *water.Interface
}

func (is *InterfaceService) Name() string {
	return "iface"
}

func (is *InterfaceService) Init(ln *LocalNode) (err error) {
	is.instance, err = network.CreateTunInterface("meshbird0")
	if err != nil {
		return
	}
	return nil
}

func (is *InterfaceService) Run() error {
	for {
		buf := make([]byte, 1500)
		n, err := is.instance.Read(buf)
		if err != nil {
			return err
		}
		log.Printf("[iface] read packet %d bytes", n)
	}
	return nil
}

