package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Tree struct {
}

// Global variables
var linkCache = make(map[string][]string)
var cacheMutex = &sync.Mutex{}
var sharedMutex = &sync.Mutex{}

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

func DLS(currentURL string, targetURL string, limit int, result *[]string, numOfArticles *int, wg *sync.WaitGroup) bool {
	defer wg.Done()

	if limit <= 1 {
		return false
	}

	sharedMutex.Lock()
	*numOfArticles++
	*result = append(*result, currentURL)
	sharedMutex.Unlock()
	if currentURL == targetURL {
		return true
	}

	links := cacheLinks(currentURL)
	//links := getAllLinks(currentURL)

	for _, link := range links {
		wg.Add(1)
		//fmt.Println("Cek link", link)
		if DLS(link, targetURL, limit-1, result, numOfArticles, wg) {
			return true
		}
	}
	sharedMutex.Lock()
	*result = (*result)[:len(*result)-1]
	sharedMutex.Unlock()
	return false
}

func IDS(startURL string, targetURL string, maxDepth int, result *[]string, numOfArticles *int) bool {
	i := 1
	var wg sync.WaitGroup
	success := false
	for {
		wg.Add(1)
		success = DLS(startURL, targetURL, i, result, numOfArticles, &wg)
		fmt.Println(i)
		i++
		wg.Wait()
		if success {
			return true
		}
		if i > maxDepth { // Safe condition only
			break
		}
	}
	return false
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
