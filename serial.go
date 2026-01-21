package jambel

import (
	"fmt"

	"github.com/google/gousb"
)

const (

	// Bus 020 Device 018: ID 0403:6001
	// Future Technology Devices International Limited usb serial converter
	// Serial: ftDIMF19
	vendor  = 0x0403
	product = 0x6001

	// Available endpoints: [0x02(2,OUT) 0x81(1,IN)]
	epIn, epOut = 1, 2
)

type SerialConnection struct {
	// Fields for interacting with the USB connection
	context *gousb.Context
	device  *gousb.Device
	intf    *gousb.Interface

	read  *gousb.InEndpoint
	write *gousb.OutEndpoint

	// deferred interface cleanup
	_releaseInterface func()
}

// Send sends command to Jambel.
// Don't forget to terminate commands with "\n"!
func (jmb *SerialConnection) Send(cmd []byte) ([]byte, error) {
	_, err := jmb.write.Write(cmd)
	if err != nil {
		return nil, err
	}
	return readUntil(jmb.read, []byte("\n"))
}

// Close releases claimed interface, config, context, all associated
// resources and closes the device.
func (jmb *SerialConnection) Close() {
	jmb._releaseInterface()
	_ = jmb.device.Close()
	_ = jmb.context.Close()
}

func NewSerialJambel() (*Jambel, error) {

	// Initialize a new Context.
	ctx := gousb.NewContext()

	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(vendor, product)
	if err != nil {
		return nil, err
	}

	// Detach the device from whichever process already has it.
	if err = dev.SetAutoDetach(true); err != nil {
		return nil, fmt.Errorf("could not auto-detach device: %w", err)
	}

	// Claim the default interface using a convenience function. The default interface is always #0 alt #0
	// in the currently active config.
	intf, done, err := dev.DefaultInterface()
	if err != nil {
		return nil, fmt.Errorf("could not claim default interface: %w", err)
	}

	readEndpoint, err := intf.InEndpoint(epIn)
	if err != nil {
		return nil, err
	}
	writeEndpoint, err := intf.OutEndpoint(epOut)
	if err != nil {
		return nil, err
	}

	conn := &SerialConnection{
		context:           ctx,
		device:            dev,
		intf:              intf,
		read:              readEndpoint,
		write:             writeEndpoint,
		_releaseInterface: done,
	}
	return &Jambel{conn: conn}, nil
}
