package network

import (
	"crypto/rand"
	log "github.com/Sirupsen/logrus"
	"net"
)

func GenerateIPAddress(cidr *net.IPNet) (net.IP, error) {
	log.Info("generate IP income: %v", cidr)
	randBytes := make([]byte, 4)
	_, err := rand.Read(randBytes)
	for i := 0; i < 4; i++ {
		randBytes[i] = randBytes[i]&^cidr.Mask[i] | cidr.IP[i]
	}
	log.Info("generate IP out: %v", net.IP(randBytes).String())
	return net.IP(randBytes), err
}
