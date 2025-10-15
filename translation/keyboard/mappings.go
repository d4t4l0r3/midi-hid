package keyboard

import (
	"fmt"
	"math"

	"github.com/charmbracelet/log"
	"gitlab.com/gomidi/midi/v2"
	"github.com/bendahl/uinput"
)

// A ButtonMapping maps a MIDI Note to a gamepad button.
type ButtonMapping struct {
	CommentStr string
	MidiChannel uint8
	MidiKey uint8
	KeyboardKey int
}

// Is checks if the MIDI message msg triggers this Mapping, without actually triggering it.
func (m ButtonMapping) Is(msg midi.Message) bool {
	var channel, key uint8

	switch {
	case msg.GetNoteOn(&channel, &key, nil), msg.GetNoteOff(&channel, &key, nil):
		return (m.MidiChannel == channel && m.MidiKey == key)
	default:
		return false
	}
}

// TriggerIfMatch checks if the MIDI message msg triggers this Mapping, and if so,
// sends the corresponding input to virtKeyboard.
func (m ButtonMapping) TriggerIfMatch(msg midi.Message, virtKeyboard uinput.Keyboard) error {
	if m.Is(msg) {
		var velocity uint8
		msg.GetNoteOn(nil, nil, &velocity)
		switch msg.Type() {
		case midi.NoteOnMsg:
			if velocity != 0 {
				log.Debug(m.CommentStr, "status", "down")
				return virtKeyboard.KeyDown(m.KeyboardKey)
			}
			fallthrough // if reached here, velocity is 0 -> NoteOff
		case midi.NoteOffMsg:
			log.Debug(m.CommentStr, "status", "up")
			return virtKeyboard.KeyUp(m.KeyboardKey)
		default:
			return fmt.Errorf("Invalid message type triggered ButtonMapping")
		}
	}

	return nil
}

// Comment returns the Mappings comment.
func (m ButtonMapping) Comment() string {
	return m.CommentStr
}

// An EncoderMapping maps a MIDI Controller to two buttons.
type EncoderMapping struct {
	CommentStr string
	MidiChannel uint8
	MidiController uint8
	KeyboardKeyPositive int
	KeyboardKeyNegative int
}

// Is checks if the MIDI message msg triggers this Mapping, without actually triggering it.
func (m EncoderMapping) Is(msg midi.Message) bool {
	var channel, controller uint8

	if msg.GetControlChange(&channel, &controller, nil) {
		return (m.MidiChannel == channel && m.MidiController == controller)
	} else {
		return false
	}
}

// TriggerIfMatch checks if the MIDI message msg triggers this Mapping, and if so,
// sends the corresponding input to virtKeyboard.
func (m EncoderMapping) TriggerIfMatch(msg midi.Message, virtKeyboard uinput.Keyboard) error {
	if m.Is(msg) {
		var valueAbsolute uint8

		msg.GetControlChange(nil, nil, &valueAbsolute)
		
		switch valueAbsolute {
		case 1:
			log.Debug(m.CommentStr, "status", "increased")
			return virtKeyboard.ButtonPress(m.KeyboardKeyPositive)
		case 127:
			log.Debug(m.CommentStr, "status", "decreased")
			return virtKeyboard.ButtonPress(m.KeyboardKeyNegative)
		default:
			return fmt.Errorf("Invalid message type triggered EncoderMapping")
		}
	}

	return nil
}

// Comment returns the Mappings comment.
func (m EncoderMapping) Comment() string {
	return m.CommentStr
}
