package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strings"
	"sync"
)

// Global variables
var linkCache = make(map[string][]string)
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
