package common

import (
	"fmt"
	"github.com/anacrolix/utp"
	"github.com/gophergala2016/meshbird/network/protocol"
	"log"
	"net"
)

type ListenerService struct {
	BaseService

	localNode *LocalNode
	socket    *utp.Socket
}

func (l ListenerService) Name() string {
	return "listener"
}

func (l *ListenerService) Init(ln *LocalNode) error {
	log.Printf("Listening on port: %d", ln.State().ListenPort+1)
	socket, err := utp.NewSocket("udp4", fmt.Sprintf("0.0.0.0:%d", ln.State().ListenPort+1))
	if err != nil {
		return err
	}

	l.localNode = ln
	l.socket = socket
	return nil
}

func (l *ListenerService) Run() error {
	for {
		conn, err := l.socket.Accept()
		if err != nil {
			continue
		}
		if err = l.process(conn); err != nil {
			log.Printf("Error on process: %s", err)
		}
	}
	return nil
}

func (l *ListenerService) Stop() {
	l.SetStatus(StatusStopping)
	l.socket.Close()
}

func (l *ListenerService) process(c net.Conn) error {
	defer c.Close()

	handshakeMsg, errHandshake := protocol.ReadDecodeHandshake(c, nil)
	if errHandshake != nil {
		return errHandshake
	}

	log.Println("Processing hansdhake...")

	if !protocol.IsMagicValid(handshakeMsg.Bytes()) {
		return fmt.Errorf("Invalid magic bytes")
	}

	log.Println("Magic bytes are correct. Preparing reply...")

	if err := protocol.WriteEncodeOk(c); err != nil {
		return err
	}
	if err := protocol.WriteEncodePeerInfo(c, l.localNode.State().PrivateIP); err != nil {
		return err
	}

	peerInfo, errPeerInfo := protocol.ReadDecodePeerInfo(c, nil)
	if errPeerInfo != nil {
		return errPeerInfo
	}

	log.Println("Processing PeerInfo...")

	rn := NewRemoteNode(c, handshakeMsg.SessionKey(), peerInfo.PrivateIP())

	netTable, ok := l.localNode.Service("net-table").(*NetTable)
	if !ok || netTable == nil {
		return fmt.Errorf("net-table is nil")
	}

	netTable.AddRemoteNode(rn)

	return nil
}
