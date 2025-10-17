package translation

import (
	"github.com/charmbracelet/log"
	"gitlab.com/gomidi/midi/v2"
	"github.com/bendahl/uinput"
)

// A ControllerList is a list of controllers. Duh.
type ControllerList []*Controller

// Stop iterates over all Controller objects and Stops their update loops and MIDI connections.
// Always call this for a clean shutdown. Meant to be deferred.
func (cl ControllerList) Stop() {
	for _, controller := range cl {
		controller.Stop()
	}
}

// A Controller object manages the translation from MIDI to uinput.
type Controller struct {
	midiInput *MidiInput
	gamepadMappings []GamepadMapping
	keyboardMappings []KeyboardMapping
	abortChan chan interface{}
	virtGamepad uinput.Gamepad
	virtKeyboard uinput.Keyboard
}

// NewController builds a new Controller object reading from the MIDI port specified by portName,
// and registers a virtual uinput-Gamepad using vendorID and productID.
func NewController(portName string, vendorID, productID uint16) (*Controller, error) {
	if vendorID == 0 && productID == 0 {
		// if no IDs were defined, imitate XBox 360 controller
		vendorID = 0x45e
		productID = 0x285
	}
	midiInput, err := NewMidiInput(portName)
	if err != nil {
		return nil, err
	}

	virtGamepad, err := uinput.CreateGamepad("/dev/uinput", []byte(portName), vendorID, productID)
	if err != nil {
		return nil, err
	}

	virtKeyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte(portName))
	if err != nil {
		return nil, err
	}

	abortChan := make(chan interface{})

	controller := &Controller{midiInput, nil, nil, abortChan, virtGamepad, virtKeyboard}

	go func() {
		for {
			select {
			case midiMessage := <-midiInput.Messages:
				controller.update(midiMessage)
			case <-abortChan:
				return
			}
		}
	}()

	return controller, nil
}

// AddGamepadMapping adds a mapping to the Controller.
func (c *Controller) AddGamepadMapping(mapping GamepadMapping) {
	c.gamepadMappings = append(c.gamepadMappings, mapping)
}

// AddKeyboardMapping adds a mapping to the Controller.
func (c *Controller) AddKeyboardMapping(mapping KeyboardMapping) {
	c.keyboardMappings = append(c.keyboardMappings, mapping)
}

// Stop quits the update loop and terminates all corresponding connections.
func (c Controller) Stop() {
	c.midiInput.Stop()
	c.abortChan <- struct{}{}
	c.virtGamepad.Close()
	c.virtKeyboard.Close()
}

func (c Controller) update(msg midi.Message) {
	for _, mapping := range c.gamepadMappings {
		err := mapping.TriggerIfMatch(msg, c.virtGamepad)
		if err != nil {
			log.Errorf("Error in Mapping \"%s\": %v", mapping.Comment(), err)
		}
	}
	for _, mapping := range c.keyboardMappings {
		err := mapping.TriggerIfMatch(msg, c.virtKeyboard)
		if err != nil {
			log.Errorf("Error in Mapping \"%s\": %v", mapping.Comment(), err)
		}
	}
}
