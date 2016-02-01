package protocol

import (
	"bytes"
	"fmt"
	"github.com/meshbird/meshbird/secure"
	"io"
)

var (
	magicKey = []byte{'M', 'E', 'S', 'H', 'B', 'I', 'R', 'D'}
)

func IsMagicValid(data []byte) bool {
	logger.Debug("decoding magic... expected %v, got %v", magicKey, data)
	return bytes.HasPrefix(data, magicKey)
}

func ExtractSessionKey(handshakeData []byte) []byte {
	return handshakeData[len(magicKey):]
}

func NewHandshakeRecord(sessionKey []byte, networkSecret *secure.NetworkSecret) *Record {
	// TODO: Encode session key
	sessionKey = append(magicKey, sessionKey...)
	return &Record{
		Type:    TypeHandshake,
		Version: CurrentVersion,
		Msg:     sessionKey,
	}
}

func ReadDecodeHandshake(r io.Reader) ([]byte, error) {
	return readDecodeMsg(r, TypeHandshake)
}

func WriteEncodeHandshake(w io.Writer, sessionKey []byte, networkSecret *secure.NetworkSecret) (err error) {
	logger.Debug("writing handshare message...")
	if err = EncodeAndWrite(w, NewHandshakeRecord(sessionKey, networkSecret)); err != nil {
		err = fmt.Errorf("error on write handshare message, %v", err)
	}
	return
}
