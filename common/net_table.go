package common

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/meshbird/meshbird/network/protocol"
	"github.com/meshbird/meshbird/secure"
	"net"
	"os"
	"sync"
	"time"
)

type NetTable struct {
	BaseService

	localNode *LocalNode
	waitGroup sync.WaitGroup
	dhtInChan chan string

	lock      sync.RWMutex
	blackList map[string]time.Time
	peers     map[string]*RemoteNode

	heartbeatTicker <-chan time.Time

	logger *log.Logger
}

func (nt NetTable) Name() string {
	return "net-table"
}

func (nt *NetTable) Init(ln *LocalNode) error {
	// TODO: Add prefix
	nt.logger = log.New()
	nt.logger.Level = ln.config.Loglevel
	nt.localNode = ln
	nt.dhtInChan = make(chan string, 10)
	nt.blackList = make(map[string]time.Time)
	nt.peers = make(map[string]*RemoteNode)
	return nil
}

func (nt *NetTable) Run() error {
	for i := 0; i < 10; i++ {
		go nt.processDHTIn()
	}
	go nt.heartbeat()
	return nil
}

func (nt *NetTable) Stop() {
	nt.SetStatus(StatusStopping)
	nt.lock.Lock()
	defer nt.lock.Unlock()
	for _, peer := range nt.peers {
		peer.Close()
	}
}

func (nt *NetTable) GetDHTInChannel() chan<- string {
	return nt.dhtInChan
}

func (nt *NetTable) RemoteNodeByIP(ip net.IP) *RemoteNode {
	nt.lock.Lock()
	defer nt.lock.Unlock()
	return nt.peers[ip.To4().String()]
}

func (nt *NetTable) AddRemoteNode(rn *RemoteNode) {
	nt.logger.WithFields(log.Fields{
		"priv": rn.privateIP.String(),
		"pub":  rn.publicAddress,
	}).Debug("Trying to add node ...")

	if nt.localNode.State().PrivateIP.Equal(rn.privateIP) {
		nt.logger.Debug("Found myself, node will not be added!")
		return
	}

	nt.lock.Lock()
	defer nt.lock.Unlock()
	nt.peers[rn.privateIP.To4().String()] = rn
	nt.logger.WithFields(log.Fields{
		"priv": rn.privateIP.String(),
		"pub":  rn.publicAddress,
	}).Info("Added remote node")
	go rn.listen(nt.localNode)
}

func (nt *NetTable) RemoveRemoteNode(ip net.IP) {
	nt.lock.Lock()
	defer nt.lock.Unlock()
	delete(nt.peers, ip.To4().String())
}

func (nt *NetTable) processDHTIn() {
	for nt.Status() != StatusStopping {
		select {
		case host, ok := <-nt.dhtInChan:
			if !ok {
				return
			}
			nt.lock.Lock()
			_, ok = nt.blackList[host]
			if !ok {
				for _, peer := range nt.peers {
					if peer.publicAddress == host {
						ok = true
						break
					}
				}
			}
			nt.lock.Unlock()

			if !ok {
				nt.tryConnect(host)
			}
		}
	}
}

func (nt *NetTable) heartbeat() {
	nt.heartbeatTicker = time.Tick(5 * time.Second)

	for {
		select {
		case _, ok := <-nt.heartbeatTicker:
			if !ok {
				break
			}
			nt.lock.Lock()
			for _, peer := range nt.peers {
				if err := peer.SendPack(protocol.NewHeartbeatMessage(nt.localNode.State().PrivateIP)); err != nil {
					nt.logger.WithError(err).Error("Error on send heartbeat")
				}
			}
			nt.lock.Unlock()
		}
	}
}

func (nt *NetTable) tryConnect(h string) {
	rn, err := TryConnect(h, nt.localNode.NetworkSecret(), nt.localNode)
	if err != nil {
		nt.addToBlackList(h)
		return
	}
	nt.logger.Debug("Adding remote node from try connect...")
	nt.AddRemoteNode(rn)
}

func (nt *NetTable) addToBlackList(h string) {
	nt.lock.Lock()
	defer nt.lock.Unlock()
	//nt.blackList[h] = time.Now()
}

func (nt *NetTable) SendPacket(dstIP net.IP, payload []byte) {

	nt.logger.WithFields(log.Fields{
		"len": len(payload),
		"dst": dstIP.String(),
	}).Debug("Sending packet....")

	rn := nt.RemoteNodeByIP(dstIP)
	if rn == nil {
		nt.logger.WithField("dst", dstIP.String()).Debug("Destination host unreachable")
		nt.logger.WithField("known", nt.knownHosts()).Debug("Known hosts")

		return
	}

	payloadEnc, err := secure.EncryptIV(payload, nt.localNode.State().Secret.Key, nt.localNode.State().Secret.Key)
	if err != nil {
		nt.logger.WithError(err).Error("Error on encrypt")
		return
	}

	if err := rn.SendToInterface(payloadEnc); err != nil {
		nt.logger.WithFields(log.Fields{
			"dst": dstIP.String(),
			"err": err,
		}).Error("Error on sending packet")
	}
}

func (nt *NetTable) PeerAddresses() map[string]string {
	nt.lock.Lock()
	defer nt.lock.Unlock()
	peers := make(map[string]string)
	for l, peer := range nt.peers {
		peers[l] = peer.publicAddress
	}
	return peers
}

func (nt *NetTable) knownHosts() []string {
	nt.lock.Lock()
	defer nt.lock.Unlock()
	ips := make([]string, len(nt.peers))
	var i int
	for k := range nt.peers {
		ips[i] = k
		i++
	}
	return ips
}
