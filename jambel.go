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
	// Send sends a single command to the Jambel and reads the returned
	// bytes. Make sure that it is terminated with "\n"
	Send(cmd []byte) (response []byte, err error)

	// Close closes the connection with the Jambel and cleans up afterward.
	Close()
}

// Reset resets Jambel to all lights off
func (jmb *Jambel) Reset() error {
	_, err := jmb.conn.Send([]byte("reset\n"))
	return err
}

// On switches colour module on.
func (jmb *Jambel) On(colour Colour) error {
	cmd := fmt.Sprintf("set=%d,on\n", colour)
	_, err := jmb.conn.Send([]byte(cmd))
	return err
}

// Off switches colour module off.
func (jmb *Jambel) Off(colour Colour) error {
	cmd := fmt.Sprintf("set=%d,off\n", colour)
	_, err := jmb.conn.Send([]byte(cmd))
	return err
}

// Blink makes colour module blink.
func (jmb *Jambel) Blink(colour Colour) error {
	cmd := fmt.Sprintf("set=%d,blink\n", colour)
	_, err := jmb.conn.Send([]byte(cmd))
	return err
}

// BlinkInverse makes colour module blink inversely.
func (jmb *Jambel) BlinkInverse(colour Colour) error {
	cmd := fmt.Sprintf("set=%d,blink_inverse\n", colour)
	_, err := jmb.conn.Send([]byte(cmd))
	return err
}

// Flash makes colour module flash.
func (jmb *Jambel) Flash(colour Colour) error {
	cmd := fmt.Sprintf("set=%d,flash\n", colour)
	_, err := jmb.conn.Send([]byte(cmd))
	return err
}

// SetAll sets all three lights (Green, Yellow, Red) to the given states.
func (jmb *Jambel) SetAll(green, yellow, red LightState) error {
	cmd := fmt.Sprintf("set_all=%d,%d,%d,0\n", red, yellow, green)
	_, err := jmb.conn.Send([]byte(cmd))
	return err
}

// Status returns the current status of all light modules.
func (jmb *Jambel) Status() ([]byte, error) {
	return jmb.conn.Send([]byte("status\n"))
}

// Version returns the device version.
func (jmb *Jambel) Version() ([]byte, error) {
	return jmb.conn.Send([]byte("version\n"))
}

// Close closes the connection to the Jambel device.
func (jmb *Jambel) Close() {
	jmb.conn.Close()
}
