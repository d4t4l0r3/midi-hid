package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/charmbracelet/log"
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

	var (
		configPath string
		printDebugMsgs bool
	)

	flag.StringVar(&configPath, "f", "$HOME/.config/midi-hid/config.yaml", "Config file")
	flag.BoolVar(&printDebugMsgs, "debug", false, "Print debug messages")
	flag.Parse()

	if printDebugMsgs {
		log.SetLevel(log.DebugLevel)
	}

	log.Info("Starting...")
	config := must(ParseConfig(configPath))
	controllerList := must(config.Construct())
	defer controllerList.Stop()

	// wait for SIGINT
	sigintChan := make(chan os.Signal, 1)
	signal.Notify(sigintChan, os.Interrupt)
	<-sigintChan
}
