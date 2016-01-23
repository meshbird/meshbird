package network

import (
	"net"
	"strings"
	"testing"
)

func TestDeviceCreate(t *testing.T) {
	iface := "tun0"
	_, err := CreateTunInterface(iface)
	if err != nil {
		t.Error(err)
	}
	niface, err := net.InterfaceByName(iface)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(niface.Flags.String(), "up") {
		t.Error("Interface not up")
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
