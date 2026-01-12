package jambel

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/reiver/go-telnet"
)

var (
	// ErrConnectionClosed is returned when trying to send data on a closed connection
	ErrConnectionClosed = errors.New("connection is closed")
	// ErrReadTimeout is returned when reading from the connection times out
	ErrReadTimeout = errors.New("read timeout")
	// ErrMaxBytesExceeded is returned when the maximum number of bytes is exceeded
	ErrMaxBytesExceeded = errors.New("maximum bytes exceeded")
)

const (
	// defaultReadTimeout is the default timeout for reading from the telnet connection
	defaultReadTimeout = 5 * time.Second
	// maxReadBytes is the maximum number of bytes to read before giving up
	maxReadBytes = 1024 * 1024 // 1MB
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
	out, err := telnetRead(c.conn, []byte("\n"))
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

// contextErrToReadErr converts context errors to our custom error types for consistency.
func contextErrToReadErr(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrReadTimeout
	}
	return err
}

// telnetRead reads from a Telnet session until the expected terminator is found.
// It implements timeout protection to prevent indefinite blocking if the connection
// stops responding or never sends the expected terminator.
func telnetRead(conn *telnet.Conn, expect []byte) (out []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultReadTimeout)
	defer cancel()

	return telnetReadWithContext(ctx, conn, expect)
}

// telnetReadWithContext reads from a Telnet session with context support.
// The context can be used to control timeout and cancellation.
func telnetReadWithContext(ctx context.Context, conn *telnet.Conn, expect []byte) ([]byte, error) {
	// Use a channel to communicate the result from the reader goroutine
	type result struct {
		data []byte
		err  error
	}
	resultChan := make(chan result, 1)

	// Start a goroutine to perform the actual reading
	go func() {
		var data []byte
		recvData := make([]byte, 1)
		
		for {
			// Check if context is cancelled before each read
			select {
			case <-ctx.Done():
				resultChan <- result{data: data, err: contextErrToReadErr(ctx.Err())}
				return
			default:
			}

			n, readErr := conn.Read(recvData)
			if readErr != nil {
				resultChan <- result{data: data, err: readErr}
				return
			}
			if n <= 0 {
				// No data read but no error - continue reading
				continue
			}
			
			data = append(data, recvData[:n]...)
			
			// Check if we've exceeded the maximum bytes limit
			if len(data) > maxReadBytes {
				resultChan <- result{data: data, err: ErrMaxBytesExceeded}
				return
			}
			
			// Check if we've found the expected terminator
			if bytes.Contains(data, expect) {
				resultChan <- result{data: data, err: nil}
				return
			}
		}
	}()

	// Wait for either the result or context cancellation
	select {
	case res := <-resultChan:
		return res.data, res.err
	case <-ctx.Done():
		return nil, contextErrToReadErr(ctx.Err())
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
