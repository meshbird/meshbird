package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"math/big"
)

//GenerateKey generates a public & private key pair
func GenerateKey() (pubKey ecdsa.PublicKey, privKey *ecdsa.PrivateKey, err error) {
	privKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // this generates a public & private key pair
	pubKey = privKey.PublicKey
	fmt.Println("Private Key :")
	fmt.Printf("%x \n", privKey)

	fmt.Println("Public Key :")
	fmt.Printf("%x \n", pubKey)
	return
}

func getSignature(r, s *big.Int) []byte {
	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	return signature
}

// SignString signs msg with ECDSA private key
func SignString(priv *ecdsa.PrivateKey, msg string) (r, s *big.Int, signhash []byte, err error) {
	r = big.NewInt(0)
	s = big.NewInt(0)
	var h hash.Hash
	h = sha256.New()
	io.WriteString(h, msg)
	signhash = h.Sum(nil)
	r, s, err = ecdsa.Sign(rand.Reader, priv, signhash)
	fmt.Printf("Signature: %x \n", getSignature(r, s))
	return
}

// Verify use publick key to verify signature
func Verify(pubKey *ecdsa.PublicKey, signhash []byte, r, s *big.Int) bool {
	return ecdsa.Verify(pubKey, signhash, r, s)
}
