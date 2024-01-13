package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
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

type DeadlineTasks = map[string][]string

type DeadlineData struct {
	Tasks DeadlineTasks `json:"deadline"`
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

type SubmitsData struct {
	Users    []User    `json:"users"`
	Contests []Contest `json:"contests"`
}

type UnsolvedData struct {
	total    int
	unsolved []TasksFromContest
	// maybe we will need more data
}

type UserValues struct {
	Name     string
	Unsolved int
	Values   []string
}

func getSubmitsData(url string) SubmitsData {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error while quering algocode: %v\n", err.Error())
	}
	parser := json.NewDecoder(res.Body)
	var data SubmitsData
	if err = parser.Decode(&data); err != nil {
		log.Fatalf("Error while parsing json from algocode: %v\n", err.Error())
	}
	return data
}

func GetDeadlineResults() ([]string, []UserValues) {
	var criterionTitles []string
	var usersValues []UserValues
	// think of making this link shorter
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
		usersValues = append(usersValues, UserValues{Name: user.Name, Values: []string{}, Unsolved: cur.total})
		for _, tasksFromContest := range cur.unsolved {
			tasksInString := strings.Join(tasksFromContest.Tasks[:], ",")
			if tasksInString == "" {
				tasksInString = "Всё решил!"
			}
			usersValues[ind].Values = append(usersValues[ind].Values, tasksInString)
		}
	}

	return criterionTitles, usersValues
}
