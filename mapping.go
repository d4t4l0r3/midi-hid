package main

import (
	"gitlab.com/gomidi/midi/v2"
)

type Mapping interface {
	Is(midi.Message) bool
}

type ButtonMapping struct {
	midiChannel uint8
	midiKey uint8
}

func (m ButtonMapping) Is(msg midi.Message) bool {
	var channel, key uint8

	switch {
	case msg.GetNoteOn(&channel, &key, nil), msg.GetNoteOff(&channel, &key, nil):
		return (m.midiChannel == channel && m.midiKey == key)
	default:
		return false
	}
}

type ControlMapping struct {
	midiChannel uint8
	midiController uint8
}

func (m ControlMapping) Is(msg midi.Message) bool {
	var channel, controller uint8

	if msg.GetControlChange(&channel, &controller, nil) {
		return (m.midiChannel == channel && m.midiController == controller)
	} else {
		return false
	}
}
