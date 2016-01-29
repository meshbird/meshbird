package secure

import (
	"testing"
	"time"
	"crypto/aes"
	"crypto/cipher"
)

var (
	original = []byte("Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird! Hello, World from MeshBird!")
)

func BenchmarkEncryptAesCbc(b *testing.B) {
	key := randomBytes(16)
	iv := make([]byte, 16)
	counter := 0
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		b.Fatal(err)
	}
	enc := cipher.NewCBCEncrypter(aesCipher, iv)
	decrypted := PKCS5Padding(original, aes.BlockSize)
	dataLen := len(decrypted)
	encrypted := make([]byte, len(decrypted))
	t := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.CryptBlocks(encrypted, decrypted)
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

	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		b.Fatal(err)
	}
	dec := cipher.NewCBCDecrypter(aesCipher, iv)
	t := time.Now()
	decrypted := make([]byte, len(encrypted))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dec.CryptBlocks(decrypted, encrypted)
		counter += dataLen
	}
	b.StopTimer()
	ts := time.Since(t)
	b.Logf("decryption speed: %.2f Mbit/s", float64(counter) * 8 / ts.Seconds() / 1024 / 1024)
}
