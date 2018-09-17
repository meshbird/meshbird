package transport

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"

	"meshbird/utils"
)

type Server struct {
	addr    string
	handler ServerHandler
	key     string
	aesgcm  cipher.AEAD
}

func NewServer(addr string, handler ServerHandler, key string) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
		key:     key,
	}
}

func (s *Server) Start() {
	go s.Run()
}

func (s *Server) Run() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("server panic: %s", err)
		}
	}()

	key := utils.SHA256([]byte(s.key))
	//log.Printf("SERVER KEY %s", utils.B64(key))
	block, err := aes.NewCipher(key)
	utils.POE(err)
	s.aesgcm, err = cipher.NewGCM(block)
	utils.POE(err)

	for {
		err := s.listen()
		if err != nil {
			log.Printf("server listen err: %s", err)
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func (s *Server) listen() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.addr)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	defer listener.Close()
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			return err
		}
		go s.handleConn(tcpConn)
	}
}

func (s *Server) handleConn(tcpConn *net.TCPConn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("server handle conn err: %s", err)
		}
		log.Printf("server client conn closed")
	}()
	defer tcpConn.Close()
	//tcpConn.SetReadBuffer(4096)
	//tcpConn.SetWriteBuffer(4096)
	tcpConn.SetNoDelay(true)

	for {
		data, err := s.read(tcpConn)
		if err != nil {
			log.Printf("read message err: %s", err)
			return
		}
		s.handler.OnData(data)
		//log.Printf("new msg from %s: %#v", tcpConn.RemoteAddr(), msg)
	}
}

func (s *Server) read(conn *net.TCPConn) ([]byte, error) {
	nonce := make([]byte, 12)
	_, err := io.ReadFull(conn, nonce)
	if err != nil {
		return nil, err
	}

	var msgLen uint16
	err = binary.Read(conn, binary.LittleEndian, &msgLen)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, int(msgLen))
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}
	//log.Printf("read %d %s %s", msgLen, utils.B64(nonce), utils.B64(buf))
	plain, err := s.aesgcm.Open(nil, nonce, buf, nil)
	if err != nil {
		return nil, err
	}
	// ep := Envelope{}
	// err = proto.Unmarshal(plain, &ep)
	// return &ep, err
	return plain, nil
}
