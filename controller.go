package main

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
	mappings []Mapping
	abortChan chan interface{}
	virtGamepad uinput.Gamepad
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

	abortChan := make(chan interface{})

	controller := &Controller{midiInput, nil, abortChan, virtGamepad}

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

// AddMapping adds a mapping to the Controller.
func (c *Controller) AddMapping(mapping Mapping) {
	c.mappings = append(c.mappings, mapping)
}

// Stop quits the update loop and terminates all corresponding connections.
func (c Controller) Stop() {
	c.midiInput.Stop()
	c.abortChan <- struct{}{}
	c.virtGamepad.Close()
}

func (c Controller) update(msg midi.Message) {
	for _, mapping := range c.mappings {
		err := mapping.TriggerIfMatch(msg, c.virtGamepad)
		if err != nil {
			log.Errorf("Error in Mapping \"%s\": %v", mapping.Comment(), err)
		}
	}
}
