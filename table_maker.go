package main

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func GetDeadlineResults(config *Config) ([]string, map[string]*UserValues) {
	criterionTitles := []string{"Не решено"}

	mapUsersValues := make(map[string]*UserValues)
	_map := treemap.NewWith(utils.StringComparator)
	//_map.tree

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
		}
		criterionTitles = append(criterionTitles, contest.Title)
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
	for _, user := range data.Users {
		cur := result[strconv.Itoa(user.Id)]

		// why not just strings.Title()? Just because. https://pkg.go.dev/strings#Title
		// user.Name = strings.Title(user.Name)
		user.Name = cases.Title(language.Russian).String(user.Name)

		_map.Put(user.Name, &UserValues{
			Name:   user.Name,
			Values: []*Value{},
		})
		mapUsersValues[user.Name] = &UserValues{
			Name:   user.Name,
			Values: []*Value{},
		}

		unsolvedColor := config.GetColorByCount(cur.total)
		//_map.
		mapUsersValues[user.Name].Values = append(mapUsersValues[user.Name].Values,
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
			mapUsersValues[user.Name].Values = append(mapUsersValues[user.Name].Values,
				&Value{
					Value: tasksInString,
					Color: valueColor,
				},
			)
		}
	}

	return criterionTitles, mapUsersValues
}
