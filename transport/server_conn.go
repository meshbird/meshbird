package transport

import (
	"crypto/cipher"
	"encoding/binary"
	"io"
	"log"
	"net"
	"bufio"

	"github.com/meshbird/meshbird/utils"
)

type ServerConn struct {
	conn    *net.TCPConn
	key     string
	nonce   []byte
	buf     []byte
	aesgcm  cipher.AEAD
	handler ServerHandler
	reader  *bufio.Reader
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
		sc.conn.Close()
	}()
	var err error
	err = sc.crypto()
	utils.POE(err)

	sc.conn.SetReadBuffer(1024 * 1024)
        sc.conn.SetWriteBuffer(1024 * 1024)
        sc.conn.SetNoDelay(true)


	sc.reader = bufio.NewReader(sc.conn)
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

func (sc *ServerConn) crypto() error {
	if sc.key == "" {
		log.Printf("incoming encryption disabled for %s", sc.conn.RemoteAddr())
		return nil
	}
	var err error
	sc.aesgcm, err = makeAES128GCM(sc.key)
	return err
}

func (sc *ServerConn) read() ([]byte, error) {
	var err error
	var secure uint8 = 0
	reader := sc.reader
	err = binary.Read(reader, binary.LittleEndian, &secure)
	if err != nil {
		return nil, err
	}
	var dataLen uint16
	err = binary.Read(reader, binary.LittleEndian, &dataLen)
	if err != nil {
		return nil, err
	}
	_, err = io.ReadFull(reader, sc.buf[:dataLen])
	if err != nil {
		return nil, err
	}
	if secure == 0 {
		return sc.buf[:dataLen], err
	} else {
		_, err = io.ReadFull(reader, sc.nonce)
		if err != nil {
			return nil, err
		}
		plain, err := sc.aesgcm.Open(nil, sc.nonce, sc.buf[:dataLen], nil)
		if err != nil {
			return nil, err
		}
		return plain, nil
	}
}
