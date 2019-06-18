package transport

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/meshbird/meshbird/utils"
)

func makeAES128GCM(key string) (cipher.AEAD, error) {
	bkey := utils.MD5([]byte(key))
	block, err := aes.NewCipher(bkey)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func makeAES256GCM(key string) (cipher.AEAD, error) {
	bkey := utils.SHA256([]byte(key))
	block, err := aes.NewCipher(bkey)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
