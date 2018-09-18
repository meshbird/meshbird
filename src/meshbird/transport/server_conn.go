package transport

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"io"
	"log"
	"net"

	"meshbird/utils"
)

type ServerConn struct {
	conn    *net.TCPConn
	key     string
	nonce   []byte
	buf     []byte
	aesgcm  cipher.AEAD
	handler ServerHandler
}

func NewServerConn(conn *net.TCPConn, key string, handler ServerHandler) *ServerConn {
	return &ServerConn{
		conn:    conn,
		key:     key,
		handler: handler,
		nonce:   make([]byte, 12),
		buf:     make([]byte, 65536),
	}
}

func (sc *ServerConn) run() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("server conn run err: %s", err)
		}
	}()
	key := utils.SHA256([]byte(sc.key))
	block, err := aes.NewCipher(key)
	utils.POE(err)
	sc.aesgcm, err = cipher.NewGCM(block)
	utils.POE(err)
	for {
		data, err := sc.read()
		if err != nil {
			log.Printf("server conn read err: %s", err)
			return
		}
		if sc.handler != nil {
			sc.handler.OnData(data)
		}
	}
}

func (sc *ServerConn) read() ([]byte, error) {
	var err error
	var n int
	_, err = io.ReadFull(sc.conn, sc.nonce)
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			return nil, nil
		}
		return nil, err
	}
	var msgLen uint16
	err = binary.Read(sc.conn, binary.LittleEndian, &msgLen)
	if err != nil {
		return nil, err
	}
	n, err = io.ReadFull(sc.conn, sc.buf[:msgLen])
	if err != nil {
		return nil, err
	}
	plain, err := sc.aesgcm.Open(nil, sc.nonce, sc.buf[:n], nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}
