package configs

import (
	"fmt"
	"log/slog"
	"os"
	"slices"

	"gopkg.in/yaml.v3"
)

func ParseConfig(filepath string) *Config {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Error(fmt.Sprintf("Error during opening config: %v", err.Error()))
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	parser := yaml.NewDecoder(file)
	var res Config
	if err := parser.Decode(&res); err != nil {
		slog.Error(fmt.Sprintf("Error during parsing config: %v", err.Error()))
		panic(err)
	}
	return &res
}

func (config *Config) GetColorByCount(count int) string {
	ind := slices.IndexFunc(config.UnsolvedBorders, func(border *UnsolvedBorder) bool {
		return border.Count >= count
	})
	if ind < 0 {
		slog.Warn("Strange config: GetColorByCount returned -1; replacing it with len() - 1")
		ind = len(config.UnsolvedBorders) - 1
	}
	return config.UnsolvedBorders[ind].Color
}

func ParseDeadlineTasks(filepath string) DeadlineData {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Error(fmt.Sprintf("Error during opening deadline tasks file: %v", err.Error()))
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
		slog.Error(fmt.Sprintf("Error during parsing deadline tasks: %v", err.Error()))
		panic(err)
	}
	return res
}
