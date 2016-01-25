package secure

import "testing"

func TestEncryptIV(t *testing.T) {
	key := randomBytes(16)
	iv := randomBytes(16)
	decrypted := []byte("Hello, World from MeshBird!")
	encrypted, err := EncryptIV(decrypted, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("encrypted: %x", encrypted)

	decrypted2, err := DecryptIV(encrypted, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("decrypted: %s", string(decrypted2))
}
