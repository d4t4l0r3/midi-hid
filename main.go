package main

import (
	"log"
	"time"

	"gitlab.com/gomidi/midi/v2"
)

func must[T any](obj T, err error) T {
	if err != nil {
		log.Fatal(err)
	}

	return obj
}

func main() {
	defer midi.CloseDriver()

	log.Println("Starting...")
	controller := must(NewController("DJControl Inpulse 500 MIDI 1"))
	controller.AddMapping(ButtonMapping{1, 7}) // play left
	controller.AddMapping(ControlMapping{1, 0}) // volume left

	time.Sleep(time.Second * 20)
	log.Println("Stopping...")
	controller.Stop()
	log.Println("Stopped.")
}
