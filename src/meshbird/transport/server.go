package transport

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Server struct {
	addr    string
	handler ServerHandler
	key     string

	listener *net.TCPListener
}

func NewServer(addr string, handler ServerHandler, key string) *Server {
	srv := &Server{
		addr:    addr,
		handler: handler,
		key:     key,
	}
	return srv
}

func (s *Server) Start() {
	go s.Run()
}

func (s *Server) Run() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("server panic: %s", err)
		}
		log.Printf("server closed")
	}()
	var err error
	for {
		err = s.listen()
		if err != nil {
			log.Printf("server listen err: %s", err)
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func (s *Server) listen() error {
	defer func() {
		log.Printf("server listener closed")
	}()
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("net resolve tcp addr err: %s", err)
	}
	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return fmt.Errorf("net listen tcp err: %s", err)
	}
	defer s.listener.Close()
	for {
		tcpConn, err := s.listener.AcceptTCP()
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
