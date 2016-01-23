package protocol_test

import (
	"encoding/hex"
	"github.com/gophergala2016/meshbird/network/protocol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeHandShake(t *testing.T) {
	key := []byte("nNJkE2LUGZ7AVpjw")
	data := []byte{
		0, 21, // length
		1,              // version
		0,              // type
		77, 69, 83, 72, // MESH
		110, 78, 74, 107, 69, 50, 76, 85, 71, 90, 55, 65, 86, 112, 106, 119, // nNJkE2LUGZ7AVpjw
	}

	pack, err := protocol.Decode(data)
	if assert.Nil(t, err) && assert.NotNil(t, pack) {
		assert.Equal(t, uint16(21), pack.Head.Length)
		assert.Equal(t, uint8(1), pack.Head.Version)
		assert.Equal(t, protocol.TypeHandshake, pack.Data.Type)
		assert.Empty(t, pack.Data.Vector)
		msg, ok := pack.Data.Msg.(protocol.HandshakeMessage)
		if assert.True(t, ok) {
			assert.Equal(t, []byte("MESH"), msg.Magic)
			assert.Equal(t, key, msg.Key)
		}
	}
}

func TestEncodeHandShake(t *testing.T) {
	key := []byte("MCPqt8z2DcyhQzfj")
	expected := []byte{
		0, 21, // length
		1,              // version
		0,              // type
		77, 69, 83, 72, // MESH
		77, 67, 80, 113, 116, 56, 122, 50, 68, 99, 121, 104, 81, 122, 102, 106, // MCPqt8z2DcyhQzfj
	}

	pack := protocol.NewHandshakePacket(key)
	data, err := protocol.Encode(pack)
	if assert.Nil(t, err) {
		t.Logf("Data: %s", hex.Dump(data))
		assert.Equal(t, expected, data)
	}
}
