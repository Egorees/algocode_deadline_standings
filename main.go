package main

import (
	"algocode_deadline_standings/configs"
	"algocode_deadline_standings/server"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
)

func main() {
	// parse config
	parser := argparse.NewParser("Algocode Deadline Table",
		"Deadline table for Yandex Kruzhok's parallel B'")
	// this argparse sucks...
	configPath := parser.String("c", "config", &argparse.Options{
		Required: true,
		Help:     "Path to config file",
	})
	// maybe should add path to the deadline file too
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Printf(parser.Usage(err))
		os.Exit(1)
	}
	config := configs.ParseConfig(*configPath)

	server.RunServer(config)
}
