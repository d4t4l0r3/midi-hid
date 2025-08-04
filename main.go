package main

import (
	"log"
	"time"

	"gitlab.com/gomidi/midi/v2"
)

func must[T any](obj T, err error) T {
	if err != nil {
		log.Fatal(err)
	}

	return obj
}

func main() {
	defer midi.CloseDriver()

	log.Println("Starting...")
	config := must(ParseConfig("config.yaml"))
	controllerList := must(config.Construct())
	defer controllerList.Stop()

	time.Sleep(time.Second * 20)
}
