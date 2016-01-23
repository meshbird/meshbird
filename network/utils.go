package network

import (
	"crypto/rand"
	"net"
)

func GenerateIPAddress(network string) (string, error) {
	_, mask, err := net.ParseCIDR(network)

	randBytes := make([]byte, 4)
	_, err = rand.Read(randBytes)

	for i := 0; i < 4; i++ {

		randBytes[i] = randBytes[i]&^mask.Mask[i] | mask.IP[i]
	}

	return net.IP(randBytes).String(), err
}
