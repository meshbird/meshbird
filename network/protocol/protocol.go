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
	}

	typeNames = map[uint8]string{
		TypeHandshake: "Handshake",
		TypeOk:        "Ok",
		TypeHeartbeat: "Heartbeat",
		TypeGone:      "Gone",
		TypeTransfer:  "Transfer",
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

func Decode(data []byte, sessionKey []byte) (*Packet, error) {
	// TODO: sessionKey
	if len(data) < 4 { // Len(2) + Ver(1) + Type(1)
		return nil, ErrorToShort
	}

	var pack Packet
	reader := bytes.NewBuffer(data)

	if binary.Read(reader, binary.BigEndian, &pack.Head.Length) != nil {
		return nil, ErrorUnableToReadLength
	}
	if binary.Read(reader, binary.BigEndian, &pack.Head.Version) != nil {
		return nil, ErrorUnableToReadVersion
	}
	if binary.Read(reader, binary.BigEndian, &pack.Data.Type) != nil {
		return nil, ErrorUnableToReadType
	}
	if !isKnownType(pack.Data.Type) {
		return nil, ErrorUnknownType
	}

	remainLength := int(pack.Head.Length) - 1 // minus type

	if TypeHandshake != pack.Data.Type && TypeOk != pack.Data.Type {
		pack.Data.Vector = reader.Next(bodyVectorLen)
		if len(pack.Data.Vector) != bodyVectorLen {
			return nil, ErrorUnableToReadVector
		}
		remainLength -= bodyVectorLen
	}

	message := reader.Next(remainLength)
	if len(message) != remainLength {
		return nil, ErrorUnableToReadMessage
	}

	switch pack.Data.Type {
	case TypeHandshake:
		pack.Data.Msg = HandshakeMessage(message)
	case TypeOk:
		pack.Data.Msg = OkMessage(message)
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

func ReadAndDecode(r io.Reader, sessionKey []byte) (*Packet, error) {
	buf := make([]byte, 1500)
	n, errRead := r.Read(buf)

	if errRead != nil {
		if errRead != io.EOF {
			log.Printf("Error on read from connection: %s", errRead)
			return nil, errRead
		}
	}

	buf = buf[:n]
	log.Printf("Readed: %v", buf)

	pack, errDecode := Decode(buf, sessionKey)
	if errDecode != nil {
		log.Printf("Unable to decode packet: %s", errDecode)
		return nil, errDecode
	}

	return pack, nil
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
