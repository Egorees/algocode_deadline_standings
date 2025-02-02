package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"slices"
)

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

func ParseDeadlineProblems(filepath string) DeadlineData {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Error("Error during opening deadline problems file: %v", err.Error())
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
		slog.Error("Error during parsing deadline problems: %v", err.Error())
		panic(err)
	}
	for contest, problems := range res.RequiredProblems {
		notRequiredProblems, ok := res.Problems[contest]
		if !ok {
			panic(fmt.Sprintf("Can't find contest '%v' in deadline", contest))
		}
		for _, task := range problems {
			if !slices.Contains(notRequiredProblems, task) {
				panic(fmt.Sprintf("Can't find problem '%v' in contest '%v'", task, contest))
			}
		}
	}
	return res
}
