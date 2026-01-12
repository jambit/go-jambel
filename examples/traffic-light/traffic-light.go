package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jambit/go-jambel"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s ADDR\n", os.Args[0])
		fmt.Printf("Example: %s 192.168.1.100:10001\n", os.Args[0])
		os.Exit(1)
	}

	addr := os.Args[1]
	jmb, err := jambel.NewNetworkJambel(addr)
	if err != nil {
		fmt.Printf("Failed to connect to Jambel at %s: %v\n", addr, err)
		os.Exit(1)
	}
	defer jmb.Close()

	// traffic light phases to iterate through
	var commands = []func() error{
		func() error { return jmb.SetAll(jambel.On, jambel.Off, jambel.Off) },
		func() error { return jmb.SetAll(jambel.Off, jambel.On, jambel.Off) },
		func() error { return jmb.Blink(jambel.Yellow) },
		func() error { return jmb.SetAll(jambel.Off, jambel.Off, jambel.On) },
		func() error { return jmb.SetAll(jambel.Off, jambel.On, jambel.On) },
	}

	version, _ := jmb.Version()
	fmt.Printf("jambel version:\n%s\n", version)
	_ = jmb.Reset()

	var i = 0
	for {

		if err := commands[i](); err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)

		i++
		if i >= len(commands) {
			i = 0
		}
	}

}
