package common

import (
	"fmt"
	"github.com/meshbird/meshbird/network"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

type InterfaceService struct {
	BaseService

	localnode *LocalNode
	instance  *network.Interface
	netTable  *NetTable
	logger    *log.Logger
}

func (is *InterfaceService) Name() string {
	return "iface"
}

func (is *InterfaceService) Init(ln *LocalNode) (err error) {
	is.logger = log.New()
	is.localnode = ln
	is.netTable = ln.NetTable()
	netsize, _ := ln.State().Secret.Net.Mask.Size()
	IPAddr := fmt.Sprintf("%s/%s", ln.State().PrivateIP, strconv.Itoa(netsize))
	is.instance, err = network.CreateTunInterfaceWithIp("", IPAddr)
	if err != nil {
		return fmt.Errorf("create interface %s err: %s", "", err)
	}
	err = network.SetMTU(is.instance.Name(), 1400)

	if err != nil {
			is.logger.Warn(err.Error())
	}
	return nil
}

func (is *InterfaceService) Run() error {
	for {
		buf := make([]byte, 1500)
		n, err := is.instance.Read(buf)
		if err != nil {
			is.logger.Error(err.Error())
			return err
		}
		packet := buf[:n]

		dst := network.IPv4Destination(packet)
		is.netTable.SendPacket(dst, packet)
			is.logger.Debug(fmt.Sprintf("Read packet %d bytes", n))
	}
	return nil
}

func (is *InterfaceService) WritePacket(packet []byte) {
		is.logger.Debug(fmt.Sprintf("Package for writing received, length %d bytes\n", len(packet)))
	if _, err := is.instance.Write(packet); err != nil {
		is.logger.Error(fmt.Sprintf("Error on twite packet: %v", err))
	}
}
