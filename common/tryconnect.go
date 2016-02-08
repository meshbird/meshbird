package common

import (
	"fmt"
	"github.com/anacrolix/utp"
	"github.com/meshbird/meshbird/log"
	"github.com/meshbird/meshbird/network/protocol"
	"github.com/meshbird/meshbird/secure"
	"net"
	"strconv"
	"time"
)

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

	rn.logger = log.L(fmt.Sprintf("public %s", rn.publicAddress))
	rn.logger.Debug("trying to connect...")

	s, errSocket := utp.NewSocket("udp4", ":0")
	if errSocket != nil {
		rn.logger.Error("unable to crete a socket, %v", errSocket)
		return nil, errSocket
	}

	conn, errDial := s.DialTimeout(rn.publicAddress, 10*time.Second)
	if errDial != nil {
		rn.logger.Error("unable to dial, %s", errDial)
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

	rn.privateIP = net.IP(peerInfo)

	// create new logger
	log.RemoveLogger(rn.logger.Name())
	rn.logger = log.L(fmt.Sprintf(rnLoggerFormat, rn.privateIP.String()))

	if err := protocol.WriteEncodePeerInfo(rn.conn, ln.State().PrivateIP()); err != nil {
		return nil, err
	}

	rn.logger.Info("connected, with public address %q", rn.publicAddress)
	return rn, nil
}
