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
	"os"
	"io"
)

type RemoteNode struct {
	Node
	conn          net.Conn
	sessionKey    []byte
	privateIP     net.IP
	publicAddress string
	logger        *log.Logger
}

func NewRemoteNode(conn net.Conn, sessionKey []byte, privateIP net.IP) *RemoteNode {
	return &RemoteNode{
		conn:          conn,
		sessionKey:    sessionKey,
		privateIP:     privateIP,
		publicAddress: conn.RemoteAddr().String(),
		logger:        log.New(os.Stderr, fmt.Sprintf("[remote priv/%s] ", privateIP.To4().String()), log.LstdFlags),
	}
}

func (rn *RemoteNode) SendPacket(payload []byte) error {
	return protocol.WriteEncodeTransfer(rn.conn, payload)
}

func (rn *RemoteNode) listen(ln *LocalNode) {
	defer rn.logger.Printf("EXIT LISTEN")
	defer func() {
		netTable, ok := ln.Service("net-table").(*NetTable)
		if !ok || netTable == nil {
			return
		}
		netTable.RemoveRemoteNode(rn.privateIP)
	}()
	iface, ok := ln.Service("iface").(*InterfaceService)
	if !ok {
		rn.logger.Printf("InterfaceService not found")
		return
	}

	rn.logger.Printf("Listening...")

	for {
		pack, err := protocol.Decode(rn.conn)
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				rn.logger.Printf("Decode error: %v", err)
				break
			}
			continue
		}

		rn.logger.Printf("Received package: %+v", pack)

		switch pack.Data.Type {
		case protocol.TypeTransfer:
			rn.logger.Printf("Writing to interface ... ")
			iface.WritePacket(pack.Data.Msg.(protocol.TransferMessage).Bytes())
		}
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
	rn.logger = log.New(os.Stderr, fmt.Sprintf("[remote pub/%s] ", rn.publicAddress), log.LstdFlags)

	rn.logger.Printf("Trying to connection to: %s", rn.publicAddress)

	s, errSocket := utp.NewSocket("udp4", ":0")
	if errSocket != nil {
		rn.logger.Printf("Unable to crete a socket: %s", errSocket)
		return nil, errSocket
	}

	conn, errDial := s.DialTimeout(rn.publicAddress, 10*time.Second)
	if errDial != nil {
		rn.logger.Printf("Unable to dial to %s: %s", rn.publicAddress, errDial)
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
	rn.logger = log.New(os.Stderr, fmt.Sprintf("[remote priv/%s] ", rn.privateIP.To4().String()), log.LstdFlags)

	if err := protocol.WriteEncodePeerInfo(rn.conn, rn.privateIP); err != nil {
		return nil, err
	}

	rn.logger.Printf("Connected to node: %s/%s", rn.privateIP.String(), rn.publicAddress)

	return rn, nil
}
