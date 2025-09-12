package config

import (
	"os"

	"github.com/d4t4l0r3/midi-hid/config/gamepad"
	"github.com/d4t4l0r3/midi-hid/translation"

	"github.com/goccy/go-yaml"
)

// Config is the root type of a config, consisting of an arbitrary number of controller configs.
type Config struct {
	Controller []gamepad.ControllerConfig `yaml:"controller"`
}

// ParseConfig takes the path to a config file and returns the parsed Config object,
// or an error if thrown.
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

// Construct iterates over all ControllerConfigs and constructs the Controller objects.
// In case of a failure, it aborts and returns an error.
func (config Config) Construct() (translation.ControllerList, error) {
	var controllerList translation.ControllerList

	for _, controllerConfig := range config.Controller {
		actualController, err := controllerConfig.Construct()
		if err != nil {
			return nil, err
		}

		controllerList = append(controllerList, actualController)
	}

	return controllerList, nil
}
