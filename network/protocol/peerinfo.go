package protocol

import (
	"io"
	"net"
)

type (
	PeerInfoMessage []byte
)

func NewPeerInfoMessage(privateIP net.IP) *Packet {
	body := Body{
		Type: TypeOk,
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
