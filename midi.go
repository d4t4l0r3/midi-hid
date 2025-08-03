package main

import (
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type MidiInput struct {
	input drivers.In
	Messages chan midi.Message
	Stop func()
}

func NewMidiInput(portName string) (*MidiInput, error) {
	input, err := midi.FindInPort(portName)
	if err != nil {
		return nil, err
	}

	messages := make(chan midi.Message)

	midiListener := func(msg midi.Message, timestampMs int32) {
		if msg.IsOneOf(midi.NoteOnMsg, midi.NoteOffMsg, midi.ControlChangeMsg) {
			messages <- msg
		}
	}

	stopFunc, err := midi.ListenTo(input, midiListener, midi.UseSysEx())
	if err != nil {
		return nil, err
	}

	return &MidiInput{input, messages, stopFunc}, nil
}
