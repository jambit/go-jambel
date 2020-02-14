package jambel

import (
	"github.com/reiver/go-telnet"
)

type TelnetConnection struct {
	addr string
}

func (c *TelnetConnection) Send(cmd []byte) error {
	conn, err := telnet.DialTo(c.addr)
	defer conn.Close()
	if err != nil {
		return err
	}
	_, err = conn.Write(cmd)
	return err
}

func (c *TelnetConnection) Close() {
	// nothing to do here since we do not keep the connection open
}

func NewNetworkJambel(url string) *Jambel {
	conn := &TelnetConnection{url}
	return &Jambel{Connection: conn}
}
