package data_processors

import (
	"algocode_deadline_standings/algocode"
	"algocode_deadline_standings/configs"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log/slog"
	"slices"
	"strconv"
	"strings"
)

func GetDeadlineResults(config *configs.Config) ([]*CriterionTitle, []*UserValues, error) {
	criterionTitles := make([]*CriterionTitle, 0)

	var usersValues []*UserValues
	data := algocode.GetSubmitsData(config.SubmitsLink)
	if data == nil {
		slog.Error("Submits data is nil")
		return nil, nil, &DeadlineResultsError{Reason: "Submits data is nil"}
	}

	result := make(map[string]*UnsolvedData, len(data.Users))

	for _, user := range data.Users {
		result[strconv.Itoa(user.Id)] = &UnsolvedData{}
	}

	needTasks := configs.ParseDeadlineTasks(config.DeadlineFilepath)

	for _, contest := range data.Contests {
		needTasksInds := make([]int, len(needTasks.Tasks[contest.Title]))
		if len(needTasksInds) == 0 {
			continue
		}
		criterionTitles = append(criterionTitles, &CriterionTitle{
			Title:    contest.Title,
			EjudgeId: contest.EjudgeId,
		})
		for ind, needTask := range needTasks.Tasks[contest.Title] {
			taskInd := slices.IndexFunc(contest.Problems, func(problem *algocode.Problem) bool {
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

		// name parse
		// why not just strings.Title()? Just because. https://pkg.go.dev/strings#Title
		fullName := strings.Split(user.Name, " ")
		firstName := cases.Title(language.Russian).String(fullName[1])
		secondName := cases.Title(language.Russian).String(fullName[0])

		usersValues = append(usersValues,
			&UserValues{
				FirstName:  firstName,
				SecondName: secondName,
				FullName:   secondName + " " + firstName,
				Values:     []*Value{},
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
			tasksInString := strings.Join(tasksFromContest.Tasks[:], ", ")
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

	slices.SortFunc(usersValues, func(a, b *UserValues) int {
		if n := strings.Compare(a.SecondName, b.SecondName); n != 0 {
			return n
		}
		return strings.Compare(a.FirstName, b.FirstName)
	})

	return criterionTitles, usersValues, nil
}
