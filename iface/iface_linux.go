package iface

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"

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
	log.Printf("iface name: %s", i.Name())
	log.Printf("ip: %s", ip.String())
	cmd := exec.Command("ifconfig", i.Name(),
		ip.String(), "netmask", netmask,
		"mtu", strconv.Itoa(i.mtu), "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ifconfig run err: %s %s", err, string(output))
	}
	return nil
}
