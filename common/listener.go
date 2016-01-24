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

	buf := make([]byte, 1500)
	n, errRead := c.Read(buf)
	if errRead != nil {
		if errRead != io.EOF {
			log.Printf("Error on read from socket: %s", errRead)
			return errRead
		}
		return nil
	}

	pack, errDecode := protocol.Decode(buf[:n], nil)
	if errDecode != nil {
		log.Printf("Unable to decode packet: %s", errDecode)
		return errDecode
	}

	log.Printf("Received: %+v", pack)

	if pack.Data.Type == protocol.TypeHandshake {
		log.Println("Handshake!!!")

		msg := pack.Data.Msg.(protocol.HandshakeMessage)
		if protocol.IsMagicValid([]byte(msg)) {
			log.Println("Ja, ja supa magic!")

			replyPack := protocol.NewOkMessage(l.localNode.State().PrivateIP)
			reply, errEncode := protocol.Encode(replyPack)
			if errEncode != nil {
				log.Printf("Error on encode: %v", errEncode)
				return nil
			}

			log.Printf("Sending reply: %+v", replyPack)

			if _, err := c.Write(reply); err != nil {
				log.Printf("Erro on write: %v", err)
			}
		} else {
			log.Println("Magic is not valid")
		}
	} else {
		log.Printf("Expected handshake but got: %d", pack.Data.Type)
	}

	return nil
}
