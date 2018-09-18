package iface

import (
	"net"
)

type PacketIP []byte

func NewPacketIP(size int) PacketIP {
	return PacketIP(make([]byte, size))
}

func (p PacketIP) GetSourceIP() net.IP {
	return net.IP(p[12:16])
}

func (p PacketIP) GetDestinationIP() net.IP {
	return net.IP(p[16:20])
}