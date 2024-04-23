package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"net/http"
	"strings"
	"time"
)

// Global variables
var linkCache = make(map[string][]string)

// Jika mengandung salah satu dari identifier seperti File: pada URL, return true
func checkIgnoredLink(url string) bool {
	ignoredLinks := [...]string{"/File:", "/Special:", "/Template:", "/Template_page:", "/Help:", "/Category:", "Special:", "/Wikipedia:", "/Portal:", "/Talk:"}
	for _, st := range ignoredLinks {
		if strings.Contains(url, st) {
			return true
		}
	}
	return false
}

func getAllLinks(url string) []string {
	//c := colly.NewCollector(colly.CacheDir("./cache"))
	c := colly.NewCollector()

	// Inisialisasi array
	var links []string

	// Cari semua link dan kalau berawalan /wiki/ ditambahkan, dan jika ada yang mengandung ignoredLinks diabaikan
	c.OnHTML("div#mw-content-text a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasPrefix(link, "/wiki/") && !checkIgnoredLink(link) {
			links = append(links, "https://en.wikipedia.org"+link)
		}
	})

	err := c.Visit(url)
	if err != nil {
		return nil
	}	
	return links
}


func cacheLinks(url string) []string {
	links, exists := linkCache[url]

	if exists {
		//fmt.Println("Menggunakan cache untuk", url)
		return links
	}

	links = getAllLinks(url)
	linkCache[url] = links
	return links
}


func DLS(currentURL string, targetURL string, limit int, result *[]string, numOfArticles *int) bool {
	*numOfArticles++
	*result = append(*result, currentURL)
	if currentURL == targetURL {
		return true
	}

	if limit <= 1 {
		*result = (*result)[:len(*result)-1]
		return false
	}

	links := cacheLinks(currentURL)
	//links := getAllLinks(currentURL)

	for _, link := range links {
		//fmt.Println("Cek link", link, "di level", limit)
		if DLS(link, targetURL, limit-1, result, numOfArticles) {
			return true
		}
	}
	*result = (*result)[:len(*result)-1]
	return false
}

func IDS(startURL string, targetURL string, maxDepth int, result *[]string, numOfArticles *int) bool {
	i := 1
	for {
		if DLS(startURL, targetURL, i, result, numOfArticles) {
			return true
		}
		fmt.Println(i)
		i++
		if i == maxDepth {
			return false // Safe condition only
		}
	}
}

// Not used
func callIDS() {
	router := gin.Default()
	router.GET("/IDS", func(c *gin.Context) {
		source := c.Query("source")
		target := c.Query("target")
		fmt.Println("Source", source, "Target", target)
		var startURL, targetURL string
		startURL = "https://en.wikipedia.org/wiki/" + source
		targetURL = "https://en.wikipedia.org/wiki/" + target
		fmt.Println("Start URL", startURL)

		start := time.Now()
		numOfArticles := 0
		result := make([]string, 0)
		var elapsedTime time.Duration
		var end time.Time
		if IDS(startURL, targetURL, 15, &result, &numOfArticles) {
			end = time.Now()
			elapsedTime = end.Sub(start)
			fmt.Println("Waktu eksekusi", end.Sub(start))
			for i, r := range result {
				fmt.Printf("Link ke %d: %s\n", i, r)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"numOfArticles": numOfArticles,
			"result":        result,
			"length":        len(result),
			"elapsedTime":   elapsedTime,
		})
	})
	err := router.Run(":8080")
	if err != nil {
		return
	}
	//err = os.Remove("./cache")
	//if err != nil {
	//	fmt.Println("Error removing cache")
	//}
}
