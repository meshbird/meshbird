package meshbird

import (
	"log"
	"strings"
	"sync"

	"meshbird/config"
	"meshbird/iface"
	"meshbird/protocol"
	"meshbird/transport"

	"github.com/golang/protobuf/proto"
)

type App struct {
	config config.Config
	peers  map[string]*Peer
	routes map[string]Route
	mutex  sync.RWMutex
	server *transport.Server
	iface  *iface.Iface
}

func NewApp(config config.Config) *App {
	return &App{
		config: config,
		peers:  make(map[string]*Peer),
		routes: make(map[string]Route),
	}
}

func (a *App) Run() error {
	a.server = transport.NewServer(a.config.LocalAddr, a.config.LocalPrivateAddr, a, a.config.Key)
	err := a.bootstrap()
	if err != nil {
		return err
	}
	go a.server.Start()
	return a.runIface()
}

func (a *App) runIface() error {
	a.iface = iface.New("", a.config.Ip, a.config.Mtu)
	err := a.iface.Start()
	if err != nil {
		return err
	}
	pkt := iface.NewPacketIP(a.config.Mtu)
	if a.config.Verbose == 1 {
		log.Printf("interface name: %s", a.iface.Name())
	}
	for {
		n, err := a.iface.Read(pkt)
		if err != nil {
			return err
		}
		src := pkt.GetSourceIP().String()
		dst := pkt.GetDestinationIP().String()
		if a.config.Verbose == 1 {
			log.Printf("packet: src=%s dst=%s len=%d", src, dst, n)
		}
		a.mutex.RLock()
		peer, ok := a.peers[a.routes[dst].LocalAddr]
		a.mutex.RUnlock()
		if !ok {
			if a.config.Verbose == 1 {
				log.Printf("unknown destination, packet dropped")
			}
		} else {
			peer.SendPacket(pkt)
		}
	}
}

func (a *App) bootstrap() error {
	seedAddrs := strings.Split(a.config.SeedAddrs, ",")
	for _, seedAddr := range seedAddrs {
		parts := strings.Split(seedAddr, "/")
		seedDC := parts[0]
		seedAddr = parts[1]
		if seedAddr == a.config.LocalAddr {
			log.Printf("skip seed addr %s because it's local addr", seedAddr)
			continue
		}
		peer := NewPeer(seedDC, seedAddr, a.config, a.getRoutes)
		peer.Start()
		a.mutex.Lock()
		a.peers[seedAddr] = peer
		a.mutex.Unlock()
	}
	return nil
}

func (a *App) getRoutes() []Route {
	a.mutex.Lock()
	routes := make([]Route, len(a.routes))
	i := 0
	for _, route := range a.routes {
		routes[i] = route
		i++
	}
	a.mutex.Unlock()
	return routes
}

func (a *App) OnData(buf []byte) {
	ep := protocol.Envelope{}
	err := proto.Unmarshal(buf, &ep)
	if err != nil {
		log.Printf("proto unmarshal err: %s", err)
		return
	}
	switch ep.Type.(type) {
	case *protocol.Envelope_Ping:
		ping := ep.GetPing()
		//log.Printf("received ping: %s", ping.String())
		a.mutex.Lock()
		a.routes[ping.GetIP()] = Route{
			LocalAddr:        ping.GetLocalAddr(),
			LocalPrivateAddr: ping.GetLocalPrivateAddr(),
			IP:               ping.GetIP(),
			DC:               ping.GetDC(),
		}
		if _, ok := a.peers[ping.GetLocalAddr()]; !ok {
			var peer *Peer
			if a.config.Dc == ping.GetDC() {
				peer = NewPeer(ping.GetDC(), ping.GetLocalPrivateAddr(),
					a.config, a.getRoutes)
			} else {
				peer = NewPeer(ping.GetDC(), ping.GetLocalAddr(),
					a.config, a.getRoutes)
				peer.Start()
			}
			a.peers[ping.GetLocalAddr()] = peer
			if a.config.Verbose == 1 {
				log.Printf("new peer %s", ping)
			}
		}
		if a.config.Verbose == 1 {
			log.Printf("routes %s", a.routes)
		}
		a.mutex.Unlock()
	case *protocol.Envelope_Packet:
		pkt := iface.PacketIP(ep.GetPacket().GetPayload())
		if a.config.Verbose == 1 {
			log.Printf("received packet: src=%s dst=%s len=%d",
				pkt.GetSourceIP(), pkt.GetDestinationIP(), len(pkt))
		}
		a.iface.Write(pkt)
	}
}
