package main

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"slices"
)

type UnsolvedBorder struct {
	Count int    `json:"count"`
	Color string `json:"color"`
}

type Config struct {
	CacheTime         float64           `yaml:"cache_time"`
	FullSolveText     string            `yaml:"full_solve_text"`
	UnsolvedBorders   []*UnsolvedBorder `yaml:"unsolved_borders"`
	ServerAddressPort string            `yaml:"server_address_port"`
	SubmitsLink       string            `yaml:"submits_link"`
	DeadlineFilepath  string            `yaml:"deadline_filepath"`
}

func ParseConfig(filepath string) *Config {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Error("Error during opening config: %v", err.Error())
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
		slog.Error("Error during parsing config: %v", err.Error())
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
