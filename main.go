package main

import (
	"log"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"github.com/bendahl/uinput"
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
	controller := must(NewController("DJControl Inpulse 500 MIDI 1", 0x45e, 0x285)) // mimics xbox 360 controller
	defer controller.Stop()

	controller.AddMapping(ButtonMapping{"Play left", 1, 7, uinput.ButtonSouth})
	controller.AddMapping(ControlMapping{"Volume left", 1, 0, LeftY, false})

	time.Sleep(time.Second * 20)
}
