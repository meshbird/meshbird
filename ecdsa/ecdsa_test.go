package ecdsa

import (
	"log"
	"testing"

	"github.com/gophergala2016/meshbird/ecdsa"
)

func generationTest(t *testing.T) {
	key, _ := ecdsa.GenerateKey()

	sg, hs, err := ecdsa.Sign(key.Private, []byte("Hello World!"))
	if err != nil {
		log.Println("ECDSA sign string error: ", err)
	}

	if !ecdsa.Verify(key.Public, hs, &sg) {
		log.Println("ECDSA verification test failed!")
	}

}
