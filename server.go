package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"slices"
	"time"
)

func main() {
	config := ParseConfig("config.json")
	slices.SortFunc(config.UnsolvedBorders, func(a, b UnsolvedBorder) int {
		return a.Count - b.Count
	})

	router := gin.Default()
	router.LoadHTMLGlob("templates/page.html")
	criterionTitles, usersValues := GetDeadlineResults(&config)
	lastUpdate := time.Now()
	router.GET("/", func(c *gin.Context) {
		if time.Since(lastUpdate).Seconds() > config.CacheTime {
			criterionTitles, usersValues = GetDeadlineResults(&config)
			lastUpdate = time.Now()
		}
		c.HTML(http.StatusOK, "page.html", gin.H{
			"CriterionTitles": criterionTitles,
			"Users":           usersValues,
		})
	})

	err := router.Run("127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Server down with error: %s", err)
		return
	}
}
