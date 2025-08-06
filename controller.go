package main

import (
	"github.com/charmbracelet/log"
	"gitlab.com/gomidi/midi/v2"
	"github.com/bendahl/uinput"
)

type ControllerList []*Controller

func (cl ControllerList) Stop() {
	for _, controller := range cl {
		controller.Stop()
	}
}

type Controller struct {
	midiInput *MidiInput
	mappings []Mapping
	abortChan chan interface{}
	virtGamepad uinput.Gamepad
}

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

func (c *Controller) AddMapping(mapping Mapping) {
	c.mappings = append(c.mappings, mapping)
}

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
