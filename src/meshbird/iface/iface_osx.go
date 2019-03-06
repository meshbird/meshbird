package iface

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

type Iface struct {
	name string
	ip   string
	mtu  int
	ifce *water.Interface
}

func New(name, ip string, mtu int) *Iface {
	return &Iface{
		name: name,
		ip:   ip,
		mtu:  mtu,
	}
}

func (i *Iface) Start() error {
	ip, netIP, err := net.ParseCIDR(i.ip)
	if err != nil {
		return fmt.Errorf("parse IP err: %s", err)
	}
	config := water.Config{
		DeviceType: water.TUN,
	}
	//config.Name = i.name
	i.ifce, err = water.New(config)
	if err != nil {
		return fmt.Errorf("iface create err: %s", err)
	}
	mask := netIP.Mask
	netmask := fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
	cmd := exec.Command("ifconfig", i.Name(), "inet", ip.String(), netmask, "mtu", fmt.Sprintf("%d", i.mtu))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("iface run err: %s %s", err, string(output))
	}
	return nil
}

func (i *Iface) Name() string {
	return i.ifce.Name()
}

func (i *Iface) Read(pkt PacketIP) (int, error) {
	return i.ifce.Read(pkt)
}

func (i *Iface) Write(pkt PacketIP) (int, error) {
	return i.ifce.Write(pkt)
}
