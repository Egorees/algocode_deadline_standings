package data_processors

import (
	"algocode_deadline_standings/algocode"
	"algocode_deadline_standings/configs"
	"fmt"
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

	needProblems := configs.ParseDeadlineProblems(config.DeadlineFilepath)

	for _, contest := range data.Contests {
		needProblemsInds := make([]int, len(needProblems.Problems[contest.Title]))
		if len(needProblemsInds) == 0 {
			continue
		}
		criterionTitles = append(criterionTitles, &CriterionTitle{
			Title:    contest.Title,
			EjudgeId: contest.EjudgeId,
		})
		for ind, needProblem := range needProblems.Problems[contest.Title] {
			problemInd := slices.IndexFunc(contest.Problems, func(problem *algocode.Problem) bool {
				return problem.Short == needProblem
			})
			if problemInd == -1 {
				// skipping it actually, maybe that will break something :)
				// well, it seems it's ok
				slog.Error(fmt.Sprintf("Not found problem '%v' in '%v'\n", needProblem, contest.Title))
			} else {
				needProblemsInds[ind] = problemInd
			}
		}
		for user, problems := range contest.Users {
			problemsFromContest := &ProblemsFromContest{
				ContestTitle: contest.Title,
				Problems:     make([]string, 0),
			}
			for indInNeedProblems, needProblem := range needProblemsInds {
				if problems[needProblem].Score == 0 {
					problemsFromContest.Problems = append(problemsFromContest.Problems,
						needProblems.Problems[contest.Title][indInNeedProblems])
				}
			}
			result[user].unsolved = append(result[user].unsolved, problemsFromContest)
			result[user].total += len(problemsFromContest.Problems)
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

		for _, problemsFromContest := range cur.unsolved {
			var valueColor string
			specialValueColor := ""
			requiredProblemsFromContest := ProblemsFromContest{
				ContestTitle: problemsFromContest.ContestTitle,
				Problems:     make([]string, 0),
			}
			notRequiredProblemsFromContest := ProblemsFromContest{
				ContestTitle: problemsFromContest.ContestTitle,
				Problems:     make([]string, 0),
			}
			for _, problem := range problemsFromContest.Problems {
				if slices.Contains(needProblems.RequiredProblems[problemsFromContest.ContestTitle], problem) {
					requiredProblemsFromContest.Problems = append(requiredProblemsFromContest.Problems, problem)
				} else {
					notRequiredProblemsFromContest.Problems = append(notRequiredProblemsFromContest.Problems, problem)
				}
			}

			requiredProblemsInString := strings.Join(requiredProblemsFromContest.Problems, ", ")
			notRequiredProblemsInString := strings.Join(notRequiredProblemsFromContest.Problems, ", ")
			if requiredProblemsInString == "" && notRequiredProblemsInString == "" {
				notRequiredProblemsInString = config.FullSolveText
				valueColor = config.UnsolvedBorders[0].Color
			} else {
				valueColor = config.UnsolvedBorders[len(config.UnsolvedBorders)-1].Color
				if requiredProblemsInString != "" && notRequiredProblemsInString != "" {
					notRequiredProblemsInString += ", "
				}
				specialValueColor = "#FF0000"
			}
			usersValues[ind].Values = append(usersValues[ind].Values,
				&Value{
					Value:        notRequiredProblemsInString,
					Color:        valueColor,
					SpecialValue: requiredProblemsInString,
					SpecialColor: specialValueColor,
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
