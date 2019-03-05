package meshbird

import (
	"log"
	"net"
	"time"

	"meshbird/config"
	"meshbird/iface"
	"meshbird/protocol"
	"meshbird/transport"
	"meshbird/utils"

	"github.com/golang/protobuf/proto"
)

type Peer struct {
	publicAddrs []string
	config      config.Config
	client      *transport.Client
}

func NewPeer(publicAddrs []string, cfg config.Config, getRoutes func() []Route) *Peer {
	peer := &Peer{
		publicAddrs: publicAddrs,
		config:      cfg,
	}
	peer.client = transport.NewClient(publicAddrs, cfg.Key, cfg.TransportThreads)
	return peer
}

func (p *Peer) Start() {
	p.client.Start()
	go p.process()
}

func (p *Peer) process() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("peer process panic: %s", err)
		}
	}()
	tickerPing := time.NewTicker(time.Second)
	defer tickerPing.Stop()
	for range tickerPing.C {
		p.SendPing()
	}
}

func (p *Peer) SendPing() {
	ip, _, err := net.ParseCIDR(p.config.Ip)
	utils.POE(err)
	env := &protocol.Envelope{
		Type: &protocol.Envelope_Ping{
			Ping: &protocol.MessagePing{
				Timestamp:   time.Now().UnixNano(),
				PublicAddrs: p.config.PublicAddrs,
				IP:          ip.String(),
			},
		},
	}
	data, err := proto.Marshal(env)
	utils.POE(err)
	p.client.Write(data)
}

func (p *Peer) SendPacket(pkt iface.PacketIP) {
	data, _ := proto.Marshal(&protocol.Envelope{
		Type: &protocol.Envelope_Packet{
			Packet: &protocol.MessagePacket{Payload: pkt},
		},
	})
	p.client.Write(data)
}
