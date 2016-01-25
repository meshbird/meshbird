package common

import (
	"fmt"
	"github.com/gophergala2016/meshbird/network"
	"github.com/miolini/water"
	"log"
	"os"
	"strconv"
	"net"
)

type InterfaceService struct {
	BaseService

	localnode *LocalNode
	instance  *water.Interface
	netTable  *NetTable
	logger    *log.Logger
}

func (is *InterfaceService) Name() string {
	return "iface"
}

func (is *InterfaceService) Init(ln *LocalNode) (err error) {
	is.logger = log.New(os.Stderr, "[iface] ", log.LstdFlags)
	is.localnode = ln
	is.netTable = ln.NetTable()
	netsize, _ := ln.State().Secret.Net.Mask.Size()
	IPAddr := fmt.Sprintf("%s/%s", ln.State().PrivateIP, strconv.Itoa(netsize))
	is.instance, err = network.CreateTunInterfaceWithIp("", IPAddr)
	if err != nil {
		return fmt.Errorf("create interface %s err: %s", "", err)
	}
	tunIface, err := net.InterfaceByName(is.instance.Name())
	if err != nil {
		fmt.Println(err)
	}
	tunIface.MTU = 1400

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

		dst := network.IPv4Destination(packet)
		is.netTable.SendPacket(dst, packet)
		is.logger.Printf("Read packet %d bytes", n)
	}
	return nil
}

func (is *InterfaceService) WritePacket(packet []byte) {
	is.logger.Printf("Package for writing received, length %d bytes", len(packet))
	if _, err := is.instance.Write(packet); err != nil {
		is.logger.Printf("Error on twite packet: %v", err)
	}
}
