package common

import (
	"fmt"
	"github.com/meshbird/meshbird/log"
	"github.com/meshbird/meshbird/network/protocol"
	"github.com/meshbird/meshbird/secure"
	"net"
	"time"
)

var (
	rnLoggerFormat = "remote %s"
)

type RemoteNode struct {
	Node
	conn          net.Conn
	sessionKey    []byte
	privateIP     net.IP
	publicAddress string
	logger        log.Logger
	lastHeartbeat time.Time
}

func NewRemoteNode(conn net.Conn, sessionKey []byte, privateIP net.IP) *RemoteNode {
	return &RemoteNode{
		conn:          conn,
		sessionKey:    sessionKey,
		privateIP:     privateIP,
		publicAddress: conn.RemoteAddr().String(),
		logger:        log.L(fmt.Sprintf(rnLoggerFormat, privateIP.String())),
		lastHeartbeat: time.Now(),
	}
}

func (rn *RemoteNode) SendToInterface(payload []byte) error {
	return protocol.WriteEncodeTransfer(rn.conn, payload)
}

func (rn *RemoteNode) SendPack(rec *protocol.Record) (err error) {
	if err = protocol.EncodeAndWrite(rn.conn, rec); err != nil {
		err = fmt.Errorf("error on write transfer message, %v", err)
	}
	return
}

func (rn *RemoteNode) Close() {
	defer rn.conn.Close()
	rn.logger.Debug("closing...")
}

func (rn *RemoteNode) listen(ln *LocalNode) {
	defer rn.logger.Debug("listener stopped...")
	defer func() {
		ln.NetTable().RemoveRemoteNode(rn.privateIP)
	}()
	defer rn.conn.Close()

	iface, ok := ln.Service("iface").(*InterfaceService)
	if !ok {
		rn.logger.Error("interface service not found")
		return
	}

	rn.logger.Debug("listening...")

	for {
		rec, err := protocol.Decode(rn.conn)
		if err != nil {
			rn.logger.Error("decode error, %v", err)
			break
		}
		rn.logger.Debug("received, %+v", rec)

		switch rec.Type {
		case protocol.TypeTransfer:
			rn.logger.Debug("Writing to interface...")
			payload, errDec := secure.DecryptIV(rec.Msg, ln.State().Secret.Key, ln.State().Secret.Key)
			if errDec != nil {
				rn.logger.Error("error on decrypt, %v", err)
				break
			}
			srcAddress := net.IP(payload[12:16])
			dstAddress := net.IP(payload[16:20])
			rn.logger.Info("received packet from %s to %s", srcAddress.String(), dstAddress.String())
			if err = iface.WritePacket(payload); err != nil {
				rn.logger.Error("write packet err: %s", err)
			}

		case protocol.TypeHeartbeat:
			rn.logger.Debug("heardbeat received, %v", rec.Msg)
			rn.lastHeartbeat = time.Now()
		}
	}
}
