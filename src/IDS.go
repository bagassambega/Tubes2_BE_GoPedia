package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strings"
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

func cacheLinks(url string) ([]string, bool) {
	links, exists := linkCache[url]

	if exists {
		//fmt.Println("Menggunakan cache untuk", url)
		return links, true
	}

	links = getAllLinks(url)

	linkCache[url] = links
	return links, false
}

func convertToTitleCase(s string) string {
	return strings.ReplaceAll(s, " ", "_")
}

func DLS(currentURL string, targetURL string, limit int, result []string, visited map[string]bool, numOfArticles *int) ([]string, bool) {
	if currentURL == targetURL {
		return result, true
	}

	if limit <= 1 || visited[currentURL] {
		return nil, false
	}

	visited[currentURL] = true
	defer delete(visited, currentURL)
	links, cached := cacheLinks(currentURL)
	if !cached {
		*numOfArticles++
	}

	for _, link := range links {
		//fmt.Println("Cek link", link)
		newPath, found := DLS(link, targetURL, limit-1, append(result, link), visited, numOfArticles)
		if found {
			return newPath, true
		}
	}
	return nil, false
}

// IDSGoroutine Return path, number of articles, and whether the target is found
func IDSGoroutine(startURL, targetURL string, maxDepth int, numOfArticles *int) ([]string, int, bool) {
	i := 1
	var result []string
	var visited = make(map[string]bool)

	ch := make(chan []string, maxDepth)
	go func(ch chan []string) {
		for {
			result, success := DLS(startURL, targetURL, i, result, visited, numOfArticles)
			if success {
				ch <- result
				close(ch)
				return
			}
			fmt.Println(i)
			i++
			if i > maxDepth { // Safe condition only
				ch <- nil
				close(ch)
				return
			}
		}
	}(ch)
	result = <-ch
	return result, *numOfArticles, result != nil
}

//func IDS(startURL string, targetURL string, maxDepth int, numOfArticles *int) ([]string, bool) {
//	i := 1
//	var result []string
//	var visited = make(map[string]bool)
//	for {
//		result, success := DLS(startURL, targetURL, i, result, visited, numOfArticles)
//		if success {
//			return result, true
//		}
//		fmt.Println(i)
//		i++
//		if i > maxDepth { // Safe condition only
//			break
//		}
//	}
//	return nil, false
//}
