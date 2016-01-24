package common

import (
	"net"
	"fmt"
)

type NetworkSecret struct {
	Key    []byte
	IP     net.IP
	IPMask net.IPMask
}

func (ns *NetworkSecret) Marshal() string {
	data := append(ns.Key, append(ns.IP, ns.IPMask...)...)
	return string(data)
}

func NetworkSecretUnmarshal(v string) (*NetworkSecret, error) {
	data := []byte(v)
	if len(data) != 24 {
		return nil, fmt.Errorf("mismatch secret length: 24 != %d", len(data))
	}
	secret := &NetworkSecret{
		Key: data[:16],
		IP: data[16:20],
		IPMask: data[20:],
	}
	return secret, nil
}