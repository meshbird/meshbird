package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

//GeneratePublicKey generates a elliptic curve for further pub+priv key pair generation using GeneratePrivateKey method
func GeneratePublicKey() elliptic.Curve {
	return elliptic.P256()
}

//GeneratePrivateKey generates private key
func GeneratePrivateKey(pub elliptic.Curve) (privKey *ecdsa.PrivateKey, err error) {
	privKey, err = ecdsa.GenerateKey(pub, rand.Reader) // this generates a public & private key pair
	fmt.Println("Private Key :")
	fmt.Printf("%x \n", privKey)

	fmt.Println("Public Key :")
	fmt.Printf("%x \n", pub)
	return
}

/*// Sign signs msg with ECDSA private key
func Sign(priv *ecdsa.PrivateKey, msg []byte) (signature []byte, err error) {
	r := big.NewInt(0)
	s := big.NewInt(0)
	var h hash.Hash
	h = sha256.New()
	signhash := h.Sum(nil)
	r, s, err = ecdsa.Sign(rand.Reader, priv, signhash)
	signature = r.Bytes()
	signature = append(signature, s.Bytes()...)
	return
}

func Verify() {}
*/
