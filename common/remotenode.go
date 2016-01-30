package common

import (
	"fmt"
	"github.com/anacrolix/utp"
	"github.com/meshbird/meshbird/log"
	"github.com/meshbird/meshbird/network/protocol"
	"github.com/meshbird/meshbird/secure"
	"io"
	"net"
	"strconv"
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

func (rn *RemoteNode) SendPack(pack *protocol.Packet) (err error) {
	if err = protocol.EncodeAndWrite(rn.conn, pack); err != nil {
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

	iface, ok := ln.Service("iface").(*InterfaceService)
	if !ok {
		rn.logger.Error("interface service not found")
		return
	}

	rn.logger.Debug("listening...")

	for {
		pack, err := protocol.Decode(rn.conn)
		if err != nil {
			rn.logger.Error("decode error, %v", err)
			if err == io.EOF {
				break
			}
			continue
		}
		rn.logger.Debug("received, %+v", pack)

		switch pack.Data.Type {
		case protocol.TypeTransfer:
			rn.logger.Debug("Writing to interface...")
			payloadEncrypted := pack.Data.Msg.(protocol.TransferMessage).Bytes()
			payload, errDec := secure.DecryptIV(payloadEncrypted, ln.State().Secret.Key, ln.State().Secret.Key)
			if errDec != nil {
				rn.logger.Error("error on decrypt, %v", err)
				break
			}
			srcAddr := payload[12:16].(net.IPAddr)
			dstAddr := payload[16:20].(net.IPAddr)
			rn.logger.Debug("received packet from %s to %s", srcAddr.String(), dstAddr.String())
			iface.WritePacket(payload)
		case protocol.TypeHeartbeat:
			rn.logger.Debug("heardbeat received, %v", pack.Data.Msg)
			rn.lastHeartbeat = time.Now()
		}
	}
}

func TryConnect(h string, networkSecret *secure.NetworkSecret, ln *LocalNode) (*RemoteNode, error) {
	host, portStr, errSplit := net.SplitHostPort(h)
	if errSplit != nil {
		return nil, errSplit
	}

	port, errConvert := strconv.Atoi(portStr)
	if errConvert != nil {
		return nil, errConvert
	}

	rn := new(RemoteNode)
	rn.lastHeartbeat = time.Now()
	rn.publicAddress = fmt.Sprintf("%s:%d", host, port + 1)

	rn.logger = log.L(fmt.Sprintf("public %s", rn.publicAddress))
	rn.logger.Debug("trying to connect...")

	s, errSocket := utp.NewSocket("udp4", ":0")
	if errSocket != nil {
		rn.logger.Error("unable to crete a socket, %v", errSocket)
		return nil, errSocket
	}

	conn, errDial := s.DialTimeout(rn.publicAddress, 10 * time.Second)
	if errDial != nil {
		rn.logger.Error("unable to dial, %v", errDial)
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

	// create new logger
	log.RemoveLogger(rn.logger.Name())
	rn.logger = log.L(fmt.Sprintf(rnLoggerFormat, rn.privateIP.String()))

	if err := protocol.WriteEncodePeerInfo(rn.conn, ln.State().PrivateIP); err != nil {
		return nil, err
	}

	rn.logger.Info("connected, with public address %q", rn.publicAddress)
	return rn, nil
}
