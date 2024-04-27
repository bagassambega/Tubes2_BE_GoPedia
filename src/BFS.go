package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"

	"sync"
)

func HapusAntrian(queue []string, parent *string) []string {
	if len(queue) <= 1 {
		return []string{}
	} else {
		*parent = queue[0]
	}
	return queue[1:]
}

func getResult(history map[string]string, start string, goal string) []string {
	var result []string
	key := start
	for key != goal {
		result = append(result, "//en.wikipedia.org/wiki/" + key)
		key = history[key]
	}
	result = append(result, "//en.wikipedia.org/wiki/" + goal)
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

func (rm *SafeBoolMap) Load(key string) (bool, bool) {
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

func (rm *SafeStringMap) Load(key string) (string, bool) {
	rm.RLock()
	result,ok  := rm.SafeMap[key]
	rm.RUnlock()
	return result, ok
}

func scrape (currLink string, visited *SafeBoolMap, history *SafeStringMap, urlVisited *int, goal string, found *bool) ([]string){
	tempQueue := []string{}
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnRequest(func(r *colly.Request) {
		*urlVisited++
	})

	c.OnHTML("div#mw-content-text a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/wiki/") && !checkIgnoredLink(href) {
			kode := href[6:]
			if href == "/wiki/" + goal {
				*found = true
				history.Store(kode, currLink)
				return
			} else {
				if _, exists := history.Load(kode); !exists {
					history.Store(kode, currLink)
				}
				tempQueue = append(tempQueue, kode)
				visited.Store(kode, false)
			}
		}
	})
	
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL.String())
		fmt.Println("Error:", err)
	})
	
	c.Visit("https://en.wikipedia.org/wiki/" + currLink)

	return tempQueue
}

// func BFS(start string, goal string, urlVisited *int) ([]string, bool) {
// 	var shortestPath []string
// 	var tempQueue []string
// 	var queue []string
// 	var parent string
// 	found := false
// 	visited := SafeBoolMap{SafeMap : make(map[string]bool)}
// 	history := SafeStringMap{SafeMap : make(map[string]string)}

// 	startTime := time.Now()

// 	tempQueueChan := make(chan []string)

// 	go func() {
// 		tempQueue := scrape(start, &visited, &history, urlVisited, goal, &found)
// 		tempQueueChan <- tempQueue
// 	}()
// 	tempQueue = <-tempQueueChan
// 	visited.Store(parent, true)
	
// 	limiter := make(chan int, 100)
// 	var wg sync.WaitGroup
// 	for !found {
// 		queue = []string{}
// 		queue = append(queue, tempQueue...)
// 		tempQueue = []string{}
// 		for _, element := range queue{
// 			wg.Add(1)
// 			limiter <- 1
// 			go func(link string) {
// 				defer wg.Done()
// 				if isVisited, _ := visited.Load(link); !isVisited {
// 					tempQueue = append(tempQueue, scrape(link, &visited, &history, urlVisited, goal, &found)...)
// 					visited.Store(parent, true)
// 				}
// 				<-limiter
// 			}(element)
// 			if (found) {
// 				break
// 			}
// 		}
// 		wg.Wait()
// 	}

// 	end := time.Now()
// 	fmt.Println("Waktu eksekusi:", end.Sub(startTime))
// 	fmt.Println("Url visited:", urlVisited)
// 	if (found) {
// 		fmt.Println(goal)
// 		shortestPath = getResult(history.SafeMap, goal, start)
// 		fmt.Println(shortestPath)
// 	} else {
// 		fmt.Println("Goal not found")
// 	}
// 	return shortestPath, found
// }

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

	// root := NewTreeNode(" ")
	queue = append(queue, start)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println(r.URL)
		*urlVisited++
	})

	fmt.Println("Start BFS")

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
				if _,exists := history[kode]; !exists {
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
			if found {
				break
			}

		}
		queue = HapusAntrian(queue, &parent)
	}

	if found {
		fmt.Println("Goal found")
		shortestPath = getResult(history, goal, start)
	} else {
		fmt.Println("Goal not found")
	}

	return shortestPath, found
}