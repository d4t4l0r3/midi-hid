package gamepad

import (
	"fmt"

	"github.com/d4t4l0r3/midi-hid/translation"
	"github.com/d4t4l0r3/midi-hid/translation/gamepad"

	"github.com/bendahl/uinput"
	"github.com/charmbracelet/log"
)

// A ControllerConfig represents the data needed to later construct a Controller object.
type ControllerConfig struct {
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
	Button ButtonName `yaml:"button"`
	ButtonNegative ButtonName `yaml:"buttonNegative"`
	Axis AxisName `yaml:"axis"`
	IsSigned bool `yaml:"isSigned"`
	Deadzone float64 `yaml:"deadzone"`
}

type MappingType string
type ButtonName string
type AxisName string

const (
	ButtonMappingType MappingType = "button"
	AxisMappingType MappingType = "axis"
	EncoderMappingType MappingType = "encoder"
	ButtonNorth ButtonName = "north"
	ButtonEast ButtonName = "east"
	ButtonSouth ButtonName = "south"
	ButtonWest ButtonName = "west"
	ButtonL1 ButtonName = "l1"
	ButtonL2 ButtonName = "l2"
	ButtonL3 ButtonName = "l3"
	ButtonR1 ButtonName = "r1"
	ButtonR2 ButtonName = "r2"
	ButtonR3 ButtonName = "r3"
	ButtonSelect ButtonName = "select"
	ButtonStart ButtonName = "start"
	ButtonDpadUp ButtonName = "dpad-up"
	ButtonDpadDown ButtonName = "dpad-down"
	ButtonDpadLeft ButtonName = "dpad-left"
	ButtonDpadRight ButtonName = "dpad-right"
	AxisLeftX AxisName = "left-x"
	AxisLeftY AxisName = "left-y"
	AxisRightX AxisName = "right-x"
	AxisRightY AxisName = "right-y"
)

// Construct builds a Controller object and its corresponding mappings.
// Aborts and returns an error if the midi port was not found or one of
// the Mappings is invalid.
func (cc ControllerConfig) Construct() (*translation.Controller, error) {
	actualController, err := translation.NewController(cc.PortName, cc.VendorID, cc.ProductID)
	if err != nil {
		return actualController, err
	}

	for _, mappingConfig := range cc.Mappings {
		actualMapping, err := mappingConfig.Construct()
		if err != nil {
			return nil, err
		}

		actualController.AddMapping(actualMapping)
	}

	return actualController, nil
}

// Construct builds the Mapping object. Returns an error if config is invalid.
func (mc MappingConfig) Construct() (translation.Mapping, error) {
	switch mc.Type {
	case ButtonMappingType:
		button, err := mc.Button.Construct()
		if err != nil {
			return gamepad.ButtonMapping{}, err
		}

		log.Debug("Parsed button mapping", "comment", mc.Comment, "midiChannel", mc.MidiChannel, "midiKey", mc.MidiKey, "button", button)

		return gamepad.ButtonMapping{mc.Comment, mc.MidiChannel, mc.MidiKey, button}, nil
	case EncoderMappingType:
		button, err := mc.Button.Construct()
		if err != nil {
			return gamepad.EncoderMapping{}, err
		}

		buttonNegative, err := mc.ButtonNegative.Construct()
		if err != nil {
			return gamepad.EncoderMapping{}, err
		}

		log.Debug("Parsed encoder mapping", "comment", mc.Comment, "midiChannel", mc.MidiChannel, "midiController", mc.MidiController, "button", button, "buttonNegative", buttonNegative)

		return gamepad.EncoderMapping{mc.Comment, mc.MidiChannel, mc.MidiController, button, buttonNegative}, nil
	case AxisMappingType:
		axis, err := mc.Axis.Construct()
		if err != nil {
			return gamepad.AxisMapping{}, err
		}

		log.Debug("Parsed axis mapping", "comment", mc.Comment, "midiChannel", mc.MidiChannel, "midiController", mc.MidiController, "axis", axis, "isSigned", mc.IsSigned, "deadzone", mc.Deadzone)

		return gamepad.AxisMapping{mc.Comment, mc.MidiChannel, mc.MidiController, axis, mc.IsSigned, mc.Deadzone}, nil
	default:
		return gamepad.ButtonMapping{}, fmt.Errorf("Invalid mapping type")
	}
}

// Construct converts a ButtonName to its corresponding key code, or returns an error if the
// name is unknown.
func (bn ButtonName) Construct() (int, error) {
	switch bn {
	case ButtonNorth:
		return uinput.ButtonNorth, nil
	case ButtonEast:
		return uinput.ButtonEast, nil
	case ButtonSouth:
		return uinput.ButtonSouth, nil
	case ButtonWest:
		return uinput.ButtonWest, nil
	case ButtonL1:
		return uinput.ButtonBumperLeft, nil
	case ButtonL2:
		return uinput.ButtonTriggerLeft, nil
	case ButtonL3:
		return uinput.ButtonThumbLeft, nil
	case ButtonR1:
		return uinput.ButtonBumperRight, nil
	case ButtonR2:
		return uinput.ButtonTriggerRight, nil
	case ButtonR3:
		return uinput.ButtonThumbRight, nil
	case ButtonSelect:
		return uinput.ButtonSelect, nil
	case ButtonStart:
		return uinput.ButtonStart, nil
	case ButtonDpadUp:
		return uinput.ButtonDpadUp, nil
	case ButtonDpadDown:
		return uinput.ButtonDpadDown, nil
	case ButtonDpadLeft:
		return uinput.ButtonDpadLeft, nil
	case ButtonDpadRight:
		return uinput.ButtonDpadRight, nil
	default:
		return -1, fmt.Errorf("Invalid button name \"%s\"", bn)
	}
}

// Construct converts an AxisName into the internal representation for a ControllerAxis.
// Returns an error if AxisName is invalid.
func (an AxisName) Construct() (gamepad.ControllerAxis, error) {
	switch an {
	case AxisLeftX:
		return gamepad.LeftX, nil
	case AxisLeftY:
		return gamepad.LeftY, nil
	case AxisRightX:
		return gamepad.RightX, nil
	case AxisRightY:
		return gamepad.RightY, nil
	default:
		return -1, fmt.Errorf("Invalid axis name \"%s\"", an)
	}
}
