package portmapper

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type PortMapper struct {
	mutex          sync.Mutex
	pairs          [][]int
	stopChan       chan bool
	wg             sync.WaitGroup
	ifaceListeners map[string]*net.UDPConn
	baddr          *net.UDPAddr
}

func NewPortMapper() *PortMapper {
	return &PortMapper{
		stopChan:       make(chan bool),
		ifaceListeners: make(map[string]*net.UDPConn),
	}
}

func (pm *PortMapper) AddPair(localPort, publicPort int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.pairs = append(pm.pairs, []int{localPort, publicPort})
}

func (pm *PortMapper) RemovePair(localPort, publicPort int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	if len(pm.pairs) == 0 {
		return
	}
	pairs := [][]int{}
	for _, pair := range pm.pairs {
		if pair[0] == localPort && pair[1] == publicPort {
			continue
		}
		pairs = append(pairs, pair)
	}
	pm.pairs = pairs
}

func (pm *PortMapper) Start() {
	pm.wg.Add(1)
	go pm.run()
}

func (pm *PortMapper) Stop() {
	pm.stopChan <- true
	pm.wg.Wait()
}

func (pm *PortMapper) run() error {
	defer pm.wg.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var err error
	pm.baddr, err = net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	if err != nil {
		return err
	}
	for {
		select {
		case <-pm.stopChan:
			log.Printf("port mapper stop")
		case <-ticker.C:
			log.Printf("resend map requests")
			if err := pm.updateInterfaceListeners(); err != nil {
				log.Printf("update iface listeners err: %s", err)
			}
		}
	}
	return nil
}

func (pm *PortMapper) updateInterfaceListeners() error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	for _, iface := range ifaces {
		log.Printf("listen interface %s %s", iface.Name, iface.Flags.String())
		conn, found := pm.ifaceListeners[iface.Name]
		if found {
			continue
		}
		conn, err = net.ListenMulticastUDP("udp", &iface, pm.baddr)
		if err != nil {
			return fmt.Errorf("listen on %s err: %s", iface.Name, err)
		}
		pm.ifaceListeners[iface.Name] = conn
	}
	return nil
}

func (pm *PortMapper) mapAll() error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, iface := range ifaces {
		go pm.mapIface(iface)
	}
	return nil
}

func (pm *PortMapper) mapIface(iface net.Interface) error {
	searchGateway(iface)
	return nil
}

func searchGateway(iface net.Interface) error {
	//request := "M-SEARCH * HTTP/1.1\r\n" +
	//	"HOST: 239.255.255.250:1900\r\n" +
	//	"ST: urn:schemas-upnp-org:service:WANIPConnection:1\r\n" +
	//	"MAN: \"ssdp:discover\"\r\n" + "MX: 3\r\n\r\n"
	return nil
}

func localAdresses() ([]net.Addr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	addrs := []net.Addr{}
	for _, iface := range ifaces {
		ifaceAddrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, ifaceAddrs...)
	}
	return addrs, nil
}

func localBroadcastAddrs() ([]net.Addr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	addrs := []net.Addr{}
	for _, iface := range ifaces {
		ifaceAddrs, err := iface.MulticastAddrs()
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, ifaceAddrs...)
	}
	return addrs, nil
}
