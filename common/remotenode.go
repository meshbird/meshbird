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
		conn: conn,
		sessionKey: sessionKey,
		privateIP: privateIP,
		publicAddress: conn.RemoteAddr().String(),
	}
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

	pack := protocol.NewHandshakePacket(rn.sessionKey, networkSecret)
	data, errEncode := protocol.Encode(pack)
	if errEncode != nil {
		log.Printf("Error on handshake generate: %s", errEncode)
		return nil, errEncode
	}

	if _, err := rn.conn.Write(data); err != nil {
		log.Printf("Error on write: %v", err)
		return nil, err
	}

	pack, errDecode := protocol.ReadAndDecode(rn.conn, rn.sessionKey)
	if errDecode != nil {
		log.Printf("Unable to decode packet: %s", errDecode)
		return nil, errDecode
	}

	if pack.Data.Type != protocol.TypeOk {
		log.Printf("Got non OK message: %+v", pack)
		return nil, fmt.Errorf("Got non OK message")
	}

	pack, errDecode = protocol.ReadAndDecode(rn.conn, rn.sessionKey)
	if errDecode != nil {
		log.Printf("Unable to decode packet: %s", errDecode)
		return nil, errDecode
	}

	if pack.Data.Type != protocol.TypePeerInfo {
		log.Printf("Got non PeerInfo message: %+v", pack)
		return nil, fmt.Errorf("Got non PeerInfo message")
	}

	rn.privateIP = pack.Data.Msg.(protocol.PeerInfoMessage).PrivateIP()

	pack = protocol.NewPeerInfoMessage(rn.privateIP)
	data, errEncode = protocol.Encode(pack)
	if errEncode != nil {
		log.Printf("Error on PeerInfo generate: %s", errEncode)
		return nil, errEncode
	}

	if _, err := rn.conn.Write(data); err != nil {
		log.Printf("Error on write: %v", err)
		return nil, err
	}

	log.Printf("Connected to node: %s/%s", rn.privateIP.String(), rn.publicAddress)

	return rn, nil
}
