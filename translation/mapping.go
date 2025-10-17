package translation

import (
	"gitlab.com/gomidi/midi/v2"
	"github.com/bendahl/uinput"
)

// A GamepadMapping is an interface for all types of gamepad mappings.
type GamepadMapping interface {
	Is(midi.Message) bool
	TriggerIfMatch(midi.Message, uinput.Gamepad) error
	Comment() string
}

// A KeyboardMapping is an interface for all types of keyboard mappings.
type KeyboardMapping interface {
	Is(midi.Message) bool
	TriggerIfMatch(midi.Message, uinput.Keyboard) error
	Comment() string
}
