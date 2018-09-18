package transport

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"meshbird/config"
)

type Client struct {
	remoteAddr string
	conns      []*ClientConn
	mutex      sync.RWMutex
	config     config.Config
	serial     int64
	wg         sync.WaitGroup
}

func NewClient(remoteAddr string, cfg config.Config) *Client {
	return &Client{
		remoteAddr: remoteAddr,
		config:     cfg,
		conns:      make([]*ClientConn, cfg.TransportThreads),
	}
}

func (c *Client) Start() {
	c.mutex.Lock()
	for connIndex := 0; connIndex < c.config.TransportThreads; connIndex++ {
		c.wg.Add(1)
		conn := NewClientConn(c.remoteAddr, c.config.Key, connIndex, &c.wg)
		c.conns[connIndex] = conn
		go c.conns[connIndex].run()
	}
	c.mutex.Unlock()
}

func (c *Client) Stop() {
	defer log.Printf("client stopped")
	c.mutex.RLock()
	for _, conn := range c.conns {
		conn.Close()
	}
	c.mutex.RUnlock()
	c.wg.Wait()
}

func (c *Client) ConnectWait() {
	for {
		count := 0
		c.mutex.RLock()
		for _, conn := range c.conns {
			if conn.IsConnected() {
				count++
			}
		}
		c.mutex.RUnlock()
		if count == c.config.TransportThreads {
			return
		}
		time.Sleep(time.Second)
	}
}

func (c *Client) WriteNow(data []byte) {
	if c.config.TransportThreads == 1 {
		conn := c.conns[0]
		if err := conn.WriteNow(data); err != nil {
			conn.Write(data)
		}
		return
	}
	serial := atomic.AddInt64(&c.serial, 1)
	next := int(serial) % c.config.TransportThreads
	conn := c.conns[next]
	if err := conn.WriteNow(data); err != nil {
		conn.Write(data)
	}
}

func (c *Client) Write(data []byte) {
	serial := atomic.AddInt64(&c.serial, 1)
	next := int(serial) % c.config.TransportThreads
	conn := c.conns[next]
	conn.Write(data)
}
