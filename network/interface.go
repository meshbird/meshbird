package network

import (
	"os/exec"

	"github.com/hsheth2/water"
)

type InterfaceHandler struct {
}

func CreateTunInterface(iface string) (*water.Interface, error) {
	ifce, err := water.NewTUN(iface)
	err = UpInterface(iface)
	return ifce, err
}

func CreateTunInterfaceWithIp(iface string, IpAddr string) (*water.Interface, error) {
	ifce, err := CreateTunInterface(iface)
	if err == nil {
		err = AssignIpAddress(iface, IpAddr)
	}
	return ifce, err
}
func AssignIpAddress(iface string, IpAddr string) error {
	err := exec.Command("ip", "addr", "add", IpAddr, "dev", iface).Run()
	return err
}

func UpInterface(iface string) error {
	err := exec.Command("ip", "link", "set", iface, "up").Run()
	return err
}

func PacketsHandler() {
	// TODO: Make handler
}
