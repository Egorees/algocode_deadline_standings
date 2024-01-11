package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/page.html")
	router.GET("/", func(c *gin.Context) {
		criterionTitles, usersValues := GetDeadlineResults()
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
