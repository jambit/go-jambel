# go-jambel

A Go library for controlling jambit's Jambel fast feedback device - a USB/Network-connected traffic light with red, yellow, and green modules.

## Features

- Control individual light modules (Red, Yellow, Green)
- Multiple light states: On, Off, Blink, Flash, BlinkInverse
- Support for both USB serial and network (Telnet) connections
- Simple and intuitive API
- Built-in examples

## Installation

### Prerequisites

- Go 1.24.0 or higher
- For USB connections: libusb (see platform-specific instructions below)

### Install the library

```bash
go get github.com/jambit/go-jambel
```

### Platform-specific setup

#### macOS
```bash
brew install libusb
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get install libusb-1.0-0-dev
```

#### Linux (Fedora/RHEL)
```bash
sudo dnf install libusb-devel
```

## Quick Start

### Network Connection (Telnet)

```go
package main

import (
	"fmt"

	"github.com/jambit/go-jambel"
)

func main() {
	// Connect to Jambel via network
	jmb, err := jambel.NewNetworkJambel("192.168.1.100:10001")
	if err != nil {
		panic(err)
	}
	defer jmb.Close()

	// Turn on the green light
	if err := jmb.On(jambel.Green); err != nil {
		panic(err)
	}

	// Get device version
	if version, err := jmb.Version(); err == nil {
		fmt.Printf("Jambel version: %s\n", version)
	}
}
```

### USB Serial Connection

```go
package main

import (
    "github.com/jambit/go-jambel"
)

func main() {
    // Connect to Jambel via USB
    jmb, err := jambel.NewSerialJambel()
    if err != nil {
        panic(err)
    }
    defer jmb.Close()

    // Turn on the red light
    if err := jmb.On(jambel.Red); err != nil {
        panic(err)
    }
}
```

## Examples

### Traffic Light Simulation

Run the included traffic light example:

```bash
cd examples/traffic-light
go run traffic-light.go 192.168.1.100:10001
```

This example cycles through typical traffic light phases:
1. Green light on
2. Yellow light on
3. Yellow light blinking
4. Red light on
5. Red and Yellow lights on
6. Repeat

### Simple Status Check

```go
package main

import (
	"fmt"

	"github.com/jambit/go-jambel"
)

func main() {
	jmb, err := jambel.NewNetworkJambel("192.168.1.100:10001")
	if err != nil {
		panic(err)
	}
	defer jmb.Close()

	// Get version
	version, err := jmb.Version()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Version: %s\n", version)

	// Get status
	status, err := jmb.Status()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Status: %s\n", status)
}
```

### Build Status Indicator

```go
package main

import (
    "github.com/jambit/go-jambel"
    "time"
)

func main() {
    jmb, err := jambel.NewNetworkJambel("192.168.1.100:10001")
    if err != nil {
        panic(err)
    }
    defer jmb.Close()

    // Reset to clean state
    jmb.Reset()

    // Building - yellow blinking
    jmb.Blink(jambel.Yellow)
    time.Sleep(5 * time.Second)

    // Check your build status here...
    buildSuccess := true

    // Show result
    if buildSuccess {
        // Success - green on
        jmb.SetAll(jambel.On, jambel.Off, jambel.Off)
    } else {
        // Failed - red on
        jmb.SetAll(jambel.Off, jambel.Off, jambel.On)
    }
}
```
