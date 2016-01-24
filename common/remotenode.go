package common

import (
	"fmt"
	"github.com/anacrolix/utp"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gophergala2016/meshbird/network/protocol"
	"github.com/gophergala2016/meshbird/secure"
)

type RemoteNode struct {
	Node
	conn          net.Conn
	sessionKey    []byte
	privateIP     net.IP
	publicAddress string
}

func NewRemoteNode(conn net.Conn, sessionKey []byte, privateIP net.IP) *RemoteNode {
	return &RemoteNode{
		conn:          conn,
		sessionKey:    sessionKey,
		privateIP:     privateIP,
		publicAddress: conn.RemoteAddr().String(),
	}
}

func (rn *RemoteNode) SendPacket(dstIP net.IP, payload []byte) error {
	return protocol.WriteEncodeTransfer(rn.conn, payload)
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

	rn := new(RemoteNode)
	rn.publicAddress = fmt.Sprintf("%s:%d", host, port+1)

	log.Printf("Trying to connection to: %s", rn.publicAddress)

	s, errSocket := utp.NewSocket("udp4", ":0")
	if errSocket != nil {
		log.Printf("Unable to crete a socket: %s", errSocket)
		return nil, errSocket
	}

	conn, errDial := s.DialTimeout(rn.publicAddress, 10*time.Second)
	if errDial != nil {
		log.Printf("Unable to dial to %s: %s", rn.publicAddress, errDial)
		return nil, errDial
	}

	rn.conn = conn
	rn.sessionKey = RandomBytes(16)

	if err := protocol.WriteEncodeHandshake(rn.conn, rn.sessionKey, networkSecret); err != nil {
		return nil, err
	}
	if _, okError := protocol.ReadDecodeOk(rn.conn); okError != nil {
		return nil, okError
	}

	peerInfo, errPeerInfo := protocol.ReadDecodePeerInfo(rn.conn)
	if errPeerInfo != nil {
		return nil, errPeerInfo
	}

	rn.privateIP = peerInfo.PrivateIP()

	if err := protocol.WriteEncodePeerInfo(rn.conn, rn.privateIP); err != nil {
		return nil, err
	}

	log.Printf("Connected to node: %s/%s", rn.privateIP.String(), rn.publicAddress)

	return rn, nil
}
