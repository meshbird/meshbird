package iface

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

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
