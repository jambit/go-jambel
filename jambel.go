package jambel

import "fmt"

const (
	Red Colour = iota + 1
	Yellow
	Green
)

// LightState represents the state of a light module.
type LightState int

const (
	Off LightState = iota
	On
	Blink
	Flash
	BlinkInverse
)

// Connector implements how to communicate with a Jambel device.
type Connector interface {
	// Send sends a single command to the Jambel. Make sure that it is
	// terminated with "\n".
	Send(cmd []byte) error

	// Close closes the connection with the Jambel and cleans up afterwards.
	Close()
}

type Jambel struct {
	Connection ConnectionType
}

// Reset resets Jambel to all lights off
func (jmb *Jambel) Reset() error {
	return jmb.Connection.Send([]byte("reset\n"))
}

// On switches [colour] on where colour is one of GREEN, RED or YELLOW
func (jmb *Jambel) On(colour int) error {
	cmd := fmt.Sprintf("set=%d,on\n", colour)
	return jmb.Connection.Send([]byte(cmd))
}

// Off switches colour module off.
func (jmb *Jambel) Off(colour Colour) error {
	cmd := fmt.Sprintf("set=%d,off\n", colour)
	return jmb.Connection.Send([]byte(cmd))
}

func (jmb *Jambel) Blink(colour int) error {
	cmd := fmt.Sprintf("set=%d,blink\n", colour)
	return jmb.Connection.Send([]byte(cmd))
}

func (jmb *Jambel) BlinkInverse(colour int) error {
	cmd := fmt.Sprintf("set=%d,blink_inverse\n", colour)
	return jmb.Connection.Send([]byte(cmd))
}

func (jmb *Jambel) Flash(colour int) error {
	cmd := fmt.Sprintf("set=%d,flash\n", colour)
	return jmb.Connection.Send([]byte(cmd))
}

func (jmb *Jambel) SetAll(green, yellow, red int) error {
	cmd := fmt.Sprintf("set_all=%d,%d,%d,0\n", red, yellow, green)
	return jmb.Connection.Send([]byte(cmd))
}
