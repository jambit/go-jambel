package jambel

import (
	"fmt"
	"log"

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
	context  *gousb.Context
	device   *gousb.Device
	intf     *gousb.Interface
	endpoint *gousb.OutEndpoint
	// deferred interface cleanup
	_releaseInterface func()
}

// send sends command to Jambel.
// Don't forget to terminate commands with "\n"!
func (jmb *SerialConnection) Send(cmd []byte) error {
	numBytes, err := jmb.endpoint.Write(cmd)
	if err != nil {
		log.Printf("only %d bytes written, returned error is %v", numBytes, err)
	} else {
		log.Printf("%d bytes successfully written", numBytes)
	}
	return err
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
		fmt.Errorf("could not auto-detach device: %v", err)
	}

	// Claim the default interface using a convenience function.
	// The default interface is always #0 alt #0 in the currently active
	// config.
	intf, done, err := dev.DefaultInterface()
	if err != nil {
		return nil, err
	}
	// Open an OUT endpoint.
	ep, err := intf.OutEndpoint(epOut)
	if err != nil {
		return nil, err
	}

	conn := &SerialConnection{
		context:           ctx,
		device:            dev,
		intf:              intf,
		endpoint:          ep,
		_releaseInterface: done,
	}
	return &Jambel{Connection: conn}, nil
}
