package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/bendahl/uinput"
)

type Config struct {
	Controller []ControllerConfig `yaml:"controller"`
}

type ControllerConfig struct {
	PortName string `yaml:"portName"`
	VendorID uint16 `yaml:"vendorID"`
	ProductID uint16 `yaml:"productID"`
	Mappings []MappingConfig `yaml:"mappings"`
}

type MappingConfig struct {
	Comment string `yaml:"comment"`
	Type MappingType `yaml:"type"`
	MidiChannel uint8 `yaml:"midiChannel"`
	MidiKey uint8 `yaml:"midiKey"`
	MidiController uint8 `yaml:"midiController"`
	Button ButtonName `yaml:"button"`
	Axis AxisName `yaml:"axis"`
	IsSigned bool `yaml:"isSigned"`
}

type MappingType string
type ButtonName string
type AxisName string

const (
	ButtonMappingType MappingType = "button"
	ControlMappingType MappingType = "control"
	ButtonNorth ButtonName = "north"
	ButtonEast ButtonName = "east"
	ButtonSouth ButtonName = "south"
	ButtonWest ButtonName = "west"
	AxisLeftX AxisName = "left-x"
	AxisLeftY AxisName = "left-y"
	AxisRightX AxisName = "right-x"
	AxisRightY AxisName = "right-y"
)

func ParseConfig(path string) (Config, error) {
	var config Config

	buffer, err := os.ReadFile(os.ExpandEnv(path))
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(buffer, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (config Config) Construct() (ControllerList, error) {
	var controllerList ControllerList

	for _, controllerConfig := range config.Controller {
		actualController, err := controllerConfig.Construct()
		if err != nil {
			return nil, err
		}

		controllerList = append(controllerList, actualController)
	}

	return controllerList, nil
}

func (cc ControllerConfig) Construct() (*Controller, error) {
	actualController, err := NewController(cc.PortName, cc.VendorID, cc.ProductID)
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

func (mc MappingConfig) Construct() (Mapping, error) {
	switch mc.Type {
	case ButtonMappingType:
		button, err := mc.Button.Construct()
		if err != nil {
			return ButtonMapping{}, err
		}

		return ButtonMapping{mc.Comment, mc.MidiChannel, mc.MidiKey, button}, nil
	case ControlMappingType:
		axis, err := mc.Axis.Construct()
		if err != nil {
			return ControlMapping{}, err
		}

		return ControlMapping{mc.Comment, mc.MidiChannel, mc.MidiController, axis, mc.IsSigned}, nil
	default:
		return ButtonMapping{}, fmt.Errorf("Invalid mapping type")
	}
}

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
	default:
		return -1, fmt.Errorf("Invalid button name \"%s\"", bn)
	}
}

func (an AxisName) Construct() (ControllerAxis, error) {
	switch an {
	case AxisLeftX:
		return LeftX, nil
	case AxisLeftY:
		return LeftY, nil
	case AxisRightX:
		return RightX, nil
	case AxisRightY:
		return RightY, nil
	default:
		return -1, fmt.Errorf("Invalid axis name \"%s\"", an)
	}
}
