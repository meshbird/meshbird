package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/lytics/base62"
)

//Signature stores ecdsa r,s sign points
type Signature struct {
	r *big.Int
	s *big.Int
}

//Signature.Cat
func (sg *Signature) Cat() []byte {
	signature := sg.r.Bytes()
	signature = append(signature, sg.s.Bytes()...)
	return signature
}

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

// Sign signs msg with ECDSA private key
func Sign(priv *ecdsa.PrivateKey, msg []byte) (sg Signature, signhash []byte, err error) {
	sg.r = big.NewInt(0)
	sg.s = big.NewInt(0)
	//var h hash.Hash
	//h = sha256.New()
	//signhash = h.Sum(msg)
	signhash = make([]byte, base62.StdEncoding.EncodedLen(len(msg)))
	base62.StdEncoding.Encode(signhash, msg)

	sg.r, sg.s, err = ecdsa.Sign(rand.Reader, priv, signhash)
	fmt.Printf("Signature: %x \n", sg.Cat())
	return
}

// Verify use publick key to verify signature
func Verify(pubKey *ecdsa.PublicKey, signhash []byte, sg Signature) bool {
	return ecdsa.Verify(pubKey, signhash, sg.r, sg.s)
}
