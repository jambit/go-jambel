package jambel

import (
	"bytes"
	"context"
	"io"
	"time"
)

const (
	// readTimeout is the maximum time to wait for a read operation to complete
	readTimeout = 5 * time.Second
	
	// goroutineCleanupTimeout is the maximum time to wait for goroutine cleanup
	goroutineCleanupTimeout = 10 * time.Millisecond
)

// readUntil is a thin function that reads from an io.Reader. expect is
// a byte slice used as signal to stop reading.
func readUntil(conn io.Reader, expect []byte) (out []byte, err error) {
	recvData := make([]byte, 1)

	// Create a context with timeout to prevent infinite hangs
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	defer cancel()

	// Channel to signal when read is complete
	type result struct {
		data []byte
		err  error
	}
	done := make(chan result, 1)
	
	go func() {
		var localOut []byte
		var localErr error
		
		for {
			// Check if context is canceled before continuing
			select {
			case <-ctx.Done():
				// Stop reading if context is done
				done <- result{data: localOut, err: ctx.Err()}
				return
			default:
			}
			
			n, readErr := conn.Read(recvData)
			if readErr != nil {
				localErr = readErr
				break
			}
			localOut = append(localOut, recvData[:n]...)
			if bytes.Contains(localOut, expect) {
				break
			}
		}
		
		done <- result{data: localOut, err: localErr}
	}()

	select {
	case res := <-done:
		return res.data, res.err
	case <-ctx.Done():
		// Wait a bit for the goroutine to clean up
		select {
		case res := <-done:
			return res.data, res.err
		case <-time.After(goroutineCleanupTimeout):
			return nil, ctx.Err()
		}
	}
}
