package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strings"
)

// Global variables
var linkCache = make(map[string][]string)
var visited = make(map[string]bool)

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
		return links, true
	}

	links = getAllLinks(url)

	linkCache[url] = links
	return links, false
}

func DLS(currentURL string, targetURL string, limit int, result []string, numOfArticles *int) ([]string, bool) {
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
		newPath, found := DLS(link, targetURL, limit-1, append(result, link), numOfArticles)
		if found {
			return newPath, true
		}
	}
	return nil, false
}

// IDSGoroutine Return path, number of articles, and whether the target is found
func IDS(startURL, targetURL string, maxDepth int, numOfArticles *int) ([]string, int, bool) {
	i := 1
	var result []string

	ch := make(chan []string, maxDepth)
	go func(ch chan []string) {
		for {
			result, success := DLS(startURL, targetURL, i, result, numOfArticles)
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
	linkCache = make(map[string][]string)
	visited = make(map[string]bool)
	return result, *numOfArticles, result != nil
}
