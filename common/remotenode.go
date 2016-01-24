package common

import (
	"fmt"
	"github.com/anacrolix/utp"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gophergala2016/meshbird/network/protocol"
	"github.com/gophergala2016/meshbird/secure"
)

type RemoteNode struct {
	Node
	conn net.Conn
}

func TryConnect(h string, networkSecret *secure.NetworkSecret) (*RemoteNode, error) {
	host, portStr, errSplit := net.SplitHostPort(h)
	if errSplit != nil {
		return nil, errSplit
	}

	port, errConvert := strconv.Atoi(portStr)
	if errConvert != nil {
		return nil, errConvert
	}

	log.Printf("Trying to connection to: %s:%d", host, port+1)

	s, errSocket := utp.NewSocket("udp4", ":0")
	if errSocket != nil {
		log.Printf("Unable to crete a socket: %s", errSocket)
		return nil, errSocket
	}

	conn, errDial := s.DialTimeout(fmt.Sprintf("%s:%d", host, port+1), 10*time.Second)
	if errDial != nil {
		log.Printf("Unable to dial: %s", errDial)
		return nil, errDial
	}

	rn := new(RemoteNode)
	rn.conn = conn

	sessionKey := RandomBytes(16)

	pack := protocol.NewHandshakePacket(sessionKey, networkSecret)
	data, errEncode := protocol.Encode(pack)
	if errEncode != nil {
		log.Printf("Error on handshake generate: %s", errEncode)
		return nil, errEncode
	}

	if _, err := rn.conn.Write(data); err != nil {
		log.Printf("Error on write: %v", err)
	}

	buf := make([]byte, 1500)
	n, errRead := rn.conn.Read(buf)
	if errRead != nil {
		if errRead != io.EOF {
			log.Printf("Error on read from connection: %s", errRead)
		}
		return nil, errRead
	}

	buf = buf[:n]
	log.Printf("Readed: %v", buf)

	pack, errDecode := protocol.Decode(buf, sessionKey)
	if errDecode != nil {
		log.Printf("Unable to decode packet: %s", errDecode)
		return nil, errDecode
	}

	log.Printf("Packet message: %v", pack.Data.Msg)

	return rn, nil
}
