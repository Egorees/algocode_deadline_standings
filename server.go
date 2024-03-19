package main

import (
	"github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"
)

func main() {
	// parse config
	config := ParseConfig("config.yaml")
	updPrd := time.Duration(config.CacheTime) * time.Second
	slices.SortFunc(config.UnsolvedBorders, func(a, b *UnsolvedBorder) int {
		return a.Count - b.Count
	})
	lock := sync.RWMutex{}

	// cache
	store := persist.NewMemoryStore(updPrd)

	// release mode?
	gin.SetMode(gin.ReleaseMode)

	// router
	router := gin.Default()

	// I think now we don't need this... https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies
	err := router.SetTrustedProxies(nil)
	if err != nil {
		slog.Error("Cant set trusted proxies: %s", err)
		panic(err)
	}

	// templates
	router.LoadHTMLGlob("templates/page.gohtml")
	// data
	criterionTitles, userValues := GetDeadlineResults(config)
	lastUpdate := time.Now()
	// funcs
	update := func() {
		lock.Lock()
		if time.Since(lastUpdate).Seconds() > config.CacheTime {
			criterionTitles, userValues = GetDeadlineResults(config)
			lastUpdate = time.Now()
		}
		lock.Unlock()
	}

	// routes (static)
	router.Static("/static", "./static")
	//router.StaticFile("favicon.jpg", "./static/favicon.jpg")
	// table routes
	router.GET("/", cache.CacheByRequestURI(store, updPrd), func(c *gin.Context) {
		update()
		lock.RLock()
		c.HTML(http.StatusOK, "page.gohtml", gin.H{
			"CriterionTitles": criterionTitles,
			"UserValues":      userValues,
			"Single":          false,
		})
		lock.RUnlock()
	})

	router.GET("/search/:name", cache.CacheByRequestURI(store, updPrd), func(c *gin.Context) {
		update()
		lock.RLock()
		name := c.Param("name")
		ind, found := slices.BinarySearchFunc(userValues, name, func(values *UserValues, s string) int {
			return strings.Compare(values.FullName, s)
		})
		if found {
			c.HTML(http.StatusOK, "page.gohtml", gin.H{
				"CriterionTitles": criterionTitles,
				"UserValues":      []*UserValues{userValues[ind]},
				"Single":          true,
			})
		} else {
			c.String(http.StatusNotFound, "Nothing found with name=\"%s\"", name)
		}
		lock.RUnlock()
	})

	// run server
	err = router.Run(config.ServerAddressPort)
	if err != nil {
		slog.Error("Server down with error: %s", err)
		panic(err)
	}
}
