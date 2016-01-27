package network

import (
	"fmt"
	"os/exec"
	"strconv"
)

func interfaceOpen(ifType, ifName string) (*Interface, error) {
	var err error
	if ifType != "tun" && ifType != "tap" {
		return nil, fmt.Errorf("unknown interface type: %s", ifType)
	}
	ifce := new(Interface)
	for i := 0; i < 256; i++ {
		ifPath := fmt.Sprintf("/dev/tun/%s%d", ifType, ifName)
		ifce.file, err = os.OpenFile(ifPath, os.O_RDWR, 0644)
		if err != nil {
			continue
		}
		ifce.name = ifName
	}
	if ifce.file == nil {
		return nil, fmt.Errorf("can't create network interface")
	}
	return ifce, err
}

func AssignIpAddress(iface string, IpAddr string) error {
	err := exec.Command("ip", "addr", "add", IpAddr, "dev", iface).Run()
	if err != nil {
		return fmt.Errorf("assign ip %s to %s err: %s", IpAddr, iface, err)
	}
	return err
}

func UpInterface(iface string) error {
	err := exec.Command("ip", "link", "set", iface, "up").Run()
	if err != nil {
		return fmt.Errorf("up interface %s err: %s", iface, err)
	}
	return err
}

func SetMTU(iface string, mtu int) error {
	err := exec.Command("ip", "link", "set", "mtu", strconv.Itoa(mtu), "dev", iface).Run()
	if err != nil {
		return fmt.Errorf("Can't set MTU %s to %s err: %s", iface, strconv.Itoa(mtu), err)
	}
	return nil
}
