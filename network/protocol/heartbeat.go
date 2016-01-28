package protocol

import (
	"io"
	"net"
)

type (
	HeartbeatMessage []byte
)

func NewHeartbeatMessage(privateIP net.IP) *Packet {
	body := Body{
		Type: TypeHeartbeat,
		Msg:  HeartbeatMessage(privateIP.To4()),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (m HeartbeatMessage) Len() uint16 {
	return uint16(len(m))
}

func (m HeartbeatMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(m)
	return int64(n), err
}
