package transport

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"meshbird/utils"
)

var nilBuf = make([]byte, 0)

type ClientConn struct {
	remoteAddr string
	key        string
	conn       *net.TCPConn
	index      int
	mutex      sync.RWMutex
	aesgcm     cipher.AEAD
	chanWrite  chan []byte
}

func NewClientConn(remoteAddr, key string, index int) *ClientConn {
	return &ClientConn{
		remoteAddr: remoteAddr,
		key:        key,
		index:      index,
		chanWrite:  make(chan []byte),
	}
}

func (cc *ClientConn) tryConnect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", cc.remoteAddr)
	if err != nil {
		return err
	}

	cc.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	//cc.conn.SetReadBuffer(4096)
	//cc.conn.SetWriteBuffer(4096)
	cc.conn.SetNoDelay(true)

	if err != nil {
		return err
	}

	return nil
}

func (cc *ClientConn) run() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("transport thread %d addr %s panic: %s", cc.index, cc.remoteAddr, err)
		}
	}()
	cc.crypto()
	for {
		//log.Printf("transport thread %d to %s", cc.iIndex, cc.remoteAddr)
		err := cc.tryConnect()
		if err != nil {
			log.Printf("transport thread %d addr %s connect err: %s", cc.index, cc.remoteAddr, err)
			time.Sleep(time.Millisecond * 1000)
		} else {
			if err == nil {
				err = cc.process()
			} else {
				log.Printf("client err: %s", err)
			}
		}
	}
}

func (cc *ClientConn) process() (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			err = fmt.Errorf("client process panic: %s", perr)
		}
		cc.conn.Close()
		log.Printf("conn closed %s", err)
	}()
	log.Printf("connection good to %s : %d", cc.remoteAddr, cc.index)
	pingTicker := time.NewTicker(time.Second * 1)
	defer pingTicker.Stop()
	for {
		select {
		case <-pingTicker.C:
			err = cc.write(nilBuf)
		case buf := <-cc.chanWrite:
			err = cc.write(buf)
		}
		if err != nil {
			return err
		}
	}
}

func (cc *ClientConn) crypto() error {
	key := utils.SHA256([]byte(cc.key))
	//log.Printf("CLIENT KEY %s", utils.B64(key))
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	cc.aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		return err
	}
	return nil
}

func (cc *ClientConn) write(data []byte) error {
	if cc.conn == nil {
		return fmt.Errorf("no connection")
	}
	nonce := make([]byte, 12)
	_, err := io.ReadFull(crand.Reader, nonce)
	if err != nil {
		return err
	}
	data = cc.aesgcm.Seal(nil, nonce, data, nil)
	dataLen := uint16(len(data))
	//log.Printf("write %d %s %s", dataLen, utils.B64(nonce), utils.B64(data))
	_, err = cc.conn.Write(nonce)
	if err == nil {
		err = binary.Write(cc.conn, binary.LittleEndian, &dataLen)
		if err == nil {
			_, err = cc.conn.Write(data)
		}
	}
	return err
}

func (cc *ClientConn) Write(data []byte) {
	cc.chanWrite <- data
}

func (cc *ClientConn) WriteNow(data []byte) error {
	return cc.write(data)
}
