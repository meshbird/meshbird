package network

import (
	"net"
	"testing"
)

func TestDeviceCreate(t *testing.T) {
	iface := "tun0"
	_, err := CreateTunInterface(iface)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateDeviceWithIpAddr(t *testing.T) {
	iface := "tun0"
	IpAddr := "10.0.0.1/24"
	_, err := CreateTunInterfaceWithIp(iface, IpAddr)
	if err != nil {
		t.Error(err)
	}
	ifce, err := net.InterfaceByName(iface)

	if err != nil {
		t.Error(err)
	}
	addr, err := ifce.Addrs()
	if err != nil {
		t.Error(err)
	}
	if addr[0].String() != IpAddr {
		t.Error("Wrong Ip address on device")
	}
}
