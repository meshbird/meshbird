package common

import ()

type Node struct {
	ID         string
	PublicKey  []byte
	Host       string
	Port       int
	InternalIP string
	LastSeen   uint64
}
