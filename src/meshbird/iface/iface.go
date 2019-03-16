package iface

import "github.com/songgao/water"

type Iface struct {
	name string
	ip   string
	mtu  int
	ifce *water.Interface
}

func New(name, ip string, mtu int) *Iface {
	return &Iface{
		name: name,
		ip:   ip,
		mtu:  mtu,
	}
}

func (i *Iface) Name() string {
	return i.ifce.Name()
}

func (i *Iface) Read(pkt PacketIP) (int, error) {
	return i.ifce.Read(pkt)
}

func (i *Iface) Write(pkt PacketIP) (int, error) {
	return i.ifce.Write(pkt)
}
