package main

import (
	"log/slog"
	"slices"
	"strconv"
	"strings"
)

type TasksFromContest struct {
	ContestTitle string
	Tasks        []string
}

type UnsolvedData struct {
	total    int
	unsolved []*TasksFromContest
	// maybe we will need more data
}

type Value struct {
	Value string
	Color string
}

type UserValues struct {
	Name   string
	Values []*Value
}

func GetDeadlineResults(config *Config) ([]string, []*UserValues, map[string]*UserValues) {
	criterionTitles := []string{"Не решено"}

	var usersValues []*UserValues
	data := getSubmitsData(config.SubmitsLink)

	result := make(map[string]*UnsolvedData, len(data.Users))

	for _, user := range data.Users {
		result[strconv.Itoa(user.Id)] = &UnsolvedData{}
	}

	needTasks := ParseDeadlineTasks(config.DeadlineFilepath)

	for _, contest := range data.Contests {
		needTasksInds := make([]int, len(needTasks.Tasks[contest.Title]))
		if len(needTasksInds) == 0 {
			continue
		} else {
			criterionTitles = append(criterionTitles, contest.Title)
		}
		for ind, needTask := range needTasks.Tasks[contest.Title] {
			taskInd := slices.IndexFunc(contest.Problems, func(problem *Problem) bool {
				return problem.Short == needTask
			})
			if taskInd == -1 {
				// skipping it actually, maybe that will break something :)
				// well, it seems it's ok
				slog.Error("Not found task %v in %v\n", needTask, contest.Title)
			} else {
				needTasksInds[ind] = taskInd
			}
		}
		for user, tasks := range contest.Users {
			tasksFromContest := &TasksFromContest{
				ContestTitle: contest.Title,
				Tasks:        make([]string, 0),
			}
			for indInNeedTasks, needTask := range needTasksInds {
				if tasks[needTask].Score == 0 {
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
		usersValues = append(usersValues,
			&UserValues{
				Name:   user.Name,
				Values: []*Value{},
			},
		)

		unsolvedColor := config.GetColorByCount(cur.total)
		usersValues[ind].Values = append(usersValues[ind].Values,
			&Value{
				Value: strconv.Itoa(cur.total),
				Color: unsolvedColor,
			},
		)

		for _, tasksFromContest := range cur.unsolved {
			var valueColor string
			tasksInString := strings.Join(tasksFromContest.Tasks[:], ",")
			if tasksInString == "" {
				tasksInString = config.FullSolveText
				valueColor = config.UnsolvedBorders[0].Color
			} else {
				valueColor = config.UnsolvedBorders[len(config.UnsolvedBorders)-1].Color
			}
			usersValues[ind].Values = append(usersValues[ind].Values,
				&Value{
					Value: tasksInString,
					Color: valueColor,
				},
			)
		}
	}

	// some changes on usersValues
	newUsersValues := make(map[string]*UserValues)

	for _, uv := range usersValues {
		newUsersValues[uv.Name] = uv
	}

	return criterionTitles, usersValues, newUsersValues
}
