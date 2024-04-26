package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"

	// "sync"
	"time"
)

type Pair struct {
	First  string
	Second bool
}

type NodeHistory struct {
	Link   string
	Parent *NodeHistory
}

func MasukAntrian(queue *[]string, start string) {
	*queue = append(*queue, start)
}

func AntrianKosong(queue []string) bool {
	return len(queue) == 0
}

func HapusAntrian(queue []string, parent *string) []string {
	if len(queue) <= 1 {
		return []string{}
	} else {
		*parent = queue[0]
	}
	return queue[1:]
}

func checkIgnoredLink(url string) bool {
	ignoredLinks := [...]string{"/File:", "/Main_Page", "/Special:", "/Template:", "/Template_page:", "/Help:", "/Category:", "Special:", "/Wikipedia:", "/Portal:", "/Talk:"}
	for _, st := range ignoredLinks {
		if strings.Contains(url, st) {
			return true
		}
	}
	return false
}

func getResult(history map[string]string, start string, goal string) []string {
	var result []string
	key := start
	for key != goal {
		result = append(result, key)
		key = history[key]
	}
	result = append(result, goal)
	return result
}

func main() {
	var start string
	var shortestPath []string
	var currLink string
	var goal string
	var queue []string
	var parent string
	urlVisited := 0
	found := false
	visited := make(map[string]bool)
	history := make(map[string]string)
	// var mutex sync.Mutex

	fmt.Print("Awal: ")
	fmt.Scan(&start)
	fmt.Print("Akhir: ")
	fmt.Scan(&goal)

	startTime := time.Now()
	// root := NewTreeNode(" ")
	queue = append(queue, start)

	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println(r.URL)
		urlVisited++
	})

	c.OnHTML("div#mw-content-text a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/wiki/") && !checkIgnoredLink(href) {
			kode := href[6:]
			if href == "/wiki/"+goal {
				found = true
				history[kode] = currLink
				e.Request.Abort()
			} else {
				queue = append(queue, kode)
				// mutex.Lock()
				if _,exists := history[kode]; !exists {
					history[kode] = currLink
				}
				// mutex.Unlock()
				// mutex.Lock()
				visited[kode] = false
				// mutex.Unlock()
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL.String())
		fmt.Println("Error:", err)
	})
	// limiter := make(chan int, 200)
	for !found {
		// mutex.Lock()
		visited[parent] = true
		// mutex.Unlock()

		for _, element := range queue {
			// limiter <- 1
			// go func(link string) {
			currLink = element
			// mutex.Lock()
			if !visited[currLink] {
				// mutex.Unlock()
				c.Visit("https://en.wikipedia.org/wiki/" + currLink)
				queue = HapusAntrian(queue, &parent)
			}
			// <-limiter
			// }(currLink)
			if found {
				break
			}

		}
		queue = HapusAntrian(queue, &parent)
	}
	
	end := time.Now()
	fmt.Println("Waktu eksekusi", end.Sub(startTime))
	fmt.Println("Url visited: ", urlVisited)
	if found {
		fmt.Println(history[goal])
		shortestPath = getResult(history, goal, start)
		for _, X := range shortestPath {
			fmt.Println(X)
		}
	} else {
		fmt.Println("Goal not found")
	}
}
