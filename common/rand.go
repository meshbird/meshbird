package common

import (
	"crypto/rand"
)

func RandomBytes(len int) []byte {
	buf := make([]byte, len)
	if n, err := rand.Read(buf); err != nil || n != len {
		panic("something wrong")
	}
	return buf
}