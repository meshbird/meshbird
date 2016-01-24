package network

import (
	"os/exec"

	"github.com/miolini/water"
)

const DEFAULT_MTU = 1500

var MTU int

func Init() {
	MTU = 0
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
	err := exec.Command("ifconfig", iface, IpAddr).Run()
	return err
}

func UpInterface(iface string) error {
	err := exec.Command("ifconfig", iface, "up").Run()
	return err
}
func SetMTU(mtu int) {
	MTU = mtu
}

func NextNetworkPacket(iface *water.Interface) ([]byte, error) {
	if MTU == 0 {
		MTU = DEFAULT_MTU
	}

	raw_data := make([]byte, MTU)

	_, err := iface.Read(raw_data)
	return raw_data, err
}
