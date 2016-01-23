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
	r, s, hs, err := ecdsa.SignString(private, "Hello World!")
	if err != nil {
		t.Error("ECDSA sign string error: ", err)
	}
	if !ecdsa.Verify(&public, hs, r, s) {
		t.Error("ECDSA verification test failed!")
	}
}
