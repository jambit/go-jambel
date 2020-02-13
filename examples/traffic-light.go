package main

import (
	"time"

	"github.com/jambit/go-jambel"
)

func main() {

	jmb := jambel.NewNetworkJambel("ampel10.dev.jambit.com:10001")

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
