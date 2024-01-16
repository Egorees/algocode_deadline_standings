package main

import (
	"encoding/json"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type DeadlineTasks = map[string][]string

type DeadlineData struct {
	Tasks DeadlineTasks `json:"deadline"`
}

type TasksFromContest struct {
	ContestTitle string
	Tasks        []string
}

type UnsolvedData struct {
	total    int
	unsolved []TasksFromContest
	// maybe we will need more data
}

type Value struct {
	Value string
	Color string
}

type UserValues struct {
	Name   string
	Values []Value
}

func parseDeadlineTasks(filepath string) DeadlineData {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error during opening deadline tasks file: %v", err.Error())
	}
	parser := json.NewDecoder(file)
	var res DeadlineData
	if err := parser.Decode(&res); err != nil {
		log.Fatalf("Error during parsing deadline tasks: %v", err.Error())
	}
	return res
}

func GetDeadlineResults(config *Config) ([]string, []UserValues) {
	criterionTitles := []string{"Не решено"}

	var usersValues []UserValues
	// think of making this link shorter/passing it
	data := getSubmitsData("https://algocode.ru/standings_data/bp_fall_2023/")

	result := make(map[string]*UnsolvedData, len(data.Users))

	for _, user := range data.Users {
		result[strconv.Itoa(user.Id)] = &UnsolvedData{}
	}

	needTasks := parseDeadlineTasks("deadline.json")

	for _, contest := range data.Contests {
		needTasksInds := make([]int, len(needTasks.Tasks[contest.Title]))
		if len(needTasksInds) == 0 {
			continue
		} else {
			criterionTitles = append(criterionTitles, contest.Title)
		}
		for ind, needTask := range needTasks.Tasks[contest.Title] {
			taskInd := slices.IndexFunc(contest.Problems, func(problem Problem) bool {
				return problem.Short == needTask
			})
			if taskInd == -1 {
				log.Fatalf("Not found task %v in %v\n", needTask, contest.Title)
			} else {
				needTasksInds[ind] = taskInd
			}
		}
		for user, tasks := range contest.Users {
			tasksFromContest := TasksFromContest{contest.Title, make([]string, 0)}
			for indInNeedTasks, needTask := range needTasksInds {
				// Maybe this should be changed to tasks[needTask].Score == 1
				if tasks[needTask].Verdict != "OK" {
					tasksFromContest.Tasks = append(tasksFromContest.Tasks,
						needTasks.Tasks[contest.Title][indInNeedTasks])
				}
			}
			result[user].unsolved = append(result[user].unsolved, tasksFromContest)
			result[user].total += len(tasksFromContest.Tasks)
		}
	}
	for ind, user := range data.Users {
		cur := result[strconv.Itoa(user.Id)]
		usersValues = append(usersValues, UserValues{Name: user.Name, Values: []Value{}})

		unsolvedColor := GetColorByCount(config, cur.total)
		usersValues[ind].Values = append(usersValues[ind].Values, Value{Value: strconv.Itoa(cur.total), Color: unsolvedColor})

		for _, tasksFromContest := range cur.unsolved {
			var valueColor string
			tasksInString := strings.Join(tasksFromContest.Tasks[:], ",")
			if tasksInString == "" {
				tasksInString = config.FullSolveText
				valueColor = config.UnsolvedBorders[0].Color
			} else {
				valueColor = config.UnsolvedBorders[len(config.UnsolvedBorders)-1].Color
			}
			usersValues[ind].Values = append(usersValues[ind].Values, Value{Value: tasksInString, Color: valueColor})
		}
	}

	return criterionTitles, usersValues
}
