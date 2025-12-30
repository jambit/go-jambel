package main

import (
	"log"
	"strings"

	"github.com/reiver/go-telnet"
	"github.com/sirupsen/logrus"
)

type MockServer struct {
	Addr string
}

// ServeTELNET implements the
func (srv *MockServer) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

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

func (srv *MockServer) Run() error {
	log.Printf("Listening on %s...\n", srv.Addr)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	server := &telnet.Server{Addr: srv.Addr, Handler: srv, Logger: logger}
	return server.ListenAndServe()
}

func main() {
	srv := &MockServer{Addr: "127.0.0.1:5555"}
	err := srv.Run()
	if err != nil {
		panic(err)
	}
}
