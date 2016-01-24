package common

import (
	"log"
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

	logger *log.Logger
}

func (nt NetTable) Name() string {
	return "net-table"
}

func (nt *NetTable) Init(ln *LocalNode) error {
	nt.logger = log.New(os.Stderr, "[net-table] ", log.LstdFlags)
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
	return nil
}

func (nt *NetTable) Stop() {
	nt.SetStatus(StatusStopping)
}

func (nt *NetTable) GetDHTInChannel() chan<- string {
	return nt.dhtInChan
}

func (nt *NetTable) RemoteNodeByIP(ip net.IP) *RemoteNode {
	nt.lock.Lock()
	defer nt.lock.Unlock()
	return nt.peers[ip.String()]
}

func (nt *NetTable) AddRemoteNode(rn *RemoteNode) {
	nt.lock.Lock()
	defer nt.lock.Unlock()

	nt.logger.Printf("Added remote node: %s/%s", rn.privateIP.String(), rn.publicAddress)

	go rn.listen(nt.localNode)
	nt.peers[rn.privateIP.String()] = rn

	for addr := range nt.peers {
		log.Printf("NET-TABLE ENTRY %s", addr)
	}
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

func (nt *NetTable) tryConnect(h string) {
	rn, err := TryConnect(h, nt.localNode.NetworkSecret())
	if err != nil {
		nt.addToBlackList(h)
		return
	}

	nt.logger.Printf("Adding remote node from try connect...")
	nt.AddRemoteNode(rn)
}

func (nt *NetTable) addToBlackList(h string) {
	nt.lock.Lock()
	defer nt.lock.Unlock()
	//nt.blackList[h] = time.Now()
}

func (nt *NetTable) SendPacket(dstIP net.IP, payload []byte) {
	nt.logger.Printf("Sending to %s packet len %d", dstIP.String(), len(payload))

	rn := nt.RemoteNodeByIP(dstIP)
	if rn == nil {
		nt.logger.Printf("Destination host unreachable: %s", dstIP.String())
		return
	}

	if err := rn.SendPacket(payload); err != nil {
		nt.logger.Printf("Send packet to %s err: %s", dstIP.String(), err)
	}
}
