package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

func EncryptIV(decrypted []byte, key []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Println("[CRYPT][AES][ENC] Problem %s", err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Println("[CRYPT][AES][ENC] Problem %s", err.Error())
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("[CRYPT][AES][NONCE] Problem %s", err.Error())
		return nil, err
	}

	return gcm.Seal(nonce, nonce, decrypted, nil), nil

}

func DecryptIV(ciphertext []byte, key []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Println("[DECRYPT][AES] Problem %s", err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Println("[DECRYPT][AES] Problem %s", err.Error())
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		log.Println("[DECRYPT][AES] Problem %s", "Cyphertext too short")
		return nil, err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)

}
