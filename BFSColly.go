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

func scrape(currLink string, found* bool, goal* string, queue *[]string, visited map[string]bool, history map[string]string) {
	c := colly.NewCollector()
	var mutexhis sync.Mutex
	var mutexvis sync.Mutex

	c.OnHTML("div#mw-content-text a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		mutexvis.Lock()
		if strings.HasPrefix(href, "/wiki/") && !checkIgnoredLink(href) && !(visited[href[6:]]){
			if href == "/wiki/"+ *goal {
				*found = true
				mutexhis.Lock()
				history[href[6:]] = currLink
				mutexhis.Unlock()
				e.Request.Abort()
				} else {
					*queue = append(*queue, href[6:])
					mutexhis.Lock()
					history[href[6:]] = currLink
					mutexhis.Unlock()
					mutexvis.Lock()
					visited[href[6:]] = false
					mutexvis.Unlock()
			}
		}
		mutexvis.Unlock()
	})
	
	c.Visit("https://en.wikipedia.org/wiki/" + currLink)
}

func main() {
	var currLink string
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

	scrape(start, &found, &goal, &queue, visited, history)
	queue = HapusAntrian(queue, &parent)
	
	limiter := make(chan int, 150)
	for !found {
		// mutex.Lock()
		
		// fmt.Println(queue[0])
		visited[parent] = true
		// mutex.Unlock()
		for _, element := range queue {
			limiter <- 1
			go func(link string) {
				currLink = element
				// defer func() {
			// 	<-limiter // Release the limiter token
			// }()
			mutex.Lock()
			if !visited[currLink] {
					mutex.Unlock()
					scrape(currLink, &found, &goal, &queue, visited, history)
					// c.Visit("https://en.wikipedia.org/wiki/" + currLink)
					queue = HapusAntrian(queue, &parent)
				} else {
					mutex.Unlock()
				}
				<-limiter
			}(currLink)
			if found {
				break
			}
			// wg.Wait()
		}
	}

	if found {
		// key := goal
		// fmt.Println(history["10th_edition_of_Systema_Naturae"])
		// for key != start {
		// mutex.Lock()
		// mutex.Unlock()
		// mutex.Lock()
		// 	key = history[key]
		// mutex.Unlock()
		// }
		// fmt.Println(key)
	} else {
		fmt.Println("Goal not found")
	}
	// fmt.Println(history)
	fmt.Println(history[goal])
	end := time.Now()
	fmt.Println("Waktu eksekusi", end.Sub(startTime))
	fmt.Println("Url visited: ", urlVisited)
}
