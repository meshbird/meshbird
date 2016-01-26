package secure

import (
	"testing"
	"bytes"
	"crypto/rand"
	"time"
)

var (
	original = []byte("Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! ")
)

func TestEncryptIV(t *testing.T) {
	key := randomBytes(16)
	iv := randomBytes(16)
	encrypted, err := EncryptIV(original, key, iv)
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
	b.Logf("encrypt speed: %.2f GB/s", float64(counter) / ts.Seconds() / 1024 / 1024 / 1024)
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
	b.Logf("encrypt speed: %.2f GB/s", float64(counter) / ts.Seconds() / 1024 / 1024 / 1024)
}
