package network

import (
	"crypto/rand"
	"net"
)

func GenerateIPAddress(cidr *net.IPNet) (net.IP, error) {
	randBytes := make([]byte, 4)
	_, err := rand.Read(randBytes)
	for i := 0; i < 4; i++ {
		randBytes[i] = randBytes[i]&^cidr.Mask[i] | cidr.IP[i]
	}
	return net.IP(randBytes), err
}
