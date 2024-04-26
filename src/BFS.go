package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strings"
	// "sync"
	"time"
)

type Pair struct {
	First  string
	Second bool
}

func BuatAntrian(queue *[]string, start string) {
	*queue = append(*queue, start)
}

func MasukAntrian(queue *[]string, link string) {
	*queue = append(*queue, link)
}

func AntrianKosong(queue []string) bool {
	return len(queue) == 0
}

func HapusAntrian(queue []string, parent *string) []string {
	if len(queue) <= 1 {
		queue = []string{}
	} else {
		*parent = queue[0]
	}
	return queue
}

func getResult(history map[string]string, start string, goal string) []string {
	var result []string
	key := goal
	for key != start {
		result = append(result, key)
		key = history[key]
	}
	result = append(result, start)
	return result
}

func BFS(start string, goal string, urlVisited *int) ([]string, bool) {
	var shortestPath []string
	// var allPath [][]string
	var currLink string
	var queue []string
	var parent string
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
		*urlVisited++
	})

	c.OnHTML("div#mw-content-text a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/wiki/") && !checkIgnoredLink(href){
			kode := href[6:]
			if href == "/wiki/"+goal {
				found = true
				history[kode] = currLink
				e.Request.Abort()
			} else {
				queue = append(queue, kode)
				if (!visited[kode]) {
					history[kode] = currLink
				}
				visited[kode] = false
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
			if (found) {
				break
			}

		}
		queue = HapusAntrian(queue, &parent)
	}

	if found {
		shortestPath = getResult(history, goal, start)
	} else {
		fmt.Println("Goal not found")
	}
	end := time.Now()
	fmt.Println("Waktu eksekusi", end.Sub(startTime))
	fmt.Println("Url visited: ", urlVisited)
	return shortestPath, found
}
