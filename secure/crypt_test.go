package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"
	"time"
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
	b.SetBytes(int64(len(decrypted)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.CryptBlocks(encrypted, decrypted)
		counter += dataLen
	}
	b.StopTimer()
	ts := time.Since(t)
	b.Logf("encryption speed: %.2f Mbit/s", float64(counter)*8/ts.Seconds()/1024/1024)
}

func BenchmarkDescryptAesCbc(b *testing.B) {
	key := randomBytes(16)
	iv := randomBytes(16)
	encrypted, err := EncryptIV(original, key)
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
	b.SetBytes(int64(len(encrypted)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dec.CryptBlocks(decrypted, encrypted)
		counter += dataLen
	}
	b.StopTimer()
	ts := time.Since(t)
	b.Logf("decryption speed: %.2f Mbit/s", float64(counter)*8/ts.Seconds()/1024/1024)
}

func BenchmarkEncryptAESGCM(b *testing.B) {
	key := make([]byte, 32)
	a, _ := aes.NewCipher(key)
	c, _ := cipher.NewGCM(a)
	benchmarkAEAD(b, c)
}

const benchSize = 1500

func benchmarkAEAD(b *testing.B, c cipher.AEAD) {
	b.SetBytes(benchSize)
	input := make([]byte, benchSize)
	output := make([]byte, benchSize)
	nonce := make([]byte, c.NonceSize())
	counter := 0
	t := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Seal(output[:0], nonce, input, nil)
		counter += len(output)
	}
	ts := time.Since(t)
	b.Logf("aes-gcm speed: %.2f Mbit/s", float64(counter)*8/ts.Seconds()/1024/1024)
}
