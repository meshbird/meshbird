package common

import (
	"fmt"
	"github.com/anacrolix/utp"
	"github.com/meshbird/meshbird/network/protocol"
	log "github.com/mgutz/logxi/v1"
	"net"
	"os"
)

type ListenerService struct {
	BaseService

	localNode *LocalNode
	socket    *utp.Socket

	logger log.Logger
}

func (l ListenerService) Name() string {
	return "listener"
}

func (l *ListenerService) Init(ln *LocalNode) error {
	l.logger = log.NewLogger(log.NewConcurrentWriter(os.Stderr), "[listener] ")

	if l.logger.IsInfo() {
		l.logger.Info("Listening on port: %d", ln.State().ListenPort+1)
	}
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
			break
		}

		if l.logger.IsDebug() {
			l.logger.Debug("Has new connection: %s", conn.RemoteAddr().String())

		}

		if err = l.process(conn); err != nil {
			l.logger.Error("Error on process: %s", err)
		}
	}
	return nil
}

func (l *ListenerService) Stop() {
	l.SetStatus(StatusStopping)
	l.socket.Close()
}

func (l *ListenerService) process(c net.Conn) error {
	//defer c.Close()

	handshakeMsg, errHandshake := protocol.ReadDecodeHandshake(c)
	if errHandshake != nil {
		return errHandshake
	}

	if l.logger.IsDebug() {

		l.logger.Debug("Processing hansdhake...")
	}

	if !protocol.IsMagicValid(handshakeMsg.Bytes()) {
		return fmt.Errorf("Invalid magic bytes")
	}
	if l.logger.IsDebug() {

		l.logger.Debug("Magic bytes are correct. Preparing reply...")
	}

	if err := protocol.WriteEncodeOk(c); err != nil {
		return err
	}
	if err := protocol.WriteEncodePeerInfo(c, l.localNode.State().PrivateIP); err != nil {
		return err
	}

	peerInfo, errPeerInfo := protocol.ReadDecodePeerInfo(c)
	if errPeerInfo != nil {
		return errPeerInfo
	}

	if l.logger.IsDebug() {
		l.logger.Debug("Processing PeerInfo...")

	}

	rn := NewRemoteNode(c, handshakeMsg.SessionKey(), peerInfo.PrivateIP())

	if l.logger.IsDebug() {
		l.logger.Debug("Adding remote node from listener...")

	}
	l.localNode.NetTable().AddRemoteNode(rn)

	return nil
}
