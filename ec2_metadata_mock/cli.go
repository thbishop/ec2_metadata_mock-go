package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

type config struct {
	Address string `short:"a" long:"address" description:"Address to bind to" default:"127.0.0.1" value-name:"ADDRESS"`
	File    string `short:"f" long:"file" description:"File to serve (json format)" default:"metadata.json" value-name:"PATH_TO_FILE"`
	Port    string `short:"p" long:"port" description:"Port to bind to" default:"8080" value-name:"PORT"`
	Version func() `short:"v" long:"version" description:"Display the version"`
}

func parseCliArgs() *config {
	c := &config{}

	c.Version = func() {
		fmt.Println(version)
		os.Exit(0)
	}

	parser := flags.NewParser(c, flags.Default)

	args, err := parser.Parse()
	if err != nil {
		helpDisplayed := false

		for _, i := range args {
			if i == "-h" || i == "--help" {
				helpDisplayed = true
				break
			}
		}

		if !helpDisplayed {
			parser.WriteHelp(os.Stderr)
		}
		os.Exit(1)
	}

	return c
}
