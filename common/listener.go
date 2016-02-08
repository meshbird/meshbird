package common

import (
	"fmt"
	"github.com/anacrolix/utp"
	"github.com/meshbird/meshbird/log"
	"github.com/meshbird/meshbird/network/protocol"
	"net"
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
	l.logger = log.L(l.Name())

	port := ln.State().ListenPort() + 1
	l.logger.Info("listening on %d port", port)
	socket, err := utp.NewSocket("udp4", fmt.Sprintf("0.0.0.0:%d", port))
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

		l.logger.Debug("new connection from %q", conn.RemoteAddr().String())

		if err = l.process(conn); err != nil {
			l.logger.Error("error on process, %v", err)
		}
	}
	return nil
}

func (l *ListenerService) Stop() {
	l.SetStatus(StatusStopping)
	l.socket.Close()
}

func (l *ListenerService) process(c net.Conn) error {
	handshakeMsg, errHandshake := protocol.ReadDecodeHandshake(c)
	if errHandshake != nil {
		return errHandshake
	}

	l.logger.Debug("processing hansdhake...")

	if !protocol.IsMagicValid(handshakeMsg) {
		return fmt.Errorf("invalid magic bytes")
	}

	l.logger.Debug("maginc is correct, replying...")

	if err := protocol.WriteEncodeOk(c); err != nil {
		return err
	}
	if err := protocol.WriteEncodePeerInfo(c, l.localNode.State().PrivateIP()); err != nil {
		return err
	}

	peerInfo, errPeerInfo := protocol.ReadDecodePeerInfo(c)
	if errPeerInfo != nil {
		return errPeerInfo
	}

	l.logger.Debug("processing peer info...")

	rn := NewRemoteNode(c, protocol.ExtractSessionKey(handshakeMsg), net.IP(peerInfo))

	l.logger.Debug("adding remote node...")
	l.localNode.NetTable().AddRemoteNode(rn)

	return nil
}
