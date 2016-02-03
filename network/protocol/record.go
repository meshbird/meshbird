package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type (
	Record struct {
		Type    uint8
		Version uint8
		Length  uint16
		Vector  []byte
		Msg     []byte
	}
)

func NewRecord(t uint8, m []byte) *Record {
	return &Record{
		Type:    t,
		Version: CurrentVersion,
		Msg:     m,
	}
}

func Decode(r io.Reader) (*Record, error) {
	var rec Record

	// Type
	if err := binary.Read(r, binary.BigEndian, &rec.Type); err != nil {
		return nil, err
	}
	if !isKnownType(rec.Type) {
		return nil, ErrorUnknownType
	}

	// Version
	if err := binary.Read(r, binary.BigEndian, &rec.Version); err != nil {
		return nil, err
	}

	// Length
	if err := binary.Read(r, binary.BigEndian, &rec.Length); err != nil {
		return nil, err
	}

	remainLength := int64(rec.Length)

	// Only `Transfer` has vector
	if TypeTransfer == rec.Type {
		vector := make([]byte, bodyVectorLen)
		if n, err := r.Read(vector); err != nil || n != bodyVectorLen {
			if n != bodyVectorLen {
				err = ErrorUnableToReadVector
			}
			return nil, err
		}
		rec.Vector = vector
		remainLength -= bodyVectorLen
	}

	buf := bytes.Buffer{}
	if n, err := io.CopyN(&buf, r, remainLength); err != nil || n != remainLength {
		if n != remainLength {
			err = ErrorUnableToReadMessage
		}
		return nil, err
	}

	rec.Msg = buf.Bytes()
	return &rec, nil
}

func Encode(rec *Record) ([]byte, error) {
	var err error
	writer := new(bytes.Buffer)

	// Type
	if err = binary.Write(writer, binary.BigEndian, rec.Type); err != nil {
		return nil, err
	}

	// Version
	if err = binary.Write(writer, binary.BigEndian, rec.Version); err != nil {
		return nil, err
	}

	// Length
	if rec.Length == 0 {
		if len(rec.Vector) > 0 {
			rec.Length += uint16(bodyVectorLen)
		}
		rec.Length += uint16(len(rec.Msg))
	}

	if err = binary.Write(writer, binary.BigEndian, rec.Length); err != nil {
		return nil, err
	}

	// Vector
	if len(rec.Vector) > 0 {
		if _, err = writer.Write(rec.Vector); err != nil {
			return nil, err
		}
	}

	// Body
	if _, err = writer.Write(rec.Msg); err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func ReadAndDecode(r io.Reader) (*Record, error) {
	rec, errDecode := Decode(r)
	if errDecode != nil {
		logger.Error("unable to decode record, %v", errDecode)
		return nil, errDecode
	}

	logger.Debug("received record: %+v", rec)
	return rec, nil
}

func EncodeAndWrite(w io.Writer, rec *Record) error {
	logger.Debug("encoding record: %+v", rec)

	reply, errEncode := Encode(rec)
	if errEncode != nil {
		logger.Error("error on encoding, %v", errEncode)
		return errEncode
	}

	logger.Debug("sending message...")

	n, err := w.Write(reply)
	if err != nil {
		logger.Error("error on write, %v", err)
		return err
	}

	logger.Debug("message sent, %d of %d bytes", n, len(reply))
	return nil
}

func readDecodeMsg(r io.Reader, expectedType uint8) ([]byte, error) {
	name := typeNames[expectedType]

	logger.Debug("reading %s message...", name)

	rec, err := ReadAndDecode(r)
	if err != nil {
		logger.Error("error on package decode, %v", err)
		return nil, fmt.Errorf("error on read %s package, %v", name, err)
	}

	if rec.Type != expectedType {
		return nil, fmt.Errorf("non %s message received, %+v", name, rec)
	}

	logger.Debug("message, %v", rec.Msg)
	return rec.Msg, nil
}
