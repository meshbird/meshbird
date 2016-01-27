package secure

import (
	"testing"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"time"
)

var (
	original = []byte("Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! ")
)

func TestEncryptIV(t *testing.T) {
	key := randomBytes(16)
	iv := randomBytes(16)

	ac, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	c := cipher.NewCBCEncrypter(ac, iv)

	encrypted := make(chan []byte, )
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("encrypted: %x", encrypted)

	decrypted, err := DecryptIV(encrypted, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, decrypted) {
		t.Fatal("original payload not equals to encrypted/decrypted")
	}
}

func TestEncryptAESGCM(t *testing.T) {
	key := randomBytes(aes.BlockSize)
	ac, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
}

func BenchmarkEncryptAesCbc(b *testing.B) {
	key := randomBytes(16)
	iv := make([]byte, 16)
	counter := 0
	dataLen := len(original)
	t := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := rand.Read(iv)
		if err != nil {
			b.Fatal(err)
		}
		_, err = EncryptIV(original, key, iv)
		if err != nil {
			b.Fatal(err)
		}
		counter += dataLen
	}
	b.StopTimer()
	ts := time.Since(t)
	b.Logf("encryption speed: %.2f Mbit/s", float64(counter) * 8 / ts.Seconds() / 1024 / 1024)
}

func BenchmarkDescryptAesCbc(b *testing.B) {
	key := randomBytes(16)
	iv := randomBytes(16)
	encrypted, err := EncryptIV(original, key, iv)
	if err != nil {
		b.Fatal(err)
	}
	counter := 0
	dataLen := len(encrypted)
	t := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DecryptIV(encrypted, key, iv)
		if err != nil {
			b.Fatal(err)
		}
		counter += dataLen
	}
	b.StopTimer()
	ts := time.Since(t)
	b.Logf("decryption speed: %.2f Mbit/s", float64(counter) * 8 / ts.Seconds() / 1024 / 1024)
}
