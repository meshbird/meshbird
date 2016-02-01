package protocol

import (
	"fmt"
	"io"
	"net"
)

func ReadDecodePeerInfo(r io.Reader) ([]byte, error) {
	return readDecodeMsg(r, TypePeerInfo)
}

func WriteEncodePeerInfo(w io.Writer, privateIP net.IP) (err error) {
	logger.Debug("writing peer info message...")
	if err = EncodeAndWrite(w, NewRecord(TypePeerInfo, privateIP.To4())); err != nil {
		err = fmt.Errorf("error on write peer info message, %v", err)
	}
	return
}
