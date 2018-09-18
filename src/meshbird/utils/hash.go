package utils

import (
	"encoding/hex"
	"encoding/base64"
	"crypto/sha256"
	"crypto/aes"
	"crypto/rand"
	"crypto/cipher"
	"io"
)

func SHA256(in []byte) []byte {
	digest := sha256.New()
	digest.Write(in)
	return digest.Sum(nil)
}

func Hex(in []byte) string {
	return hex.EncodeToString(in)
}

func B64(in []byte) string {
	return base64.RawStdEncoding.EncodeToString(in)
}

func AES256GCM(key, in []byte) []byte {
	block, err := aes.NewCipher(key)
	POE(err)
	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	POE(err)
	aesgcm, err := cipher.NewGCM(block)
	POE(err)
	ciphertext := aesgcm.Seal(nil, nonce, in, nil)
	return ciphertext
}

func POE(err interface{}) {
	if err != nil {
		panic(err)
	}
}