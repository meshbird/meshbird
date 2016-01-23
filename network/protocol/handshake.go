package protocol

import (
	"bytes"
	"errors"
	"io"
)

const (
	handshakeKeyLen = 16
)

var (
	ErrorUnknownHandshakeMagic    = errors.New("unknown handshace magic")
	ErrorUnableToReadHandshakeKey = errors.New("unable to read hankshake key")

	handshakeMagic = []byte{
		'M', 'E', 'S', 'H',
	}
)

type (
	HandshakeMessage struct {
		Message

		Magic []byte
		Key   []byte
	}
)

func NewHandshakeMessage(key []byte) HandshakeMessage {
	return HandshakeMessage{
		Magic: handshakeMagic,
		Key:   key,
	}
}

func (m HandshakeMessage) Len() uint16 {
	return uint16(len(handshakeMagic) + handshakeKeyLen)
}

func (m HandshakeMessage) WriteTo(w io.Writer) (n int64, err error) {
	if _, err = w.Write(m.Magic); err != nil {
		return
	}
	if _, err = w.Write(m.Key); err != nil {
		return
	}
	n = int64(len(m.Magic) + len(m.Key))
	return
}

func decodeHandshake(data []byte) (HandshakeMessage, error) {
	var msg HandshakeMessage
	reader := bytes.NewBuffer(data)

	msg.Magic = reader.Next(len(handshakeMagic))
	if !bytes.Equal(msg.Magic, handshakeMagic) {
		return msg, ErrorUnknownHandshakeMagic
	}

	msg.Key = reader.Next(handshakeKeyLen)
	if len(msg.Key) != handshakeKeyLen {
		return msg, ErrorUnableToReadHandshakeKey
	}

	return msg, nil
}
