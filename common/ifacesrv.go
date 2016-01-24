package common

import (
	"fmt"
	"github.com/gophergala2016/meshbird/network"
	"github.com/miolini/water"
	"log"
	"strconv"
)

type InterfaceService struct {
	BaseService

	localnode *LocalNode
	instance  *water.Interface
	netTable  *NetTable
}

func (is *InterfaceService) Name() string {
	return "iface"
}

func (is *InterfaceService) Init(ln *LocalNode) (err error) {
	is.localnode = ln
	is.netTable = ln.Service("net-table").(*NetTable)
	netsize, _ := ln.State().Secret.Net.Mask.Size()
	IPAddr := fmt.Sprintf("%s/%s", ln.State().PrivateIP, strconv.Itoa(netsize))
	is.instance, err = network.CreateTunInterfaceWithIp("", IPAddr)
	if err != nil {
		return fmt.Errorf("create interface %s err: %s", "", err)
	}
	return nil
}

func (is *InterfaceService) Run() error {
	for {
		buf := make([]byte, 1500)
		n, err := is.instance.Read(buf)
		if err != nil {
			log.Println(err)
			return err
		}
		packet := buf[:n]
		fmt.Println(len(packet))
		log.Printf("[iface] read packet %d bytes", n)
	}
	return nil
}
