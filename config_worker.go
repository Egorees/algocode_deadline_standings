package main

import (
	"encoding/json"
	"log"
	"os"
	"slices"
)

type UnsolvedBorder struct {
	Count int    `json:"count"`
	Color string `json:"color"`
}

type Config struct {
	CacheTime       float64          `json:"cache_time"`
	FullSolveText   string           `json:"full_solve_text"`
	UnsolvedBorders []UnsolvedBorder `json:"unsolved_borders"`
}

func ParseConfig(filepath string) Config {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error during opening config: %v", err.Error())
	}
	parser := json.NewDecoder(file)
	var res Config
	if err := parser.Decode(&res); err != nil {
		log.Fatalf("Error during parsing config: %v", err.Error())
	}
	return res
}

func GetColorByCount(config *Config, count int) string {
	ind := slices.IndexFunc(config.UnsolvedBorders, func(border UnsolvedBorder) bool {
		return border.Count >= count
	})
	if ind < 0 {
		log.Printf("Strange config: GetColorByCount returned -1; replacing it with len() - 1")
		ind = len(config.UnsolvedBorders) - 1
	}
	return config.UnsolvedBorders[ind].Color
}
