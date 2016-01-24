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
	host, portStr, err := net.SplitHostPort(h)
	if err != nil {
		return nil, nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, nil
	}

	conn, err := utp.DialTimeout(fmt.Sprintf("%s:%d", host, port+1), 10*time.Second)
	if err != nil {
		return nil, err
	}

	rn := new(RemoteNode)
	rn.conn = conn

	sessionKey := RandomBytes(16)

	pack := protocol.NewHandshakePacket(sessionKey, networkSecret)
	data, err := protocol.Encode(pack)
	if err != nil {
		log.Printf("Error on handshake generate: %s", err)
		return nil, nil
	}

	rn.conn.Write(data)

	buf := make([]byte, 1500)
	n, err := rn.conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Printf("Error on read from connection: %s", err)
		}
		return nil, nil
	}

	buf = buf[:n]
	//pack, err := protocol.Decode(buf, sessionKey)

	return rn, nil
}
