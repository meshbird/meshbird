package transport

import (
	"sync"
	"sync/atomic"

	"meshbird/config"
)

type Client struct {
	remoteAddr string
	conns      []*ClientConn
	mutex      sync.RWMutex
	config     config.Config
	serial     int64
}

func NewClient(remoteAddr string, cfg config.Config) *Client {
	return &Client{
		remoteAddr: remoteAddr,
		config:     cfg,
	}
}

func (c *Client) Start() {
	c.mutex.Lock()
	c.conns = make([]*ClientConn, c.config.TransportThreads)
	for connIndex := 0; connIndex < c.config.TransportThreads; connIndex++ {
		conn := NewClientConn(c.remoteAddr, c.config.Key, connIndex)
		c.conns[connIndex] = conn
		go c.conns[connIndex].run()
	}
	c.mutex.Unlock()
}

func (c *Client) WriteNow(data []byte) {
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