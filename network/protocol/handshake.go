package protocol

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/meshbird/meshbird/secure"
	"io"
)

var (
	magicKey = []byte{'M', 'E', 'S', 'H', 'B', 'I', 'R', 'D'}
)

type (
	HandshakeMessage []byte
)

func IsMagicValid(data []byte) bool {
	logger.WithFields(log.Fields{
		"magic": magicKey,
		"data":  data,
	}).Debug("Trying to check magic")
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

func ReadDecodeHandshake(r io.Reader) (HandshakeMessage, error) {
	logger.Debug("Trying to read Handshake message...")

	okPack, errDecode := ReadAndDecode(r)
	if errDecode != nil {
		logger.WithError(errDecode).Error("Unable to decode package")
		return nil, fmt.Errorf("Error on read Handshake package: %v", errDecode)
	}

	if okPack.Data.Type != TypeHandshake {
		return nil, fmt.Errorf("Got non Handshake message: %+v", okPack)
	}

	logger.WithField("msg", okPack.Data.Msg).Debug("Readed Handshake")

	return okPack.Data.Msg.(HandshakeMessage), nil
}

func WriteEncodeHandshake(w io.Writer, sessionKey []byte, networkSecret *secure.NetworkSecret) (err error) {
	logger.Debug("Trying to write Handshake message...")
	if err = EncodeAndWrite(w, NewHandshakePacket(sessionKey, networkSecret)); err != nil {
		err = fmt.Errorf("Error on write Handshake message: %v", err)
	}
	return
}
