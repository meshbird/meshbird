package network

//import (
//	"net"
//	"os/exec"
//	"strings"
//	"testing"
//)
//
//func startReadData(ch chan<- []byte, iface *Interface) {
//	go func() {
//		for {
//			data, err := NextNetworkPacket(iface)
//			if err == nil {
//				ch <- data
//			}
//		}
//	}()
//}
//func startBroadCast(t *testing.T, dst string) {
//	if err := exec.Command("ping", "-b", "-c", "3", dst).Start(); err != nil {
//		t.Fatal(err)
//	}
//}
//func TestDeviceCreate(t *testing.T) {
//	iface := "tun0"
//	_, err := CreateTunInterface(iface)
//	if err != nil {
//		t.Fatal(err)
//	}
//	niface, err := net.InterfaceByName(iface)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if !strings.Contains(niface.Flags.String(), "up") {
//		t.Fatal("Interface not up")
//	}
//}
//
//func TestCreateDeviceWithIpAddr(t *testing.T) {
//	iface := "tun0"
//	IpAddr := "10.0.0.1/24"
//	_, err := CreateTunInterfaceWithIp(iface, IpAddr)
//	if err != nil {
//		t.Fatal(err)
//	}
//	ifce, err := net.InterfaceByName(iface)
//
//	if err != nil {
//		t.Fatal(err)
//	}
//	addr, err := ifce.Addrs()
//	if err != nil {
//		t.Fatal(err)
//	}
//	if addr[0].String() != IpAddr {
//		t.Fatal("Wrong Ip address on device")
//	}
//}
//
//func TestNextNetworkPacket(t *testing.T) {
//	ifaceName := "tun1"
//	IpAddrString := "10.1.0.1/24"
//	iface, err := CreateTunInterfaceWithIp(ifaceName, IpAddrString)
//	if err != nil {
//		t.Fatal(err)
//	}
//	startBroadCast(t, "10.1.0.255")
//	raw_data := make(chan []byte, 1500)
//	startReadData(raw_data, iface)
//}
//
//func TestSetMTU(t *testing.T) {
//	mtu := 1400
//	SetMTU(mtu)
//	if MTU != mtu {
//		t.Fatal("MTU not set")
//	}
//}
