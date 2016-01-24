package protocol_test

import (
	"bytes"
	"encoding/hex"
	"github.com/gophergala2016/meshbird/network/protocol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeHandShake(t *testing.T) {
	key := []byte("m7JC5ApnVKrTLs5L2Sz8SPMRckP2Wqkj")
	data := []byte{
		0, 33, // length
		1, // version
		0, // type
		109, 55, 74, 67, 53, 65, 112, 110, 86, 75, 114, 84, 76, 115, 53, 76,
		50, 83, 122, 56, 83, 80, 77, 82, 99, 107, 80, 50, 87, 113, 107, 106,
	}

	pack, err := protocol.Decode(bytes.NewBuffer(data))
	if assert.Nil(t, err) && assert.NotNil(t, pack) {
		assert.Equal(t, uint16(33), pack.Head.Length)
		assert.Equal(t, uint8(1), pack.Head.Version)
		assert.Equal(t, protocol.TypeHandshake, pack.Data.Type)
		assert.Empty(t, pack.Data.Vector)
		msg, ok := pack.Data.Msg.(protocol.HandshakeMessage)
		if assert.True(t, ok) {
			assert.Equal(t, key, []byte(msg))
		}
	}
}

func TestEncodeHandShake(t *testing.T) {
	key := []byte("hM57uSygtTFwp4f7fpdy6fdEJEqqXrZh")
	expected := []byte{
		0, 33, // length
		1, // version
		0, // type
		104, 77, 53, 55, 117, 83, 121, 103, 116, 84, 70, 119, 112, 52, 102, 55,
		102, 112, 100, 121, 54, 102, 100, 69, 74, 69, 113, 113, 88, 114, 90, 104,
	}

	pack := protocol.NewHandshakePacket(key)
	data, err := protocol.Encode(pack)
	if assert.Nil(t, err) {
		t.Logf("Data: %s", hex.Dump(data))
		assert.Equal(t, expected, data)
	}
}
