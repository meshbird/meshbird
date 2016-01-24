package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
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
	ErrorToShort             = errors.New("data length is too short")
	ErrorUnableToReadLength  = errors.New("unable to read length")
	ErrorUnableToReadVersion = errors.New("unable to read version")
	ErrorUnableToReadType    = errors.New("unable to read type")
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

	log.Println("Decoding...")

	if binary.Read(r, binary.BigEndian, &pack.Head.Length) != nil {
		return nil, ErrorUnableToReadLength
	}

	log.Printf("Packet body length: %d", pack.Head.Length)

	if binary.Read(r, binary.BigEndian, &pack.Head.Version) != nil {
		return nil, ErrorUnableToReadVersion
	}

	log.Printf("Packet version: %d", pack.Head.Version)

	if binary.Read(r, binary.BigEndian, &pack.Data.Type) != nil {
		return nil, ErrorUnableToReadType
	}

	log.Printf("Packet type: %d (%s)", pack.Data.Type, TypeName(pack.Data.Type))

	if !isKnownType(pack.Data.Type) {
		return nil, ErrorUnknownType
	}

	remainLength := int(pack.Head.Length) - 1 // minus type

	if TypeHandshake != pack.Data.Type && TypeOk != pack.Data.Type && TypePeerInfo != pack.Data.Type {
		log.Println("Reading vector ...")

		vector := make([]byte, bodyVectorLen)
		if n, err := r.Read(vector); err != nil || n != bodyVectorLen {
			return nil, ErrorUnableToReadVector
		}
		pack.Data.Vector = vector
		remainLength -= bodyVectorLen

		log.Printf("Vector: %v", pack.Data.Vector)
	}

	log.Printf("Data length: %d", remainLength)

	message := make([]byte, remainLength)
	if n, err := r.Read(message); err != nil || n != remainLength {
		return nil, ErrorUnableToReadMessage
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
		log.Printf("Unable to decode packet: %s", errDecode)
		return nil, errDecode
	}

	log.Printf("Received packet: %+v", pack)

	return pack, nil
}

func EncodeAndWrite(w io.Writer, pack *Packet) error {
	log.Printf("Encoding package: %+v", pack)

	typeName := TypeName(pack.Data.Type)

	reply, errEncode := Encode(pack)
	if errEncode != nil {
		log.Printf("Error on encoding %s: %v", typeName, errEncode)
		return errEncode
	}

	log.Printf("Sending %s message %d bytes...", typeName, len(reply))

	n, err := w.Write(reply)
	if err != nil {
		log.Printf("Error on write %s: %v", typeName, err)
		return err
	}

	log.Printf("%d of %d bytes of %s message sent", n, len(reply), typeName)

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
