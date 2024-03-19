package main

import (
	"github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

func main() {
	// parse config
	config := ParseConfig("config.yaml")
	updPrd := time.Duration(config.CacheTime) * time.Second
	slices.SortFunc(config.UnsolvedBorders, func(a, b *UnsolvedBorder) int {
		return a.Count - b.Count
	})

	// cache
	store := persist.NewMemoryStore(updPrd)
	// release mode?
	//gin.SetMode(gin.ReleaseMode)

	// router
	router := gin.Default()

	// I think now we don't need this... https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies
	router.ForwardedByClientIP = false

	router.LoadHTMLGlob("templates/page.html")
	// data
	criterionTitles, userValues := GetDeadlineResults(config)
	lastUpdate := time.Now()
	// funcs
	update := func() {
		if time.Since(lastUpdate).Seconds() > config.CacheTime {
			criterionTitles, userValues = GetDeadlineResults(config)
			lastUpdate = time.Now()
		}
	}

	// routes
	router.Static("/static", "./static")
	router.StaticFile("favicon.jpg", "./static/favicon.jpg")
	router.GET("/", cache.CacheByRequestURI(store, updPrd), func(c *gin.Context) {
		update()
		c.HTML(http.StatusOK, "page.html", gin.H{
			"CriterionTitles": criterionTitles,
			"UsersMap":        userValues,
		})
	})

	router.GET("/search/:name", cache.CacheByRequestURI(store, updPrd), func(c *gin.Context) {
		update()
		name := c.Param("name")
		for _, el := range userValues {
			if el.FullName == name {
				c.HTML(http.StatusOK, "page.html", gin.H{
					"CriterionTitles": criterionTitles,
					"UsersMap":        []*UserValues{el},
				})
				return
			}
		}
		c.String(http.StatusNotFound, "Nothing found with name=\"%s\"", name)
	})
	err := router.Run(config.ServerAddressPort)
	if err != nil {
		slog.Error("Server down with error: %s", err)
		panic(err)
	}
}
