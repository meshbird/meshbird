package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/lytics/base62"
)

//Signature stores ecdsa r,s sign points
type Signature struct {
	r *big.Int
	s *big.Int
}

type Key struct {
	Private *ecdsa.PrivateKey
	Public  *ecdsa.PublicKey
}

//Signature.Encode return encoded r, s pair
func (sg Signature) Encode() []byte {
	//signature := sg.r.Bytes()
	//signature = append(signature, sg.s.Bytes()...)
	signature := pointsToDER(sg.r.Bytes(), sg.s.Bytes())
	return signature
}

func (sg Signature) String() string {
	return string(sg.Encode())
}

//Signature.Decode return decoded r, s pair
func (sg *Signature) Decode(data []byte) {
	r, s := pointsFromDER(data)
	sg.r.SetBytes(r)
	sg.s.SetBytes(s)
}

//GenerateKey generates a public & private key pair
func GenerateKey() (key *Key, err error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // this generates a public & private key pair
	pubKey := &privKey.PublicKey
	key = &Key{privKey, pubKey}
	return
}

func InitKey(C elliptic.Curve, D, X, Y *big.Int) ecdsa.PrivateKey {
	priv := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: C,
			X:     X,
			Y:     Y,
		},
		D: D,
	}
	return priv
}

func (k Key) GetPublic() string {
	buf := elliptic.Marshal(elliptic.P256(), k.Public.X, k.Public.Y)
	encoded := make([]byte, base62.StdEncoding.EncodedLen(len(buf)))
	base62.StdEncoding.Encode(encoded, buf)
	return string(encoded)
}

func (k *Key) ParsePublic(q string) {
	buf := []byte(q)
	decoded := make([]byte, base62.StdEncoding.DecodedLen(len(buf)))
	base62.StdEncoding.Decode(decoded, buf)
	c := elliptic.P256()
	x, y := elliptic.Unmarshal(c, decoded)
	k.Public = &ecdsa.PublicKey{
		Curve: c,
		X:     x,
		Y:     y,
	}
}

func (k Key) GetPrivate() string {
	// Private Key: D + Public
	return string(pointsToDER(k.Private.D.Bytes(), elliptic.Marshal(elliptic.P256(), k.Public.X, k.Public.Y)))
}

func (k *Key) ParsePrivate(q string) {
	D := big.NewInt(0)
	d, public := pointsFromDER([]byte(q))
	D.SetBytes(d)
	c := elliptic.P256()
	x, y := elliptic.Unmarshal(c, public)

	k.Private = &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: c,
			X:     x,
			Y:     y,
		},
		D: D,
	}
	k.Public = &k.Private.PublicKey
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
	fmt.Printf("Signature: %x \n", sg.Encode())
	return
}

// Verify use publick key to verify signature
func Verify(pubKey *ecdsa.PublicKey, signhash []byte, sg *Signature) bool {
	return ecdsa.Verify(pubKey, signhash, sg.r, sg.s)
}

// Convert an ECDSA signature (points R and S) to a byte array using ASN.1 DER encoding.
// This is a port of Bitcore's Key.rs2DER method.
func pointsToDER(r, s []byte) []byte {
	// Ensure MSB doesn't break big endian encoding in DER sigs
	prefixPoint := func(b []byte) []byte {
		if len(b) == 0 {
			b = []byte{0x00}
		}
		if b[0]&0x80 != 0 {
			paddedBytes := make([]byte, len(b)+1)
			copy(paddedBytes[1:], b)
			b = paddedBytes
		}
		return b
	}

	rb := prefixPoint(r)
	sb := prefixPoint(s)
	//rb := prefixPoint(r.Bytes())
	//sb := prefixPoint(s.Bytes())

	// DER encoding:
	// 0x30 + z + 0x02 + len(rb) + rb + 0x02 + len(sb) + sb
	length := 2 + len(rb) + 2 + len(sb)

	der := append([]byte{0x30, byte(length), 0x02, byte(len(rb))}, rb...)
	der = append(der, 0x02, byte(len(sb)))
	der = append(der, sb...)

	encoded := make([]byte, base62.StdEncoding.EncodedLen(len(der)))
	base62.StdEncoding.Encode(encoded, der)

	return encoded
}

// Get the X and Y points from a DER encoded signature
// Sometimes demarshalling using Golang's DEC to struct unmarshalling fails; this extracts R and S from the bytes
// manually to prevent crashing.
// This should NOT be a hex encoded byte array
func pointsFromDER(der []byte) ([]byte, []byte) {
	decoded := make([]byte, base62.StdEncoding.DecodedLen(len(der)))
	base62.StdEncoding.Decode(decoded, der)
	//R, S := &big.Int{}, &big.Int{}

	data := asn1.RawValue{}
	if _, err := asn1.Unmarshal(decoded, &data); err != nil {
		panic(err.Error())
	}

	// The format of our DER string is 0x02 + rlen + r + 0x02 + slen + s
	rLen := data.Bytes[1] // The entire length of R + offset of 2 for 0x02 and rlen
	r := data.Bytes[2 : rLen+2]
	// Ignore the next 0x02 and slen bytes and just take the start of S to the end of the byte array
	s := data.Bytes[rLen+4:]

	//R.SetBytes(r)
	//S.SetBytes(s)

	return r, s
}
