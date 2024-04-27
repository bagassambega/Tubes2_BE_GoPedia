package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strings"
	"sync"
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

// IDS Return path, number of articles, and whether the target is found
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

type Visited struct {
	visited2 map[string]bool
	sync.RWMutex
}

type Cache struct {
	cache map[string][]string
	sync.RWMutex
}

func (v *Visited) setVisited(url string) {
	v.Lock()
	defer v.Unlock()
	v.visited2[url] = true
}

func (v *Visited) isVisited(url string) bool {
	v.RLock()
	defer v.RUnlock()
	return v.visited2[url]
}

func (c *Cache) setCache(url string, links []string) {
	c.Lock()
	defer c.Unlock()
	c.cache[url] = links
}

func (c *Cache) getCache(url string) ([]string, bool) {
	c.RLock()
	defer c.RUnlock()
	links, exists := c.cache[url]
	return links, exists
}

func (c *Cache) cacheLinks(url string) ([]string, bool) {
	links, exists := c.getCache(url)

	if exists {
		return links, true
	}

	links = getAllLinks(url)

	c.setCache(url, links)
	return links, false
}

func (v *Visited) deleteVisited(url string) {
	v.Lock()
	defer v.Unlock()
	delete(v.visited2, url)
}

func DLS2(currentURL string, targetURL string, limit int, result []string, numOfArticles *int, visited2 *Visited, cache *Cache) ([]string, bool) {
	if currentURL == targetURL {
		return result, true
	}

	if limit <= 1 || visited2.isVisited(currentURL) {
		return nil, false
	}

	visited2.setVisited(currentURL)
	defer visited2.deleteVisited(currentURL)
	links, cached := cache.cacheLinks(currentURL)
	if !cached {
		*numOfArticles++
	}

	wg := sync.WaitGroup{}
	//fmt.Println("Panjang", len(links))
	limiter := make(chan int, 200)

	for _, link := range links {
		wg.Add(1)
		found := false
		limiter <- 1
		var newPath []string
		go func(link string, found *bool, newPath *[]string) {
			//fmt.Print("Cek link", link, "    ")
			defer func() {
				wg.Done()
				<-limiter
			}()
			*newPath, *found = DLS2(link, targetURL, limit-1, append(result, link), numOfArticles, visited2, cache)
		}(link, &found, &newPath)
		wg.Wait()
		if found {
			return newPath, true
		}

		//newPath, found := DLS2(link, targetURL, limit-1, append(result, link), numOfArticles, visited2, cache)
		//if found {
		//	return newPath, true
		//}
	}
	wg.Wait()
	return nil, false
}

func IDS2(startURL, targetURL string) ([]string, int, bool) {
	maxDepth := 10
	i := 1
	var result []string
	visited2 := Visited{visited2: make(map[string]bool)}
	cache := Cache{cache: make(map[string][]string)}
	numOfArticles := 0

	for {
		result, success := DLS2(startURL, targetURL, i, result, &numOfArticles, &visited2, &cache)
		if success {
			return result, numOfArticles, true
		}
		fmt.Println(i)
		i++
		if i > maxDepth { // Safe condition only
			break
		}
	}
	return nil, numOfArticles, false
}
