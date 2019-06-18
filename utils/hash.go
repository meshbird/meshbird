package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func SHA256(in []byte) []byte {
	digest := sha256.New()
	digest.Write(in)
	return digest.Sum(nil)
}

func SHA1(in []byte) []byte {
	digest := sha1.New()
	digest.Write(in)
	return digest.Sum(nil)
}

func MD5(in []byte) []byte {
	digest := md5.New()
	digest.Write(in)
	return digest.Sum(nil)
}

func Hex(in []byte) string {
	return hex.EncodeToString(in)
}

func B64(in []byte) string {
	return base64.RawStdEncoding.EncodeToString(in)
}

func POE(err interface{}) {
	if err != nil {
		panic(err)
	}
}
