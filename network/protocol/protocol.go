package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
)

const (
	TypeHandshake uint8 = iota
	TypeOk
	TypeHeartbeat
	TypeGone
	TypeTransfer
	TypePeerInfo
)

const (
	CurrentVersion = 1
	bodyVectorLen  = 16
)

var (
	logger = log.New(os.Stderr, "[proto] ", log.LstdFlags)

	ErrorUnableToReadVector  = errors.New("unable to read vector")
	ErrorUnableToReadMessage = errors.New("unable to read message")
	ErrorUnknownType         = errors.New("unknown type")

	knownTypes = []uint8{
		TypeHandshake,
		TypeOk,
		TypeHeartbeat,
		TypeGone,
		TypeTransfer,
		TypePeerInfo,
	}

	typeNames = map[uint8]string{
		TypeHandshake: "Handshake",
		TypeOk:        "Ok",
		TypeHeartbeat: "Heartbeat",
		TypeGone:      "Gone",
		TypeTransfer:  "Transfer",
		TypePeerInfo:  "PeerInfo",
	}
)

type (
	Message interface {
		io.WriterTo

		Len() uint16
	}

	Header struct {
		Length  uint16
		Version uint8
	}
	Body struct {
		Type   uint8
		Vector []byte
		Msg    Message
	}
	Packet struct {
		Head Header
		Data Body
	}
)

func (h Header) Len() uint16 {
	return 3
}

func (h *Header) WriteTo(w io.Writer) (n int64, err error) {
	binary.Write(w, binary.BigEndian, h.Length)
	binary.Write(w, binary.BigEndian, h.Version)
	return
}

func (b Body) Len() uint16 {
	return b.Msg.Len() + uint16(len(b.Vector)+1)
}

func (b *Body) WriteTo(w io.Writer) (n int64, err error) {
	binary.Write(w, binary.BigEndian, b.Type)
	if len(b.Vector) > 0 {
		binary.Write(w, binary.BigEndian, b.Vector)
	}
	b.Msg.WriteTo(w)
	return
}

func (p Packet) Len() uint16 {
	return p.Head.Len() + p.Data.Len()
}

func Decode(r io.Reader) (*Packet, error) {
	var pack Packet

	if err := binary.Read(r, binary.BigEndian, &pack.Head.Length); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &pack.Head.Version); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.BigEndian, &pack.Data.Type); err != nil {
		return nil, err
	}
	if !isKnownType(pack.Data.Type) {
		return nil, ErrorUnknownType
	}

	remainLength := int(pack.Head.Length) - 1 // minus type

	if TypeHandshake != pack.Data.Type && TypeOk != pack.Data.Type && TypePeerInfo != pack.Data.Type && TypeHeartbeat != pack.Data.Type {
		vector := make([]byte, bodyVectorLen)
		if n, err := r.Read(vector); err != nil || n != bodyVectorLen {
			if n != bodyVectorLen {
				err = ErrorUnableToReadVector
			}
			return nil, err
		}
		pack.Data.Vector = vector
		remainLength -= bodyVectorLen
	}

	message := make([]byte, remainLength)
	if n, err := r.Read(message); err != nil || n != remainLength {
		if n != remainLength {
			err = ErrorUnableToReadMessage
		}
		return nil, err
	}

	switch pack.Data.Type {
	case TypeHandshake:
		pack.Data.Msg = HandshakeMessage(message)
	case TypeOk:
		pack.Data.Msg = OkMessage(message)
	case TypePeerInfo:
		pack.Data.Msg = PeerInfoMessage(message)
	case TypeTransfer:
		pack.Data.Msg = TransferMessage(message)
	case TypeHeartbeat:
		pack.Data.Msg = HeartbeatMessage(message)
	case TypeGone:
		pack.Data.Msg = GoneMessage(message)
	}

	return &pack, nil
}

func Encode(pack *Packet) ([]byte, error) {
	writer := new(bytes.Buffer)
	writer.Grow(int(pack.Len()))

	pack.Head.WriteTo(writer)
	pack.Data.WriteTo(writer)

	return writer.Bytes(), nil
}

func ReadAndDecode(r io.Reader) (*Packet, error) {
	pack, errDecode := Decode(r)
	if errDecode != nil {
		logger.Printf("Unable to decode packet: %s", errDecode)
		return nil, errDecode
	}

	logger.Printf("Received packet: %+v", pack)

	return pack, nil
}

func EncodeAndWrite(w io.Writer, pack *Packet) error {
	logger.Printf("Encoding package: %+v", pack)

	typeName := TypeName(pack.Data.Type)

	reply, errEncode := Encode(pack)
	if errEncode != nil {
		logger.Printf("Error on encoding %s: %v", typeName, errEncode)
		return errEncode
	}

	logger.Printf("Sending %s message %d bytes...", typeName, len(reply))

	n, err := w.Write(reply)
	if err != nil {
		logger.Printf("Error on write %s: %v", typeName, err)
		return err
	}

	logger.Printf("%d of %d bytes of %s message sent", n, len(reply), typeName)

	return nil
}

func TypeName(t uint8) string {
	return typeNames[t]
}

func isKnownType(needle uint8) bool {
	for _, t := range knownTypes {
		if needle == t {
			return true
		}
	}
	return false
}
