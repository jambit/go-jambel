package jambel

import (
	"bytes"
	"io"
)

// readUntil is a thin function reads from Telnet session. expect is
// a string used as signal to stop reading.
func readUntil(conn io.Reader, expect []byte) (out []byte, err error) {
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
