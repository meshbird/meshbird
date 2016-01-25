package network

import (
	"fmt"
	"os/exec"
	"strconv"
)

func AssignIpAddress(iface string, IpAddr string) error {
	err := exec.Command("ifconfig", iface, IpAddr).Run()
	if err != nil {
		return fmt.Errorf("assign ip %s to %s err: %s", IpAddr, iface, err)
	}
	return err
}

func UpInterface(iface string) error {
	err := exec.Command("ifconfig", iface, "up").Run()
	return err
}

func SetMTU(iface string, mtu int) error {
	err := exec.Command("ifconfig", iface, "mtu", strconv.Itoa(mtu)).Run()
	if err != nil {
		return fmt.Errorf("Can't set MTU %s to %s err: %s", iface, strconv.Itoa(mtu), err)
	}
	return nil
}
