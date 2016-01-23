package network

import (
	"github.com/hsheth2/water"
	"os/exec"
	"fmt"
)

type InterfaceHandler struct {

}

func CreateTunInterface (iface string) (error, *water.Interface) (
	return water.NewTUN(iface)
)

func CreateTunInterfaceWithIp(iface string, IpAddr string) (error, *water.Interface) {
	iface, err := CreateTunInterface(iface)
	if err != nil {
		fmt.Println(error)
	}
	err = AssignIpAddress(iface, IpAddr)
	if err != nil {
		fmt.Println(error)
	}
	err = UpInterface(iface, IpAddr)

	if err != nil {
		fmt.Println(error)
	}
}
func AssignIpAddress (iface string, IpAddr string) error{
	err := exec.Command("ip", "addr", "add", IpAddr, "dev", iface).Run()
	return err
}

func UpInterface(iface string, IpAddr string) error {
	err := exec.Command("ip", "addr", "add", IpAddr, "dev", iface).Run()
	return err
}

func PacketsHandler() {
	// TODO: Make handler
}