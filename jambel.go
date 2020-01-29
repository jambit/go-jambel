package jambel

import (
	"fmt"
	"github.com/google/gousb"
	"log"
)

const (
	// Bus 020 Device 018: ID 0403:6001 Future Technology Devices International Limited usb serial converter  Serial: ftDIMF19
	vendor  = 0x0403
	product = 0x6001

	// Available endpoints: [0x02(2,OUT) 0x81(1,IN)]
	epIn, epOut = 1, 2

	// colour modules
	GREEN  = 3
	YELLOW = 2
	RED    = 1

	// light module status
	OFF           = 0
	ON            = 1
	BLINK         = 2
	FLASH         = 3
	BLINK_INVERSE = 4
)

type Jambel interface {

	// Reset resets Jambel to all lights off
	Reset()

	// On switches [colour] on where colour is one of GREEN, RED or YELLOW
	On(colour int)
	Off(colour int)
	Blink(colour int)
	BlinkInverse(colour int)
	Flash(colour int)
	SetAll(green, yellow, red int)
}

type SerialJambel struct {
	// Fields for interacting with the USB connection
	context  *gousb.Context
	device   *gousb.Device
	intf     *gousb.Interface
	endpoint *gousb.OutEndpoint
	// interface cleanup
	_releaseInterface func()
}

// Close releases claimed interface, config, context, all associated
// resources and closes the device.
func (jmb *SerialJambel) Close() {
	jmb._releaseInterface()
	jmb.device.Close()
	jmb.context.Close()
}

// send sends command to Jambel.
// Don't forget to terminate commands with "\n"!
func (jmb *SerialJambel) send(cmd []byte) error {
	numBytes, err := jmb.endpoint.Write(cmd)
	if err != nil {
		log.Printf("only %d bytes written, returned error is %v", numBytes, err)
	} else {
		log.Printf("%d bytes successfully written", numBytes)
	}
	return err
}

// Reset resets Jambel to all lights off
func (jmb *SerialJambel) Reset() {
	_ = jmb.send([]byte("reset\n"))
}

// On switches [colour] on where colour is one of GREEN, RED or YELLOW
func (jmb *SerialJambel) On(colour int) {
	cmd := fmt.Sprintf("set=%d,on\n", colour)
	_ = jmb.send([]byte(cmd))
}

func (jmb *SerialJambel) Off(colour int) {
	cmd := fmt.Sprintf("set=%d,on\n", colour)
	_ = jmb.send([]byte(cmd))
}

func (jmb *SerialJambel) Blink(colour int) {
	cmd := fmt.Sprintf("set=%d,blink\n", colour)
	_ = jmb.send([]byte(cmd))
}

func (jmb *SerialJambel) BlinkInverse(colour int) {
	cmd := fmt.Sprintf("set=%d,blink_inverse\n", colour)
	_ = jmb.send([]byte(cmd))
}

func (jmb *SerialJambel) Flash(colour int) {
	cmd := fmt.Sprintf("set=%d,flash\n", colour)
	_ = jmb.send([]byte(cmd))
}

func (jmb *SerialJambel) SetAll(green, yellow, red int) {
	cmd := fmt.Sprintf("set_all=%d,%d,%d,0\n", red, yellow, green)
	_ = jmb.send([]byte(cmd))
}

func NewSerialJambel() (*SerialJambel, error) {

	// Initialize a new Context.
	ctx := gousb.NewContext()

	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(vendor, product)
	if err != nil {
		return nil, err
	}

	// Detach the device from whichever process already has it.
	_ = dev.SetAutoDetach(true)

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

	jmb := SerialJambel{
		context:           ctx,
		device:            dev,
		intf:              intf,
		endpoint:          ep,
		_releaseInterface: done,
	}
	return &jmb, nil
}
