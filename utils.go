package jambel

import (
	"bytes"
	"context"
	"io"
	"time"
)

// readUntil is a thin function that reads from an io.Reader. expect is
// a byte slice used as signal to stop reading.
func readUntil(conn io.Reader, expect []byte) (out []byte, err error) {
	recvData := make([]byte, 1)
	var n int

	// Create a context with timeout to prevent infinite hangs
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Channel to signal when read is complete
	done := make(chan struct{})
	
	go func() {
		defer close(done)
		for {
			n, err = conn.Read(recvData)
			if err != nil {
				return
			}
			if n <= 0 {
				return
			}
			out = append(out, recvData...)
			if bytes.Contains(out, expect) {
				return
			}
		}
	}()

	select {
	case <-done:
		return out, err
	case <-ctx.Done():
		return out, ctx.Err()
	}
}
