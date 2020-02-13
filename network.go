package jambel

import (
	"github.com/reiver/go-telnet"
)

type TelnetConnection struct {
	addr string
}

func (jmb *TelnetConnection) Send(cmd []byte) error {
	conn, err := telnet.DialTo(jmb.addr)
	if err != nil {
		return err
	}
	_, err = conn.Write(cmd)
	return err
}

func (jmb *TelnetConnection) Close() {
	// nothing to do here since we do not keep the connection open
}

func NewNetworkJambel(url string) *Jambel {
	conn := &TelnetConnection{url}
	return &Jambel{Connection: conn}
}
