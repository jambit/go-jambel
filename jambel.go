package jambel

import "fmt"

// Jambel represents jambit's fast feedback device.
type Jambel struct {
	conn Connector
}

// Colour is the (typical) number of the colour module in the Jambel device.
type Colour int

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

// Reset resets Jambel to all lights off
func (jmb *Jambel) Reset() error {
	return jmb.conn.Send([]byte("reset\n"))
}

// On switches colour module on.
func (jmb *Jambel) On(colour Colour) error {
	cmd := fmt.Sprintf("set=%d,on\n", colour)
	return jmb.conn.Send([]byte(cmd))
}

// Off switches colour module off.
func (jmb *Jambel) Off(colour Colour) error {
	cmd := fmt.Sprintf("set=%d,off\n", colour)
	return jmb.conn.Send([]byte(cmd))
}

// Blink makes colour module blink.
func (jmb *Jambel) Blink(colour Colour) error {
	cmd := fmt.Sprintf("set=%d,blink\n", colour)
	return jmb.conn.Send([]byte(cmd))
}

// BlinkInverse makes colour module blink inversely.
func (jmb *Jambel) BlinkInverse(colour Colour) error {
	cmd := fmt.Sprintf("set=%d,blink_inverse\n", colour)
	return jmb.conn.Send([]byte(cmd))
}

// Flash makes colour module flash.
func (jmb *Jambel) Flash(colour Colour) error {
	cmd := fmt.Sprintf("set=%d,flash\n", colour)
	return jmb.conn.Send([]byte(cmd))
}

// SetAll sets all three lights (Green, Yellow, Red) to the given states.
func (jmb *Jambel) SetAll(green, yellow, red LightState) error {
	cmd := fmt.Sprintf("set_all=%d,%d,%d,0\n", red, yellow, green)
	return jmb.conn.Send([]byte(cmd))
}

// Close closes the connection to the Jambel device.
func (jmb *Jambel) Close() {
	jmb.conn.Close()
}
