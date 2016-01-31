package secure

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"crypto/aes"
	"crypto/cipher"
	"github.com/meshbird/meshbird/network"
)

type NetworkSecret struct {
	Key []byte
	Net *net.IPNet
}

func NewNetworkSecret(network *net.IPNet) *NetworkSecret {
	return &NetworkSecret{
		Net: network,
		Key: randomBytes(32),
	}
}

func (ns NetworkSecret) Bytes() []byte {
	return append(ns.Key, append(ns.Net.IP, ns.Net.Mask...)...)
}

func (ns NetworkSecret) Marshal() string {
	return hex.EncodeToString(ns.Bytes())
}

func NetworkSecretUnmarshal(v string) (*NetworkSecret, error) {
	data, err := hex.DecodeString(v)
	if err != nil {
		return nil, err
	}
	if len(data) != 40 {
		return nil, fmt.Errorf("mismatch secret length: 40 != %d", len(data))
	}
	secret := &NetworkSecret{
		Key: data[:32],
		Net: &net.IPNet{
			IP:   data[32:36],
			Mask: data[36:],
		},
	}
	return secret, nil
}

func (ns NetworkSecret) InfoHash() string {
	hashBytes := sha1.Sum(ns.Bytes())
	return hex.EncodeToString(hashBytes[:])
}

func (ns NetworkSecret) CIDR() string {
	return ns.Net.String()
}

func (ns NetworkSecret) Encode(dst []byte, data []byte, nonce []byte) error {
	aesCipher, err := aes.NewCipher(ns.Key)
	if err != nil {
		return err
	}
	c, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return err
	}
	c.Seal(dst, nonce, dst, nil)
	return nil
}

func (ns NetworkSecret) Decode(dst []byte, data []byte, nonce []byte) error {
	aesCipher, err := aes.NewCipher(ns.Key)
	if err != nil {
		return err
	}
	c, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return err
	}
	c.Open(dst, nonce, data, nil)
	return nil
}
