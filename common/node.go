package common

import (
	"time"
)

type Node struct {
	ID        string
	PublicKey []byte
	Host      string
	Port      int
	LastSeen  time.Time
}