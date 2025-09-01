package translation

import (
	"fmt"
	"math"

	"github.com/charmbracelet/log"
	"gitlab.com/gomidi/midi/v2"
	"github.com/bendahl/uinput"
)

// A Mapping is an interface for all types of Mappings.
type Mapping interface {
	Is(midi.Message) bool
	TriggerIfMatch(midi.Message, uinput.Gamepad) error
	Comment() string
}

// A ButtonMapping maps a MIDI Note to a gamepad button.
type ButtonMapping struct {
	CommentStr string
	MidiChannel uint8
	MidiKey uint8
	GamepadKey int
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
// sends the corresponding input to virtGamepad.
func (m ButtonMapping) TriggerIfMatch(msg midi.Message, virtGamepad uinput.Gamepad) error {
	if m.Is(msg) {
		var velocity uint8
		msg.GetNoteOn(nil, nil, &velocity)
		switch msg.Type() {
		case midi.NoteOnMsg:
			if velocity != 0 {
				log.Debug(m.CommentStr, "status", "down")
				return virtGamepad.ButtonDown(m.GamepadKey)
			}
			fallthrough // if reached here, velocity is 0 -> NoteOff
		case midi.NoteOffMsg:
			log.Debug(m.CommentStr, "status", "up")
			return virtGamepad.ButtonUp(m.GamepadKey)
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
	GamepadKeyPositive int
	GamepadKeyNegative int
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
// sends the corresponding input to virtGamepad.
func (m EncoderMapping) TriggerIfMatch(msg midi.Message, virtGamepad uinput.Gamepad) error {
	if m.Is(msg) {
		var valueAbsolute uint8

		msg.GetControlChange(nil, nil, &valueAbsolute)
		
		switch valueAbsolute {
		case 1:
			log.Debug(m.CommentStr, "status", "increased")
			return virtGamepad.ButtonPress(m.GamepadKeyPositive)
		case 127:
			log.Debug(m.CommentStr, "status", "decreased")
			return virtGamepad.ButtonPress(m.GamepadKeyNegative)
		default:
			return fmt.Errorf("Invalid message type triggered ButtonMapping")
		}
	}

	return nil
}

// Comment returns the Mappings comment.
func (m EncoderMapping) Comment() string {
	return m.CommentStr
}
type ControllerAxis int

const (
	LeftX ControllerAxis = iota
	LeftY
	RightX
	RightY
)

type ControlMapping struct {
	CommentStr string
	MidiChannel uint8
	MidiController uint8
	Axis ControllerAxis
	IsSigned bool
	Deadzone float64
}

// Is checks if the MIDI message msg triggers this Mapping, without actually triggering it.
func (m ControlMapping) Is(msg midi.Message) bool {
	var channel, controller uint8

	if msg.GetControlChange(&channel, &controller, nil) {
		return (m.MidiChannel == channel && m.MidiController == controller)
	} else {
		return false
	}
}

// TriggerIfMatch checks if the MIDI message msg triggers this Mapping, and if so,
// sends the corresponding input to virtGamepad.
func (m ControlMapping) TriggerIfMatch(msg midi.Message, virtGamepad uinput.Gamepad) error {
	if m.Is(msg) {
		var (
			valueAbsolute uint8
			valueNormalised float64
		)

		msg.GetControlChange(nil, nil, &valueAbsolute)

		// value is 0-127, normalise
		valueNormalised = float64(valueAbsolute) / 127
		if m.IsSigned {
			valueNormalised *= 2
			valueNormalised -= 1
		}

		if math.Abs(valueNormalised) < m.Deadzone {
			valueNormalised = 0
		}

		log.Debug(m.CommentStr, "value", valueNormalised, "deadzone", m.Deadzone)

		switch m.Axis {
		case LeftX:
			return virtGamepad.LeftStickMoveX(float32(valueNormalised))
		case LeftY:
			return virtGamepad.LeftStickMoveY(float32(valueNormalised))
		case RightX:
			return virtGamepad.RightStickMoveX(float32(valueNormalised))
		case RightY:
			return virtGamepad.RightStickMoveY(float32(valueNormalised))
		}
	}

	return nil
}

// Comment returns the Mappings comment.
func (m ControlMapping) Comment() string {
	return m.CommentStr
}
