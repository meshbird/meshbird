package transport

import (
	"bytes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/meshbird/meshbird/utils"
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
	chanClose  chan bool
	wg         sync.WaitGroup
	parentWG   *sync.WaitGroup
	connected  bool
	nonce      []byte
	buf        *bytes.Buffer
}

func NewClientConn(remoteAddr, key string, index int, parentWG *sync.WaitGroup) *ClientConn {
	return &ClientConn{
		remoteAddr: remoteAddr,
		key:        key,
		index:      index,
		chanWrite:  make(chan []byte),
		chanClose:  make(chan bool),
		parentWG:   parentWG,
		nonce:      make([]byte, 12),
		buf:        &bytes.Buffer{},
	}
}

func (cc *ClientConn) IsConnected() bool {
	cc.mutex.RLock()
	connected := cc.connected
	cc.mutex.RUnlock()
	return connected
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

	cc.conn.SetReadBuffer(1024 * 1024)
	cc.conn.SetWriteBuffer(1024 * 1024)
	cc.conn.SetNoDelay(true)

	return nil
}

func (cc *ClientConn) crypto() (err error) {
	if cc.key == "" {
		return nil
	}
	cc.aesgcm, err = makeAES128GCM(cc.key)
	return
}

func (cc *ClientConn) run() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("transport thread %d addr %s panic: %s", cc.index, cc.remoteAddr, err)
		}
		cc.wg.Done()
	}()
	cc.wg.Add(1)
	var err error
	err = cc.crypto()
	utils.POE(err)
	for {
		select {
		case <-cc.chanClose:
			return
		default:
		}
		err := cc.tryConnect()
		if err != nil {
			log.Printf("transport thread %d addr %s connect err: %s", cc.index, cc.remoteAddr, err)
			time.Sleep(time.Millisecond * 1000)
		} else {
			if err == nil {
				err = cc.process()
				if err == nil {
					break
				}
			}
			if err != nil {
				log.Printf("client err: %s", err)
			}
		}
	}
}

func (cc *ClientConn) Close() {
	cc.chanClose <- true
	cc.wg.Wait()
	if cc.parentWG != nil {
		cc.parentWG.Done()
	}
}

func (cc *ClientConn) setConnected(value bool) {
	cc.mutex.Lock()
	cc.connected = value
	cc.mutex.Unlock()
}

func (cc *ClientConn) process() (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			err = fmt.Errorf("client process panic: %s", perr)
		}
		cc.setConnected(false)
		cc.conn.Close()
		log.Printf("client conn closed")
	}()
	cc.setConnected(true)
	log.Printf("connection good to %s : %d", cc.remoteAddr, cc.index)
	pingTicker := time.NewTicker(time.Second * 1)
	defer pingTicker.Stop()
	for {
		select {
		case <-cc.chanClose:
			return nil
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

func (cc *ClientConn) write(data []byte) error {
	if cc.conn == nil {
		return fmt.Errorf("no connection")
	}
	var err error
	cc.buf.Reset()
	var secure uint8 = 0
	if cc.aesgcm != nil {
		secure = 1
	}
	err = binary.Write(cc.buf, binary.LittleEndian, &secure)
	if err != nil {
		return err
	}
	if secure == 0 {
		dataLen := uint16(len(data))
		err = binary.Write(cc.buf, binary.LittleEndian, &dataLen)
		if err != nil {
			return err
		}
		_, err = cc.buf.Write(data)
		if err != nil {
			return err
		}
	} else {
		_, err = io.ReadFull(crand.Reader, cc.nonce)
		if err != nil {
			return err
		}
		data2 := cc.aesgcm.Seal(nil, cc.nonce, data, nil)
		dataLen := uint16(len(data2))
		err = binary.Write(cc.buf, binary.LittleEndian, &dataLen)
		if err != nil {
			return err
		}
		_, err = cc.buf.Write(data2)
		if err != nil {
			return err
		}
		_, err = cc.buf.Write(cc.nonce)
		if err != nil {
			return err
		}
	}
	if err == nil {
		_, err = cc.buf.WriteTo(cc.conn)
	}
	return err
}

func (cc *ClientConn) Write(data []byte) {
	cc.chanWrite <- data
}

func (cc *ClientConn) WriteNow(data []byte) error {
	return cc.write(data)
}
