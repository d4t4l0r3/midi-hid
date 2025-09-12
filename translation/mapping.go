package translation

import (
	"gitlab.com/gomidi/midi/v2"
	"github.com/bendahl/uinput"
)

// A Mapping is an interface for all types of Mappings.
type Mapping interface {
	Is(midi.Message) bool
	TriggerIfMatch(midi.Message, uinput.Gamepad) error
	Comment() string
}
