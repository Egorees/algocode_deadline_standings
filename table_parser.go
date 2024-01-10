package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Group      string `json:"group"`
	GroupShort string `json:"group_short"`
}

type Problem struct {
	Id    string `json:"id"`
	Long  string `json:"long"`
	Short string `json:"short"`
	Index int    `json:"index"`
}

type UserSubmit struct {
	Score   int    `json:"score"`
	Penalty int    `json:"penalty"`
	Verdict string `json:"verdict"`
	Time    int    `json:"time"`
}

type Contest struct {
	Id          int                     `json:"id"`
	Date        string                  `json:"date"`
	EjudgeId    int                     `json:"ejudge_id"`
	Title       string                  `json:"title"`
	Coefficient float64                 `json:"coefficient"`
	Problems    []Problem               `json:"problems"`
	Users       map[string][]UserSubmit `json:"users"`
}

type TasksFromContest struct {
	ContestTitle string
	Tasks        []string
}

func FindProblemInd(problems *[]Problem, short string) int {
	for ind, problem := range *problems {
		if problem.Short == short {
			return ind
		}
	}
	return -1
}

func main() {

	var data struct {
		Users    []User    `json:"users"`
		Contests []Contest `json:"contests"`
	}

	dataFile, err := os.Open("data.json")
	if err != nil {
		log.Fatal(err)
	}

	jsonParser := json.NewDecoder(dataFile)
	if err = jsonParser.Decode(&data); err != nil {
		log.Fatal(err)
	}

	result := make(map[string][]TasksFromContest, len(data.Users))

	for _, user := range data.Users {
		result[strconv.Itoa(user.Id)] = []TasksFromContest{}
	}

	needTasks := map[string][]string{
		"Кратчайшие пути": {"A", "B", "C", "E"},
	}

	for _, contest := range data.Contests {
		needTasksInds := make([]int, len(needTasks[contest.Title]))
		if len(needTasksInds) == 0 {
			continue
		}
		for ind, needTask := range needTasks[contest.Title] {
			taskInd := FindProblemInd(&contest.Problems, needTask)
			if taskInd == -1 {
				log.Fatal("Not found task " + needTask + " in " + contest.Title)
			} else {
				needTasksInds[ind] = taskInd
			}
		}
		for user, tasks := range contest.Users {
			var tasksFromContest TasksFromContest
			tasksFromContest.ContestTitle = contest.Title
			for indInNeedTasks, needTask := range needTasksInds {
				if tasks[needTask].Verdict != "OK" {
					tasksFromContest.Tasks = append(tasksFromContest.Tasks, needTasks[contest.Title][indInNeedTasks])
				}
			}
			if len(tasksFromContest.Tasks) != 0 {
				result[user] = append(result[user], tasksFromContest)
			}
		}
	}
}
