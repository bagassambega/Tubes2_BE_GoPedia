package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"

	"sync"
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

type SafeBoolMap struct {
	sync.RWMutex
	SafeMap map[string]bool
}

type SafeStringMap struct {
	sync.RWMutex
	SafeMap map[string]string
}

func (rm *SafeBoolMap) Store(key string, value bool) {
	rm.Lock()
	rm.SafeMap[key] = value
	rm.Unlock()
}

func (rm *SafeBoolMap) Load(key string) (bool, bool){
	rm.RLock()
	result, ok := rm.SafeMap[key]
	rm.RUnlock()
	return result, ok
}

func (rm *SafeStringMap) Store(key string, value string) {
	rm.Lock()
	rm.SafeMap[key] = value
	rm.Unlock()
}

func (rm *SafeStringMap) Load(key string) (string, bool){
	rm.RLock()
	result,ok  := rm.SafeMap[key]
	rm.RUnlock()
	return result, ok
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
	visited := SafeBoolMap{SafeMap : make(map[string]bool)}
	history := SafeStringMap{SafeMap : make(map[string]string)}

	fmt.Print("Awal: ")
	fmt.Scan(&start)
	fmt.Print("Akhir: ")
	fmt.Scan(&goal)

	startTime := time.Now()
	queue = append(queue, start)

	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
		// colly.Async(true),
	)

	c.OnRequest(func(r *colly.Request) {
		urlVisited++
	})

	c.OnHTML("div#mw-content-text a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/wiki/") && !checkIgnoredLink(href) {
			kode := href[6:]
			if href == "/wiki/"+goal {
				found = true
				history.Store(kode, currLink)
				e.Request.Abort()
			} else {
				if _, exists := history.Load(kode); !exists {
					history.Store(kode, currLink)
					queue = append(queue, kode)
				}
				visited.Store(kode, false)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL.String())
		fmt.Println("Error:", err)
	})

	// limiter := make(chan int, 200)
	var wg sync.WaitGroup
	for !found {
		for _, element := range queue {
			wg.Add(1)
			go func(link string) {
				defer wg.Done()
				currLink = link
				isVisited, _ := visited.Load(currLink)
				if !isVisited {
					c.Visit("https://en.wikipedia.org/wiki/" + currLink)
				}
				queue = HapusAntrian(queue, &parent)
				visited.Store(parent, true)
				// <-limiter
			}(element)
			if found {
				break
			}
		}
	}

	end := time.Now()
	fmt.Println("Waktu eksekusi", end.Sub(startTime))
	fmt.Println("Url visited: ", urlVisited)
	if found {
		shortestPath = getResult(history.SafeMap, start, goal)
		for i := len(shortestPath) - 1; i >= 0; i-- {
			fmt.Println(shortestPath[i])
		}
	} else {
		fmt.Println("Goal not found")
	}
}