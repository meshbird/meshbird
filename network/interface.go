package network

import (
	"fmt"
	"github.com/hsheth2/water/waterutil"
	"github.com/miolini/water"
	"net"
)

const DEFAULT_MTU = 1500

var MTU int

func Init() {
	MTU = 0
}

func CreateTunInterface(iface string) (*water.Interface, error) {
	ifce, err := water.NewTUN(iface)
	if err != nil {
		return nil, fmt.Errorf("create new tun interface %s err: %s", iface, err)
	}
	err = UpInterface(ifce.Name())
	if err != nil {
		return nil, fmt.Errorf("tun interface %s up err: %s", iface, err)
	}
	return ifce, nil
}

func CreateTunInterfaceWithIp(iface string, IpAddr string) (*water.Interface, error) {
	ifce, err := CreateTunInterface(iface)
	if err != nil {
		return nil, err
	}
	err = AssignIpAddress(ifce.Name(), IpAddr)
	return ifce, err
}

func NextNetworkPacket(iface *water.Interface) ([]byte, error) {
	if MTU == 0 {
		MTU = DEFAULT_MTU
	}

	raw_data := make([]byte, MTU)

	_, err := iface.Read(raw_data)
	return raw_data, err
}

func IPv4Destination(packet []byte) net.IP {
	return waterutil.IPv4Destination(packet)

}
