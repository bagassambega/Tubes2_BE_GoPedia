package main

import (
	"fmt"
	"strings"
	"time"
	"github.com/gocolly/colly/v2"
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
	if (len(queue) <= 1) {
		queue = []string{}
	} else {
		*parent = queue[0]
	}
	return queue
}

func checkIgnoredLink(url string) bool {
	ignoredLinks := [...]string{"/File:", "/Special:", "/Template:", "/Template_page:", "/Help:", "/Category:", "Special:", "/Wikipedia:", "/Portal:", "/Talk:"}
	for _, st := range ignoredLinks {
		if strings.Contains(url, st) {
			return true
		}
	}
	return false
}

// func printString(queue []string) {
// 	for _, elm := range queue {
// 		fmt.Println(elm)
// 	}
// }

func main() {
	var queue []string
	var history map[string]string
	var start string
	var goal string
	var parent string
	urlVisited := 0
	found := false
	visited := make(map[string]bool)

	fmt.Print("Awal: ")
	fmt.Scan(&start)
	fmt.Print("Akhir: ")
	fmt.Scan(&goal)

	startTime := time.Now()
	BuatAntrian(&queue, start)
	history = make(map[string]string)
	// history[start] = " "
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL)
	})

	c.OnHTML("div#mw-content-text a[href]", func(e *colly.HTMLElement) {
		urlVisited++
		href := e.Attr("href")
		if strings.HasPrefix(href, "/wiki/") && !checkIgnoredLink(href) {
			// history[href] = parent
			if href == goal { 
				found = true
			} else {
				queue = append(queue, href[6:])
				history[href[6:]] = queue[0]
				visited[href[6:]] = false
			}
		}
	})

	c.Visit("https://en.wikipedia.org/wiki/" + start)
	queue = HapusAntrian(queue, &parent)
	// printString(queue)
	for !found {
		visited[parent] = true
		for _, currLink := range queue {
			if !visited[currLink] {
				c.Visit("https://en.wikipedia.org/wiki/" + currLink)
				queue = HapusAntrian(queue, &parent)
			}
		}
	}

	// key := goal
	fmt.Println(goal)
	// for key != start {
	// 	fmt.Println(history[key])
	// 	key = history[key]
	// }
	// fmt.Println(key)
	end := time.Now()
	fmt.Println("Waktu eksekusi", end.Sub(startTime))
	fmt.Println("Url visited: ", urlVisited)
}
