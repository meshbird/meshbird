package secure

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func EncryptIV(decrypted []byte, key []byte, iv []byte) ([]byte, error) {
	ac, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	c := cipher.NewCBCEncrypter(ac, iv)
	decrypted = PKCS5Padding(decrypted, ac.BlockSize())
	encrypted := make([]byte, len(decrypted))
	c.CryptBlocks(encrypted, decrypted)
	return encrypted, nil
}

func DecryptIV(encrypted []byte, key []byte, iv []byte) ([]byte, error) {
	ac, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	c := cipher.NewCBCDecrypter(ac, iv)
	decrypted := make([]byte, len(encrypted))
	c.CryptBlocks(decrypted, encrypted)
	decrypted = PKCS5UnPadding(decrypted)
	return decrypted, nil
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
