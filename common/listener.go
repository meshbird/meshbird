package common

import (
	"fmt"
	"github.com/anacrolix/utp"
	"github.com/gophergala2016/meshbird/network/protocol"
	"io"
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

	pack, errDecode := protocol.ReadAndDecode(c, nil)
	if errDecode != nil {
		return fmt.Errorf("Unable to decode packet: %s", errDecode)
	}

	log.Printf("Received: %+v", pack)

	if pack.Data.Type != protocol.TypeHandshake {
		return fmt.Errorf("Unexpected message type: %s", protocol.TypeName(pack.Data.Type))
	}

	log.Println("Processing hansdhake...")

	msg := pack.Data.Msg.(protocol.HandshakeMessage)
	if !protocol.IsMagicValid(msg.Bytes()) {
		return fmt.Errorf("Invalid magic bytes")
	}

	log.Println("Magic bytes are correct. Preparing reply...")

	reply, errEncode := protocol.Encode(
		protocol.NewOkMessage(l.localNode.State().PrivateIP),
	)

	if errEncode != nil {
		return fmt.Errorf("Error on encoding reply: %v", errEncode)
	}

	log.Printf("Sending reply...")

	if _, err := c.Write(reply); err != nil {
		return fmt.Errorf("Error on write reply: %v", err)
	}

	return nil
}
