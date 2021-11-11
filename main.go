package main

import (
	"flag"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xantios/koekeloere/moffel"
	"github.com/xantios/koekeloere/piertje"
)

// Global options
var verbose *bool
var watched []string
var hooks []string

// Create a new instance of the logger.
var log = logrus.New()

// main Parse CLI options and run watcher
func main() {

	verbose = flag.Bool("v", false, "verbose mode")
	watchDirs := flag.String("w", ".", "Dirs to watch, comma seperated")
	hooks = flag.Args()

	flag.Parse()

	if *verbose {
		log.Warn("Running in verbose mode")
		log.SetLevel(logrus.DebugLevel)
	}

	// Setup Moffel to send out messages
	moffel.SetVerbose(verbose)
	moffel.SetLogger(log)
	moffel.Init(hooks)

	// Make a channel to pull async messages from
	channel := make(chan string)

	// Setup piertje to watch files
	piertje.SetLogger(log)
	piertje.SetVerbose(verbose)
	piertje.SetPaths(*watchDirs)
	piertje.SetChannel(channel)

	go piertje.Run()

	for {
		msg := <-channel
		split := strings.SplitN(msg, ":", 2)
		moffel.Emit(split[0], split[1])
	}
}
