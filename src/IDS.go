package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"strings"
	"sync"
	"time"
)

// Global variables
var linkCache = make(map[string][]string)
var cacheMutex = sync.RWMutex{}

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
	cacheMutex.RLock()
	links, exists := linkCache[url]
	cacheMutex.RUnlock()

	if exists {
		//fmt.Println("Menggunakan cache untuk", url)
		return links
	}

	links = getAllLinks(url)
	cacheMutex.Lock()
	linkCache[url] = links
	cacheMutex.Unlock()
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
	*result = []string{}
	for i := 0; i <= maxDepth; i++ {
		if DLS(startURL, targetURL, i, result, numOfArticles) {
			return true
		}
	}
	return false
}

func main() {
	start := time.Now()
	startURL := "https://en.wikipedia.org/wiki/Russia"
	targetURL := "https://en.wikipedia.org/wiki/Joko_Widodo"
	i := 1
	numOfArticles := 0
	fmt.Println("Start URL", startURL)
	result := make([]string, 0)
	for {
		result = []string{}
		if IDS(startURL, targetURL, i, &result, &numOfArticles) {
			fmt.Println("Berhasil dengan array", len(result))
			for _, r := range result {
				fmt.Println(r)
			}
			fmt.Println("Jumlah artikel yang dikunjungi", numOfArticles)
			break
		} else {
			fmt.Println("Belum ada di level", i)
			i++
		}
	}
	end := time.Now()
	fmt.Println("Waktu eksekusi", end.Sub(start))

	router := gin.Default()
	router.GET("/IDS", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"numOfArticles": numOfArticles,
			"result":        result,
			"length":        len(result),
		})
	})
	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
