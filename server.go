package main

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

func main() {
	// parse config
	config := ParseConfig("config.yaml")
	slices.SortFunc(config.UnsolvedBorders, func(a, b *UnsolvedBorder) int {
		return a.Count - b.Count
	})

	// router
	router := gin.Default()
	router.LoadHTMLGlob("templates/page.html")
	// data
	criterionTitles, usersValues := GetDeadlineResults(config)
	lastUpdate := time.Now()
	// routes
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		if time.Since(lastUpdate).Seconds() > config.CacheTime {
			criterionTitles, usersValues = GetDeadlineResults(config)
			lastUpdate = time.Now()
		}
		c.HTML(http.StatusOK, "page.html", gin.H{
			"CriterionTitles": criterionTitles,
			"Users":           usersValues,
		})
	})

	err := router.Run(config.ServerAddressPort)
	if err != nil {
		slog.Error("Server down with error: %s", err)
		panic(err)
	}
}
