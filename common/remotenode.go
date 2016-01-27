package common

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/anacrolix/utp"
	"github.com/meshbird/meshbird/network/protocol"
	"github.com/meshbird/meshbird/secure"
)

type RemoteNode struct {
	Node
	conn          net.Conn
	sessionKey    []byte
	privateIP     net.IP
	publicAddress string
	logger        *log.Logger
	lastHeartbeat time.Time
}

func NewRemoteNode(conn net.Conn, sessionKey []byte, privateIP net.IP) *RemoteNode {
	return &RemoteNode{
		conn:          conn,
		sessionKey:    sessionKey,
		privateIP:     privateIP,
		publicAddress: conn.RemoteAddr().String(),
		logger:        log.New(),
		// TODO: Fix it
		//		logger:        log.NewLogger(log.NewConcurrentWriter(os.Stderr), fmt.Sprintf("[remote priv/%s] ", privateIP.To4().String())),
		lastHeartbeat: time.Now(),
	}
}

func (rn *RemoteNode) SendToInterface(payload []byte) error {
	return protocol.WriteEncodeTransfer(rn.conn, payload)
}

func (rn *RemoteNode) SendPack(pack *protocol.Packet) (err error) {
	if err = protocol.EncodeAndWrite(rn.conn, pack); err != nil {
		err = fmt.Errorf("Error on write Transfer message: %v", err)
	}
	return
}

func (rn *RemoteNode) Close() {
	defer rn.conn.Close()
	rn.logger.Debug("Closing...")
}

func (rn *RemoteNode) listen(ln *LocalNode) {
	defer rn.logger.Debug("EXIT LISTEN")
	defer func() {
		ln.NetTable().RemoveRemoteNode(rn.privateIP)
	}()

	iface, ok := ln.Service("iface").(*InterfaceService)
	if !ok {
		rn.logger.Error("InterfaceService not found")
		return
	}

	rn.logger.Info("Listening...")

	for {
		pack, err := protocol.Decode(rn.conn)
		if err != nil {
			rn.logger.Debug(fmt.Sprintf("Decode error: %v", err))
			if err == io.EOF {
				break
			}
			continue
		}
		rn.logger.Debug(fmt.Sprintf("Received package: %+v", pack))

		switch pack.Data.Type {
		case protocol.TypeTransfer:
			rn.logger.Debug("Writing to interface...")
			payloadEncrypted := pack.Data.Msg.(protocol.TransferMessage).Bytes()
			payload, errDec := secure.DecryptIV(payloadEncrypted, ln.State().Secret.Key, ln.State().Secret.Key)
			if errDec != nil {
				rn.logger.Debug(fmt.Sprintf("Error on decrypt: %v", errDec))
				break
			}
			iface.WritePacket(payload)
		case protocol.TypeHeartbeat:
			rn.logger.Debug(fmt.Sprintf("Received heardbeat... %v", pack.Data.Msg))
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
	rn.publicAddress = fmt.Sprintf("%s:%d", host, port+1)
	// TODO: Fix it
	//rn.logger = log.NewLogger(log.NewConcurrentWriter(os.Stderr), fmt.Sprintf("[remote priv/%s] ", fmt.Sprintf("[remote pub/%s] ", rn.publicAddress)))
	rn.logger = log.New()

	rn.logger.Debug(fmt.Sprintf("Trying to connection to: %s", rn.publicAddress))

	s, errSocket := utp.NewSocket("udp4", ":0")
	if errSocket != nil {
		rn.logger.Debug(fmt.Sprintf("Unable to crete a socket: %s", errSocket))
		return nil, errSocket
	}

	conn, errDial := s.DialTimeout(rn.publicAddress, 10*time.Second)
	if errDial != nil {
		rn.logger.Error(fmt.Sprintf("Unable to dial to %s: %s", rn.publicAddress, errDial))
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
	// TODO: Fix ot
	//rn.logger = log.NewLogger(log.NewConcurrentWriter(os.Stderr), fmt.Sprintf("[remote priv/%s] ", rn.privateIP.To4().String()))
	rn.logger = log.New()
	if err := protocol.WriteEncodePeerInfo(rn.conn, ln.State().PrivateIP); err != nil {
		return nil, err
	}
	rn.logger.Debug(fmt.Sprintf("Connected to node: %s/%s", rn.privateIP.String(), rn.publicAddress))

	return rn, nil
}
