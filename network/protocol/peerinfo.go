package protocol

import (
	"fmt"
	"io"
	"net"
)

type (
	PeerInfoMessage []byte
)

func NewPeerInfoMessage(privateIP net.IP) *Packet {
	body := Body{
		Type: TypePeerInfo,
		Msg:  PeerInfoMessage(privateIP.To4()),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (m PeerInfoMessage) Len() uint16 {
	return uint16(len(m))
}

func (m PeerInfoMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(m)
	return int64(n), err
}

func (m PeerInfoMessage) PrivateIP() net.IP {
	return net.IP(m)
}

func ReadDecodePeerInfo(r io.Reader) (PeerInfoMessage, error) {
	logger.Debug("reading peer info message...")

	peerInfoPack, errDecode := ReadAndDecode(r)
	if errDecode != nil {
		logger.Error("error on package decode, %v", errDecode)
		return nil, fmt.Errorf("error on read peer info package, %v", errDecode)
	}

	if peerInfoPack.Data.Type != TypePeerInfo {
		return nil, fmt.Errorf("non peer info message received, %+v", peerInfoPack)
	}

	logger.Debug("message, %v", peerInfoPack.Data.Msg)
	return peerInfoPack.Data.Msg.(PeerInfoMessage), nil
}

func WriteEncodePeerInfo(w io.Writer, privateIP net.IP) (err error) {
	logger.Debug("writing peer info message...")
	if err = EncodeAndWrite(w, NewPeerInfoMessage(privateIP)); err != nil {
		err = fmt.Errorf("error on write peer info message, %v", err)
	}
	return
}
