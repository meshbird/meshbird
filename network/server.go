package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/anacrolix/utp"
	"github.com/gophergala2016/meshbird/network"
	"github.com/hsheth2/water"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go server(&wg, "0.0.0.0", 6000)
	wg.Wait()
}

func server(wg *sync.WaitGroup, host string, port int) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("error: %s", r)
		}
		wg.Done()
	}()
	log.Printf("server listen %s:%d", host, port)
	s, err := utp.NewSocket("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}
	iface, err := network.CreateTunInterfaceWithIp("tun0", "10.0.0.1/24")
	go func() {
		for {
			raw_data, err := network.NextNetworkPacket(iface)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if r := recover(); r != nil {
					log.Printf("error: %s", r)
				}
				log.Printf("disconnected")
				wg.Done()
			}()
			conn, err := utp.DialTimeout(fmt.Sprintf("%s:%d", "192.168.66.43", 6000), time.Second)
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			log.Printf("connected")
			_, err = conn.Write(raw_data)
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			//		chStat <- n
		}
	}()
	if err != nil {
		log.Println(err)
	}
	defer s.Close()
	for {
		conn, err := s.Accept()
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go readConn(conn, iface)
	}

}

func readConn(conn net.Conn, iface *water.Interface) {
	defer conn.Close()
	defer log.Printf("client %s disconnected", conn.RemoteAddr().String())
	log.Printf("client %s connected", conn.RemoteAddr().String())
	data := make([]byte, 1500)
	_, err := conn.Read(data)
	if err != nil {
		fmt.Println(err)
	}
	iface.Write(data)
}
