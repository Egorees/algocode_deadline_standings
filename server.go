package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserValues struct {
	Name   string
	Values []string
}

var criterionTitles []string
var usersValues []UserValues

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/page.html")

	Main()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "page.html", gin.H{
			"CriterionTitles": criterionTitles,
			"Users":           usersValues,
		})
	})

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Server down with error: %s", err)
		return
	}
}
