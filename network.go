package jambel

import (
	"bytes"
	"errors"
	"fmt"
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

	fmt.Printf(">>> %s", cmd)
	_, err := c.conn.Write(cmd)
	if err != nil {
		return []byte{}, err
	}
	out, _ := telnetRead(c.conn, []byte("\n"))
	fmt.Printf("<<< %s", out)
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

// telnetRead is a thin function reads from Telnet session. expect is
// a string used as signal to stop reading.
func telnetRead(conn *telnet.Conn, expect []byte) (out []byte, err error) {
	recvData := make([]byte, 1)
	var n int

	for {
		n, err = conn.Read(recvData)
		if err != nil {
			return out, err
		}
		if n <= 0 {
			break
		}
		out = append(out, recvData...)
		if bytes.Contains(out, expect) {
			break
		}
	}
	return out, nil
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
