package main

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type DeadlineTasks = map[string][]string

type DeadlineData struct {
	Tasks DeadlineTasks `yaml:"deadline"`
}

func ParseDeadlineTasks(filepath string) DeadlineData {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Error("Error during opening deadline tasks file: %v", err.Error())
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	parser := yaml.NewDecoder(file)
	var res DeadlineData
	if err := parser.Decode(&res); err != nil {
		slog.Error("Error during parsing deadline tasks: %v", err.Error())
		panic(err)
	}
	return res
}
