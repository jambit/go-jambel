package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jambit/go-jambel"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s ADDR", os.Args[0])
		os.Exit(1)
	}
	addr := os.Args[1]
	jmb := jambel.NewNetworkJambel(addr)

	// traffic light phases
	var commands = []func() error{
		func() error { return jmb.On(jambel.GREEN) },
		func() error { return jmb.On(jambel.YELLOW) },
		func() error { return jmb.Blink(jambel.YELLOW) },
		func() error { return jmb.On(jambel.RED) },
		func() error { return jmb.SetAll(jambel.OFF, jambel.ON, jambel.ON) },
	}

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
