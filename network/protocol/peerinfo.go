package protocol

import (
	"fmt"
	"io"
	"log"
	"net"
)

type (
	PeerInfoMessage []byte
)

func NewPeerInfoMessage(privateIP net.IP) *Packet {
	body := Body{
		Type: TypePeerInfo,
		Msg:  PeerInfoMessage(privateIP),
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

func ReadDecodePeerInfo(r io.Reader, sessionKey []byte) (PeerInfoMessage, error) {
	log.Printf("Trying to read PeerInfo message...")

	peerInfoPack, errDecode := ReadAndDecode(r, 8, sessionKey)
	if errDecode != nil {
		log.Printf("Unable to decode package: %s", errDecode)
		return nil, fmt.Errorf("Error on read PeerInfo package: %v", errDecode)
	}

	if peerInfoPack.Data.Type != TypePeerInfo {
		return nil, fmt.Errorf("Got non PeerInfo message: %+v", peerInfoPack)
	}

	log.Printf("Readed PeerInfo: %+v", peerInfoPack.Data.Msg)

	return peerInfoPack.Data.Msg.(PeerInfoMessage), nil
}

func WriteEncodePeerInfo(w io.Writer, privateIP net.IP) (err error) {
	log.Printf("Trying to write PeerInfo message...")
	if err = EncodeAndWrite(w, NewPeerInfoMessage(privateIP)); err != nil {
		err = fmt.Errorf("Error on write PeerInfo message: %v", err)
	}
	return
}
