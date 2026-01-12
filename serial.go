package jambel

import (
	"bytes"
	"fmt"

	"github.com/karalabe/usb"
)

const (

	// Bus 020 Device 018: ID 0403:6001
	// Future Technology Devices International Limited usb serial converter
	// Serial: ftDIMF19
	vendor  = 0x0403
	product = 0x6001
)

type SerialConnection struct {
	// Fields for interacting with the USB connection
	device usb.Device
}

// Send sends command to Jambel.
// Don't forget to terminate commands with "\n"!
func (jmb *SerialConnection) Send(cmd []byte) ([]byte, error) {
	_, err := jmb.device.Write(cmd)
	fmt.Printf(">>> %s", cmd)
	out, _ := readFromDevice(jmb.device, []byte("\n"))
	fmt.Printf("<<< %s", out)
	return out, err
}

func readFromDevice(device usb.Device, readUntil []byte) (out []byte, err error) {
	recv := make([]byte, 1)
	var n int

	for {
		n, err = device.Read(recv)
		if err != nil || n <= 0 {
			break
		}
		out = append(out, recv...)
		if bytes.Contains(out, readUntil) {
			break
		}
	}
	return out, err
}

// Close releases claimed interface, config, context, all associated
// resources and closes the device.
func (jmb *SerialConnection) Close() {
	_ = jmb.device.Close()
}

// NewSerialJambel creates an instance of a USB Jambel.
func NewSerialJambel() (*Jambel, error) {

	// Open any device with a given VID/PID using a convenience function.
	devs, err := usb.Enumerate(vendor, product)
	if err != nil {
		return nil, err
	}

	// FIXME: Handle multiple USB devices
	device, err := devs[0].Open()
	if err != nil {
		return nil, err
	}

	conn := &SerialConnection{device: device}
	return &Jambel{conn: conn}, nil
}
