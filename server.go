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

	// release mode?
	//gin.SetMode(gin.ReleaseMode)

	// router
	router := gin.Default()

	// I think now we don't need this... https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies
	router.ForwardedByClientIP = false

	router.LoadHTMLGlob("templates/page.html")
	// data
	criterionTitles, mapUsersValues := GetDeadlineResults(config)
	lastUpdate := time.Now()
	// funcs
	update := func() {
		if time.Since(lastUpdate).Seconds() > config.CacheTime {
			criterionTitles, mapUsersValues = GetDeadlineResults(config)
			lastUpdate = time.Now()
		}
	}
	// routes
	router.Static("/static", "./static")
	router.StaticFile("favicon.jpg", "./static/favicon.jpg")
	router.GET("/", func(c *gin.Context) {
		update()
		c.HTML(http.StatusOK, "page.html", gin.H{
			"CriterionTitles": criterionTitles,
			"UsersMap":        mapUsersValues,
		})
	})
	router.GET("/name/:name", func(c *gin.Context) {
		update()
		name := c.Param("name")
		if val, ok := mapUsersValues[name]; ok {
			c.HTML(http.StatusOK, "page.html", gin.H{
				"CriterionTitles": criterionTitles,
				"UsersMap":        []*UserValues{val},
			})
		} else {
			c.String(http.StatusNotFound, "name \"%s\" not found", name)
		}
	})

	err := router.Run(config.ServerAddressPort)
	if err != nil {
		slog.Error("Server down with error: %s", err)
		panic(err)
	}
}
