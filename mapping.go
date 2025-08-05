package main

import (
	"fmt"
	"log"

	"gitlab.com/gomidi/midi/v2"
	"github.com/bendahl/uinput"
)

type Mapping interface {
	Is(midi.Message) bool
	TriggerIfMatch(midi.Message, uinput.Gamepad) error
	Comment() string
}

type ButtonMapping struct {
	comment string
	midiChannel uint8
	midiKey uint8
	gamepadKey int
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

func (m ButtonMapping) TriggerIfMatch(msg midi.Message, virtGamepad uinput.Gamepad) error {
	if m.Is(msg) {
		switch msg.Type() {
		case midi.NoteOnMsg:
			log.Printf("%s: Button down\n", m.comment)
			return virtGamepad.ButtonDown(m.gamepadKey)
		case midi.NoteOffMsg:
			log.Printf("%s: Button up\n", m.comment)
			return virtGamepad.ButtonUp(m.gamepadKey)
		default:
			return fmt.Errorf("Invalid message type triggered ButtonMapping")
		}
	}

	return nil
}

func (m ButtonMapping) Comment() string {
	return m.comment
}

type ControllerAxis int

const (
	LeftX ControllerAxis = iota
	LeftY
	RightX
	RightY
)

type ControlMapping struct {
	comment string
	midiChannel uint8
	midiController uint8
	axis ControllerAxis
	isSigned bool
}

func (m ControlMapping) Is(msg midi.Message) bool {
	var channel, controller uint8

	if msg.GetControlChange(&channel, &controller, nil) {
		return (m.midiChannel == channel && m.midiController == controller)
	} else {
		return false
	}
}

func (m ControlMapping) TriggerIfMatch(msg midi.Message, virtGamepad uinput.Gamepad) error {
	if m.Is(msg) {
		var (
			valueAbsolute uint8
			valueNormalised float32
		)

		msg.GetControlChange(nil, nil, &valueAbsolute)

		// value is 0-127, normalise
		valueNormalised = float32(valueAbsolute) / 127
		if m.isSigned {
			valueNormalised *= 2
			valueNormalised -= 1
		}

		log.Printf("%s: value %v\n", m.comment, valueNormalised)

		switch m.axis {
		case LeftX:
			return virtGamepad.LeftStickMoveX(valueNormalised)
		case LeftY:
			return virtGamepad.LeftStickMoveY(valueNormalised)
		case RightX:
			return virtGamepad.RightStickMoveX(valueNormalised)
		case RightY:
			return virtGamepad.RightStickMoveY(valueNormalised)
		}
	}

	return nil
}

func (m ControlMapping) Comment() string {
	return m.comment
}
