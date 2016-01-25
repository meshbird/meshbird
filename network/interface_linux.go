package network

import (
	"os"
	"strconv"
)

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
