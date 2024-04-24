package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strings"
	"sync"
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

func checkIgnoredLink(url string) bool {
	ignoredLinks := [...]string{"/File:","/Main_Page", "/Special:", "/Template:", "/Template_page:", "/Help:", "/Category:", "Special:", "/Wikipedia:", "/Portal:", "/Talk:"}
	for _, st := range ignoredLinks {
		if strings.Contains(url, st) {
			return true
		}
	}
	return false
}

// func BFSrun(start string, goal string, wg *sync.WaitGroup) map[string]string {
// 	var queue []string
// 	var history map[string]string
// 	var parent string
// 	urlVisited := 0
// 	found := false
// 	var mutex sync.Mutex

// 	BuatAntrian(&queue, start)
// 	visited := make(map[string]bool)
// 	history = make(map[string]string)

// 	c := colly.NewCollector(
// 		colly.AllowedDomains("en.wikipedia.org"),
// 	)

// 	c.OnRequest(func(r *colly.Request) {
// 		// fmt.Println(r.URL)
// 	})

// 	c.OnHTML("div#mw-content-text a[href]", func(e *colly.HTMLElement) {
// 		urlVisited++
// 		href := e.Attr("href")
// 		if strings.HasPrefix(href, "/wiki/") && !checkIgnoredLink(href) {
// 			// history[href] = parent
// 			if href == goal {
// 				found = true
// 			} else {
// 				queue = append(queue, href[6:])
// 				mutex.Lock()
// 				history[href[6:]] = queue[0]
// 				mutex.Unlock()
// 				visited[href[6:]] = false
// 			}
// 		}
// 	})

// 	// printString(queue)
// 	c.Visit("https://en.wikipedia.org/wiki/" + start)
// 	queue = HapusAntrian(queue, &parent)

// 	wg.Add(1)
// 	for !found {
// 		visited[parent] = true
// 		for _, currLink := range queue {
// 			if !visited[currLink] {
// 				c.Visit("https://en.wikipedia.org/wiki/" + currLink)
// 				queue = HapusAntrian(queue, &parent)
// 			}
// 		}
// 	}

// 	return history
// }

func getResult(history map[string]string, start string, goal string) []string {
	var result []string
	key := goal
	for key != start {
		result = append(result, key)
		fmt.Println(history[key])
		key = history[key]
	}
	result = append(result, start)
	fmt.Println(start)
	return result
}

func main() {
	var start string
	var goal string
	var queue []string
	var parent string
	urlVisited := 0
	found := false
	visited := make(map[string]bool)
	history := make(map[string]string)
	var mutex sync.Mutex

	fmt.Print("Awal: ")
	fmt.Scan(&start)
	fmt.Print("Akhir: ")
	fmt.Scan(&goal)

	startTime := time.Now()
	BuatAntrian(&queue, start)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println(r.URL)
		urlVisited++
	})

	c.OnHTML("div#mw-content-text a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/wiki/") && !checkIgnoredLink(href) {
			if href == "/wiki/"+goal {
				found = true
				e.Request.Abort()
				mutex.Lock()
				history[href[6:]] = queue[0]
				mutex.Unlock()
			} else {
				queue = append(queue, href[6:])
				mutex.Lock()
				history[href[6:]] = queue[0]
				visited[href[6:]] = false
				mutex.Unlock()
			}
		}
	})

	c.Visit("https://en.wikipedia.org/wiki/" + start)
	queue = HapusAntrian(queue, &parent)

	limiter := make(chan int, 200)
	for !found {
		fmt.Println(len(queue))
		// mutex.Lock()
		
		// fmt.Println(queue[0])
		visited[parent] = true
		// mutex.Unlock()
		for _, currLink := range queue {
			limiter <- 1
			go func(link string) {
				// defer func() {
				// 	<-limiter // Release the limiter token
				// }()
				mutex.Lock()
				if !visited[currLink] {
					mutex.Unlock()
					c.Visit("https://en.wikipedia.org/wiki/" + currLink)
					queue = HapusAntrian(queue, &parent)
				} else {
					mutex.Unlock()
				}
				<-limiter
			}(currLink)
			if (found) {
				break
			}
			// wg.Wait()
		}
	}

	if found {
		// key := goal
		// for key != start {
		// 	mutex.Lock()
		// 	fmt.Println((*history)[key])
		// 	mutex.Unlock()
		// 	mutex.Lock()
		// 	key = (*history)[key]
		// 	mutex.Unlock()
		// }
		// fmt.Println(key)
	} else {
		fmt.Println("Goal not found")
	}
	end := time.Now()
	fmt.Println("Waktu eksekusi", end.Sub(startTime))
	fmt.Println("Url visited: ", urlVisited)
	// fmt.Println(history)
}
