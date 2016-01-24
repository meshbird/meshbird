package secure

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/lytics/base62"
	"net"
)

type NetworkSecret struct {
	Key []byte
	Net *net.IPNet
}

func NewNetworkSecret(network *net.IPNet) *NetworkSecret {
	return &NetworkSecret{
		Net: network,
		Key: randomBytes(16),
	}
}

func (ns NetworkSecret) Marshal() string {
	data := append(ns.Key, append(ns.Net.IP, ns.Net.Mask...)...)
	return base62.StdEncoding.EncodeToString(data)
}

func NetworkSecretUnmarshal(v string) (*NetworkSecret, error) {
	data, err := base62.StdEncoding.DecodeString(v)
	if err != nil {
		return nil, err
	}
	if len(data) != 24 {
		return nil, fmt.Errorf("mismatch secret length: 24 != %d", len(data))
	}
	secret := &NetworkSecret{
		Key: data[:16],
		Net: &net.IPNet{
			IP:   data[16:20],
			Mask: data[20:],
		},
	}
	return secret, nil
}

func (ns NetworkSecret) InfoHash() string {
	hashBytes := sha1.Sum([]byte(ns.Marshal()))
	return hex.EncodeToString(hashBytes[:])
}

func (ns NetworkSecret) CIDR() string {
	return ns.Net.String()
}

func (ns NetworkSecret) Encode(data []byte) []byte {
	return nil
}
