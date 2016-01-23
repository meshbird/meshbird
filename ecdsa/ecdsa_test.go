package ecdsa

import (
	"testing"

	"github.com/meshbird/ecdsa"
)

func generationTest(t *testing.T) {
	public := ecdsa.GeneratePublicKey()
	_, err := ecdsa.GeneratePrivateKey(public)
	if err != nil {
		t.Error("ECDSA key generation error: ", err)
	}
}
