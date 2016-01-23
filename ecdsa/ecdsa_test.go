package ecdsa

import (
	"testing"

	"github.com/gophergala2016/meshbird/ecdsa"
)

func generationTest(t *testing.T) {
	public, private, err := ecdsa.GenerateKey()

	if err != nil {
		t.Error("ECDSA key generation error: ", err)
	}
	sg, hs, err := ecdsa.Sign(private, []byte("Hello World!"))
	if err != nil {
		t.Error("ECDSA sign string error: ", err)
	}
	return
	if !ecdsa.Verify(&public, hs, &sg) {
		t.Error("ECDSA verification test failed!")
	}
	sg1 := &ecdsa.Signature{}
	sg1.Decode(sg.Encode())
	if !ecdsa.Verify(&public, hs, sg1) {
		t.Error("ECDSA signature encoding/decoding test failed!")
	}
}
