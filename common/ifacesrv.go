package common

import (
	"fmt"
	"github.com/meshbird/meshbird/log"
	"github.com/meshbird/meshbird/network"
	"strconv"
)

type InterfaceService struct {
	BaseService

	ln       *LocalNode
	instance *network.Interface
	netTable *NetTable
	logger   log.Logger
}

func (is *InterfaceService) Name() string {
	return "iface"
}

func (is *InterfaceService) Init(ln *LocalNode) (err error) {
	is.logger = log.L(is.Name())
	is.ln = ln
	is.netTable = ln.NetTable()
	netSize, _ := ln.State().Secret.Net.Mask.Size()
	IPAddress := fmt.Sprintf("%s/%s", ln.State().PrivateIP, strconv.Itoa(netSize))

	if is.instance, err = network.CreateTunInterfaceWithIp("", IPAddress); err != nil {
		return err
	}

	if err = network.SetMTU(is.instance.Name(), 1400); err != nil {
		is.logger.Warning("unable to set mtu, %v", err)
	}
	return nil
}

func (is *InterfaceService) Run() error {
	for {
		buf := make([]byte, 1500)
		n, err := is.instance.Read(buf)
		if err != nil {
			is.logger.Error("error on read from interface, %v", err)
			return err
		}
		packet := buf[:n]

		dst := network.IPv4Destination(packet)
		is.netTable.SendPacket(dst, packet)
		is.logger.Debug("successfully been read %d bytes", n)
	}
	return nil
}

func (is *InterfaceService) WritePacket(packet []byte) error {
	is.logger.Debug("ready to write %d bytes", len(packet))
	if _, err := is.instance.Write(packet); err != nil {
		return err
	}
	return nil
}
