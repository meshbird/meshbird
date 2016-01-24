package network

import (
	"crypto/rand"
	"log"
	"net"
)

func GenerateIPAddress(cidr *net.IPNet) (string, error) {
	log.Printf("generate IP income: %v", cidr)
	randBytes := make([]byte, 4)
	_, err := rand.Read(randBytes)
	for i := 0; i < 4; i++ {
		randBytes[i] = randBytes[i]&^cidr.Mask[i] | cidr.IP[i]
	}
	log.Printf("generate IP out: %v", net.IP(randBytes).String())
	return net.IP(randBytes).String(), err
}
