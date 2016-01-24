package protocol

import (
	"bytes"
	"github.com/gophergala2016/meshbird/secure"
	"io"
)

var (
	magicKey = []byte{'M', 'E', 'S', 'H', 'B', 'I', 'R', 'D'}
)

type (
	HandshakeMessage []byte
)

func IsMagicValid(data []byte) bool {
	return bytes.HasPrefix(data, magicKey)
}

func NewHandshakePacket(sessionKey []byte, networkSecret *secure.NetworkSecret) *Packet {
	sessionKey = append(magicKey, sessionKey...)
	data := networkSecret.Encode(sessionKey)

	body := Body{
		Type: TypeHandshake,
		Msg:  HandshakeMessage(data),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (m HandshakeMessage) Len() uint16 {
	return uint16(len(m))
}

func (m HandshakeMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(m)
	return int64(n), err
}

func (m HandshakeMessage) Bytes() []byte {
	return []byte(m)
}

func (m HandshakeMessage) SessionKey() []byte {
	return m[len(magicKey):]
}