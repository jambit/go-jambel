package jambel

import (
	"errors"
	"sync"

	"github.com/reiver/go-telnet"
)

var (
	// ErrConnectionClosed is returned when trying to send data on a closed connection
	ErrConnectionClosed = errors.New("connection is closed")
)

type TelnetConnection struct {
	addr string
	conn *telnet.Conn
	mu   sync.Mutex
}

// Send implements the Connector interface.
func (c *TelnetConnection) Send(cmd []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return []byte{}, ErrConnectionClosed
	}

	_, err := c.conn.Write(cmd)
	if err != nil {
		return []byte{}, err
	}
	out, err := readUntil(c.conn, []byte("\n"))
	return out, err
}

// Close implements the Connector interface.
func (c *TelnetConnection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
}

func NewNetworkJambel(url string) (*Jambel, error) {
	conn, err := telnet.DialTo(url)
	if err != nil {
		return nil, err
	}

	telnetConn := &TelnetConnection{
		addr: url,
		conn: conn,
	}
	return &Jambel{conn: telnetConn}, nil
}
