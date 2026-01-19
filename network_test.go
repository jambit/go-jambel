package jambel

import (
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/reiver/go-telnet"
)

// startMockServer will start a Telnet server at a random free port.
func startMockServer() (server *telnet.Server, err error) {
	var handler telnet.Handler = telnet.EchoHandler
	var listener net.Listener

	listener, err = net.Listen("tcp", ":0")
	if err != nil {
		return
	}
	server = &telnet.Server{
		Handler: handler,
		Addr:    listener.Addr().String(),
	}
	err = server.Serve(listener)
	return
}

/*
// ServeTELNET implements the
func (srv *mockServer) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	p := buffer[:]

	var line []byte
	for {

		n, err := r.Read(p)
		if n > 0 {
			line = append(line, p[:n]...)
			if strings.HasSuffix(string(line), "\n") {
				ctx.Logger().Debugf("%#v", string(line))
				_, _ = w.Write([]byte("OK\n"))
				line = []byte("") // empty line
			}
		}

		if err != nil {
			break
		}
	}
}
*/

func TestTelnetConnection_Close(t *testing.T) {
	type fields struct {
		addr string
		conn *telnet.Conn
		mu   sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			addr, err := startMockServer()
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("%v\n", addr)

			c := &TelnetConnection{
				addr: tt.fields.addr,
				conn: tt.fields.conn,
				mu:   tt.fields.mu,
			}
			c.Close()
		})
	}
}

// TODO test bytes are written
// TODO test locked connection
// TODO test closed connection
// TODO test write error

func TestTelnetConnection_Send(t *testing.T) {
	server, err := startMockServer()
	if err != nil {
		t.Fatal(err)
	}
	conn, err := NewTelnetConnection(server.Addr)
	if err != nil {
		t.Fatal(err)
	}

	conn.Write([]byte("hello"))
}

func mustConnectToTelnet(addr string) *telnet.Conn {
	conn, err := telnet.DialTo(addr)
	if err != nil {
		panic(err)
	}
	return conn
}

func TestNewNetworkJambel(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    *Jambel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNetworkJambel(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNetworkJambel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNetworkJambel() got = %v, want %v", got, tt.want)
			}
		})
	}
}
