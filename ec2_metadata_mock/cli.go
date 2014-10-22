package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

type config struct {
	Address string `short:"a" long:"address" description:"Address to bind to. Default is 127.0.0.1"`
	File    string `short:"f" long:"file" description:"File to serve (json format). Default is 'metadata.json'"`
	Port    string `short:"p" long:"port" description:"Port to bind to. Default is 8080"`
	Version func() `short:"v" long:"version" description:"Display the version"`
}

func setDefaults(c *config) *config {
	if c.Address == "" {
		c.Address = "127.0.0.1"
	}

	if c.File == "" {
		c.File = "./metadata.json"
	}

	if c.Port == "" {
		c.Port = "8080"
	}

	return c
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

	return setDefaults(c)
}
