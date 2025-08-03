package main

import (
	"log"
	"gitlab.com/gomidi/midi/v2"
)

type Controller struct {
	midiInput *MidiInput
	mappings []Mapping
	abortChan chan interface{}
}

func NewController(portName string) (*Controller, error) {
	midiInput, err := NewMidiInput(portName)
	if err != nil {
		return nil, err
	}

	abortChan := make(chan interface{})

	controller := &Controller{midiInput, nil, abortChan}

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
}

func (c Controller) update(msg midi.Message) {
	for _, mapping := range c.mappings {
		if mapping.Is(msg) {
			log.Println("Mapping triggered!\n")
		}
	}
}
