package transport

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Server struct {
	addrs   []string
	handler ServerHandler
	key     string
}

func NewServer(addrs []string, handler ServerHandler, key string) *Server {
	srv := &Server{
		addrs:   addrs,
		handler: handler,
		key:     key,
	}
	return srv
}

func (s *Server) Start() {
	for _, addr := range s.addrs {
		log.Printf("run listener on %s", addr)
		go s.process(addr)
	}
}

func (s *Server) process(addr string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("public server panic: %s", err)
		}
		log.Printf("public server closed")
	}()
	for {
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			log.Printf("net resolve tcp addr err: %s", err)
			time.Sleep(time.Second * 5)
			continue
		}
		err = s.listen(tcpAddr)
		if err != nil {
			log.Printf("server listen err: %s", err)
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func (s *Server) listen(tcpAddr *net.TCPAddr) error {
	defer func() {
		log.Printf("server listener closed")
	}()
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return fmt.Errorf("net listen tcp err: %s", err)
	}
	defer listener.Close()
	for {
		tcpConn, err := listener.AcceptTCP()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			return fmt.Errorf("server accept err: %s", err)
		} else {
			log.Printf("new accept from %s", tcpConn.RemoteAddr())
			serverConn := NewServerConn(tcpConn, s.key, s.handler)
			go serverConn.run()
		}
	}
}
