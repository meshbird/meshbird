package common

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/meshbird/meshbird/network"
	"os"
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
	// TODO: Add prefix to logs
	is.logger = log.New()
	is.logger = ln.config.Loglevel
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
		is.logger.WithError(err).Warn()
	}
	return nil
}

func (is *InterfaceService) Run() error {
	for {
		buf := make([]byte, 1500)
		n, err := is.instance.Read(buf)
		if err != nil {
			is.logger.WithError(err).Error()
			return err
		}
		packet := buf[:n]

		dst := network.IPv4Destination(packet)
		is.netTable.SendPacket(dst, packet)
		is.logger.WithField("len", n).Debug("Read packet")
	}
	return nil
}

func (is *InterfaceService) WritePacket(packet []byte) {
	is.logger.WithField("len", len(packet)).Debug("Package for writing received")
	if _, err := is.instance.Write(packet); err != nil {
		is.logger.WithError(err).Error("Error on twite packet")
	}
}
