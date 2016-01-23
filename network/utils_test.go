package network

import (
	"net"
	"testing"
)

func TestGenerateIPAddress(t *testing.T) {
	_, ipnet, err := net.ParseCIDR("192.168.0.0/24")
	IPAddress, err := GenerateIPAddress("192.168.0.0/24")
	if err != nil {
		t.Fatal(err)
	}
	genIP := net.ParseIP(IPAddress)
	if !ipnet.Contains(genIP) {
		t.Fatal("Generated wrong IP address")
	}

}
