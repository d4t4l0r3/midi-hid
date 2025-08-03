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

	midiInput := must(NewMidiInput("DJControl Inpulse 500 MIDI 1"))

	time.Sleep(time.Second * 20)
	midiInput.Stop()
}
