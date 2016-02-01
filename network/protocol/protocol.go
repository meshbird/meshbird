package protocol

import (
	"errors"
	"github.com/meshbird/meshbird/log"
)

const (
	TypeHandshake uint8 = iota
	TypeOk
	TypeHeartbeat
	TypeTransfer
	TypePeerInfo
)

const (
	CurrentVersion = 1
	bodyVectorLen  = 16
)

var (
	logger = log.L("proto")

	ErrorUnableToReadVector  = errors.New("unable to read vector")
	ErrorUnableToReadMessage = errors.New("unable to read message")
	ErrorUnknownType         = errors.New("unknown type")

	knownTypes = []uint8{
		TypeHandshake,
		TypeOk,
		TypeHeartbeat,
		TypeTransfer,
		TypePeerInfo,
	}

	typeNames = map[uint8]string{
		TypeHandshake: "Handshake",
		TypeOk:        "Ok",
		TypeHeartbeat: "Heartbeat",
		TypeTransfer:  "Transfer",
		TypePeerInfo:  "PeerInfo",
	}
)

func isKnownType(needle uint8) bool {
	for _, t := range knownTypes {
		if needle == t {
			return true
		}
	}
	return false
}
