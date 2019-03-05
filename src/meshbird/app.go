package meshbird

import (
	"fmt"
	"log"
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
	if a.config.HostAddr == "" {
		log.Printf("host_addr is empty")
		if len(a.config.PublicAddrs) == 0 || a.config.PublicAddrs[0] == "" {
			return fmt.Errorf("if host_addr is e16ympty, need public_addrs")
		}
		if len(a.config.BindAddrs) == 0 || a.config.BindAddrs[0] == "" {
			return fmt.Errorf("if host_addr is empty, need bind_addrs")
		}
	} else {
		a.config.PublicAddrs = []string{a.config.HostAddr}
		a.config.BindAddrs = []string{a.config.HostAddr}
	}
	if a.config.Key == "" {
		log.Printf("key is empty, encryption disabled")
	}
	log.Printf("run listeners on %s", a.config.BindAddrs)
	a.server = transport.NewServer(a.config.BindAddrs, a, a.config.Key)
	if len(a.config.SeedAddrs) == 0 || a.config.SeedAddrs[0] == "" {
		log.Printf("seed_addrs is empty, bootstrap disabled")
	} else {
		err := a.bootstrap()
		if err != nil {
			return err
		}
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
		peer, ok := a.findPeerByIP(dst)
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
	for _, seedAddr := range a.config.SeedAddrs {
		found := false
		for _, publicAddr := range a.config.PublicAddrs {
			if seedAddr == publicAddr {
				found = true
				break
			}
		}
		if found {
			log.Printf("skip seed addr %s because it's local addr", seedAddr)
			continue
		}
		peer := NewPeer([]string{seedAddr}, a.config, a.getRoutes)
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
			PublicAddrs: ping.GetPublicAddrs(),
			IP:          ping.GetIP(),
		}
		peer, ok := a.findPeerByPublicAddrs(ping.GetPublicAddrs())
		if !ok {
			var peer *Peer
			peer = NewPeer(ping.PublicAddrs, a.config, a.getRoutes)
			peer.Start()
			a.peers[ping.GetIP()] = peer
			if a.config.Verbose == 1 {
				log.Printf("new peer %s", ping)
			}
		} else {
			peer.publicAddrs = ping.GetPublicAddrs()
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

func (a *App) findPeerByPublicAddrs(addrs []string) (peer *Peer, ok bool) {
	for _, peer := range a.peers {
		for _, peerPublicAddr := range peer.publicAddrs {
			for _, addr := range addrs {
				if peerPublicAddr == addr {
					return peer, true
				}
			}
		}
	}
	return nil, false
}

func (a *App) findPeerByIP(ip string) (peer *Peer, ok bool) {
	peer, ok = a.peers[ip]
	return
}
