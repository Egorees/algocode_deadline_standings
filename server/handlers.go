package server

import (
	"algocode_deadline_standings/configs"
	processors "algocode_deadline_standings/data-processors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func getData(config *configs.Config, needStats bool) (
	[]*processors.CriterionTitle,
	[]*processors.UserValues,
	map[int]*processors.Stats,
	error) {

	// hope we can do it on every request
	criterionTitles, userValues, err := processors.GetDeadlineResults(config)
	var stats map[int]*processors.Stats
	if needStats && err == nil {
		stats, err = processors.CreateStatistics(config, userValues)
	}
	// yes, that will return trash in stats if we don't need it or an error occurred
	return criterionTitles, userValues, stats, err
}

func errorChecker(c *gin.Context, err error) { // just shows an error if something happened
	if err == nil {
		return
	}
	c.HTML(http.StatusInternalServerError, "error.gohtml", gin.H{
		"Error": err.Error(),
	})
	c.Abort()
}

func mainPage(c *gin.Context, config *configs.Config) {
	criterionTitles, userValues, _, err := getData(config, false)
	errorChecker(c, err)
	c.HTML(http.StatusOK, "page.gohtml", gin.H{
		"CriterionTitles": criterionTitles,
		"UserValues":      userValues,
		"Single":          false,
	})
}

func studentStats(c *gin.Context, config *configs.Config) {
	// don't know how to avoid these two lines
	criterionTitles, userValues, _, err := getData(config, false)
	errorChecker(c, err)

	name := c.Param("name")
	ind, found := slices.BinarySearchFunc(userValues, name, func(values *processors.UserValues, s string) int {
		return strings.Compare(values.FullName, s)
	})
	if found {
		c.HTML(http.StatusOK, "page.gohtml", gin.H{
			"CriterionTitles": criterionTitles,
			"UserValues":      []*processors.UserValues{userValues[ind]},
			"Single":          true,
		})
	} else {
		c.String(http.StatusNotFound, "Nothing found with name=\"%s\"", name)
	}
}

func allStats(c *gin.Context, config *configs.Config) {
	_, _, stats, err := getData(config, true)
	errorChecker(c, err)

	c.HTML(http.StatusOK, "stats.gohtml", gin.H{
		"Stats": stats,
	})
}

func studentWhoPass(c *gin.Context, config *configs.Config) {
	_, _, stats, err := getData(config, true)
	errorChecker(c, err)
	problemsPass, err := strconv.Atoi(c.Param("pass"))
	errorChecker(c, err)
	totalCnt := 0
	resString := strings.Builder{}
	for problemsCnt, users := range stats {
		if problemsCnt > problemsPass {
			continue
		}
		totalCnt += users.Count
		resString.WriteString(strings.Join(users.Peoples, "\n"))
		resString.WriteByte('\n')
	}
	c.String(http.StatusOK, fmt.Sprintf("Total: %v\n%v", totalCnt, resString.String()))
}

func studentWhoNotPass(c *gin.Context, config *configs.Config) {
	_, _, stats, err := getData(config, true)
	errorChecker(c, err)
	problemsPass, err := strconv.Atoi(c.Param("pass"))
	errorChecker(c, err)
	totalCnt := 0
	resString := strings.Builder{}
	for problemsCnt, users := range stats {
		if problemsCnt <= problemsPass {
			continue
		}
		totalCnt += users.Count
		resString.WriteString(strings.Join(users.Peoples, "\n"))
		resString.WriteByte('\n')
	}
	c.String(http.StatusOK, fmt.Sprintf("Total: %v\n%v", totalCnt, resString.String()))
}
