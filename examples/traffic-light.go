package main

import (
	"log"
	"time"
)

func main() {

	jmb, err := NewSerialJambel()
	defer jmb.Close()
	if err != nil {
		log.Fatalf("could not open device: %v", err)
	}

	// traffic light phases
	var commands = []func(){
		func() { jmb.On(GREEN) },
		func() { jmb.On(YELLOW) },
		func() { jmb.Blink(YELLOW) },
		func() { jmb.On(RED) },
		func() { jmb.SetAll(OFF, ON, ON) },
	}

	var i = 0

	for {
		jmb.Reset()
		commands[i]()
		time.Sleep(2 * time.Second)

		i++
		if i >= len(commands) {
			i = 0
		}
	}

}
