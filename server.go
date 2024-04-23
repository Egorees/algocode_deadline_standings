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
	if config.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// router
	router := gin.Default()

	// variables
	var err error
	var criterionTitles []*CriterionTitle
	var userValues []*UserValues
	var lastUpdate time.Time
	var stats map[int]*Stats

	// I think now we don't need this... https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies
	err = router.SetTrustedProxies(nil)
	if err != nil {
		slog.Error("Cant set trusted proxies: %s", err)
		panic(err)
	}

	// templates
	router.LoadHTMLGlob("templates/*.gohtml")

	// funcs
	update := func() {
		lock.Lock()
		if time.Since(lastUpdate).Seconds() > config.CacheTime {
			criterionTitles, userValues, err = GetDeadlineResults(config)
			if err == nil {
				stats, err = statisticsFun(config, userValues)
			}
			lastUpdate = time.Now()
		}
		lock.Unlock()
	}
	errorCheckerHandler := func(c *gin.Context) {
		update()
		lock.RLock()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.gohtml", gin.H{
				"Error": err.Error(),
			})
			c.Abort()
		}
		lock.RUnlock()
	}

	// first update
	update()

	// routes (static)
	router.Static("/static", "./static")

	// table routes
	router.GET("/", cache.CacheByRequestURI(store, updPrd), errorCheckerHandler, func(c *gin.Context) {
		update()
		lock.RLock()
		c.HTML(http.StatusOK, "page.gohtml", gin.H{
			"CriterionTitles": criterionTitles[2:],
			"UserValues":      userValues,
			"Single":          false,
		})
		lock.RUnlock()
	})

	router.GET("/search/:name", cache.CacheByRequestURI(store, updPrd), errorCheckerHandler, func(c *gin.Context) {
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

	router.GET("/stats", cache.CacheByRequestURI(store, updPrd), errorCheckerHandler, func(c *gin.Context) {
		lock.RLock()
		c.HTML(http.StatusOK, "stats.gohtml", gin.H{
			"Stats": stats,
		})
		lock.RUnlock()
	})

	//router.GET("/tasks/stat", cache.CacheByRequestURI(store, updPrd), func(c *gin.Context) {
	//	lock.RLock()
	//	c.HTML(http.StatusOK, "tasks_stat.gohtml", gin.H{
	//		"CriterionTitles": criterionTitles,
	//		"DeadlineTasks":
	//	})
	//	lock.RUnlock()
	//})

	// run server
	err = router.Run(config.ServerAddressPort)
	if err != nil {
		slog.Error("Server down with error: %s", err)
		panic(err)
	}
}
