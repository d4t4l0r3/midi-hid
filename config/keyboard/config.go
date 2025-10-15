package keyboard

import (
	"fmt"

	"github.com/d4t4l0r3/midi-hid/translation"
	"github.com/d4t4l0r3/midi-hid/translation/keyboard"

	"github.com/charmbracelet/log"
)

// A KeyboardConfig represents the data needed to later construct a Keyboard object.
type KeyboardConfig struct {
	PortName string `yaml:"portName"`
	VendorID uint16 `yaml:"vendorID"`
	ProductID uint16 `yaml:"productID"`
	Mappings []MappingConfig `yaml:"mappings"`
}

// A MappingConfig consists of all data possibly needed to construct a mapping, both button and control.
type MappingConfig struct {
	Comment string `yaml:"comment"`
	Type MappingType `yaml:"type"`
	MidiChannel uint8 `yaml:"midiChannel"`
	MidiKey uint8 `yaml:"midiKey"`
	MidiController uint8 `yaml:"midiController"`
	KeyboardKey KeyName `yaml:"key"`
	KeyboardKeyNegative KeyName `yaml:"keyNegative"`
	IsSigned bool `yaml:"isSigned"`
	Deadzone float64 `yaml:"deadzone"`
}

type MappingType string
type KeyName string

const (
	KeyMappingType MappingType = "button"
	EncoderMappingType MappingType = "encoder"
)

// Construct builds the Mapping object. Returns an error if config is invalid.
func (mc MappingConfig) Construct() (translation.KeyboardMapping, error) {
	switch mc.Type {
	case KeyMappingType:
		key, err := mc.KeyboardKey.Construct()
		if err != nil {
			return keyboard.KeyMapping{}, err
		}

		log.Debug("Parsed key mapping", "comment", mc.Comment, "midiChannel", mc.MidiChannel, "midiKey", mc.MidiKey, "key", key)

		return keyboard.KeyMapping{mc.Comment, mc.MidiChannel, mc.MidiKey, key}, nil
	case EncoderMappingType:
		key, err := mc.KeyboardKey.Construct()
		if err != nil {
			return keyboard.EncoderMapping{}, err
		}

		keyNegative, err := mc.KeyboardKeyNegative.Construct()
		if err != nil {
			return keyboard.EncoderMapping{}, err
		}

		log.Debug("Parsed encoder mapping", "comment", mc.Comment, "midiChannel", mc.MidiChannel, "midiController", mc.MidiController, "key", key, "keyNegative", keyNegative)

		return keyboard.EncoderMapping{mc.Comment, mc.MidiChannel, mc.MidiController, key, keyNegative}, nil
	default:
		return keyboard.KeyMapping{}, fmt.Errorf("Invalid mapping type")
	}
}

// Construct converts a KeyName to its corresponding key code, or returns an error if the
// name is unknown.
func (kn KeyName) Construct() (int, error) {
	// TODO: implement
	switch kn {
	default:
		return -1, fmt.Errorf("Invalid button name \"%s\"", kn)
	}
}
