package main

import (
	"fmt"
	"strings"
	"time"
	"sync"
	"github.com/gocolly/colly/v2"
)

type TreeNode struct {
	childSize int
	children  []*TreeNode
	value     string
}

type Pair struct {
	First  string
	Second bool
}

func MasukAntrian(queue *[]*TreeNode, link string) {
	*queue = append(*queue, NewTreeNode(link))
}

func AntrianKosong(queue []*TreeNode) bool {
	return len(queue) == 0
}

func HapusAntrian(queue []*TreeNode, parent *string) []*TreeNode {
	if (len(queue) <= 1) {
		queue = []*TreeNode{}
	} else {
		queue = queue[1:]
		*parent = queue[0].value
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
	// var mutex sync.Mutex
	
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
				// mutex.Lock()
// 				history[href[6:]] = queue[0]
				// mutex.Unlock()
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

func NewTreeNode(value string) *TreeNode {
	return &TreeNode{
		value:    value,
		children: []*TreeNode{},
		childSize: 0,
	}
}

func main() {
	var start string
	var goal string
	var parent string
	urlVisited := 0
	found := false
	visited := make(map[string]bool)
	var mutex sync.Mutex

	fmt.Print("Awal: ")
	fmt.Scan(&start)
	fmt.Print("Akhir: ")
	fmt.Scan(&goal)
	
	startTime := time.Now()

	root := NewTreeNode(" ")
	root.children = append(root.children, NewTreeNode(start))
	queue := []*TreeNode{root}

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
			if href[6:] == goal {
				found = true
				// mutex.Lock()
				fmt.Println(href[6:])
				// history[href[6:]] = queue[0]
				// mutex.Unlock()
			} else {
				mutex.Lock()
				queue = append(queue, NewTreeNode(href[6:]))
				mutex.Unlock()
				// history[href[6:]] = queue[0]
				visited[href[6:]] = false
			}
			fmt.Println(len(queue))
		}
	})
	
	parents := make(map[*TreeNode]*TreeNode)
	
	c.Visit("https://en.wikipedia.org/wiki/" + start)
	queue = HapusAntrian(queue, &parent)
	// mutex.Lock()
	visited[parent] = true
	// mutex.Unlock()
	
	// limiter := make(chan int, 200)
	for !found {
		mutex.Lock()
		node := queue[0]
		queue = HapusAntrian(queue, &parent)
		mutex.Unlock()
		for _, Node := range node.children {
			currLink := Node.value
			parents[Node] = node
			if (found) {
				break
			}
			// limiter <- 1
			// go func(link string) {
				// mutex.Lock()
			if !visited[currLink] {
				// mutex.Unlock()
				c.Visit("https://en.wikipedia.org/wiki/" + currLink)
				fmt.Println(len(queue))
				mutex.Lock()
				queue = HapusAntrian(queue, &parent)
				visited[parent] = true
				mutex.Unlock()
			} else {
				// mutex.Unlock()
			}
				// <- limiter
			// }(currLink)
		}
	}

	if found {
		// key := goal
		// for key != start {
		// 	fmt.Println(key)
		// 	mutex.Lock()
		// // 	// key = history[key]
		// 	mutex.Unlock()
		// }
		// fmt.Println(key)
	} else {
		fmt.Println("Goal not found")
	}
	end := time.Now()
	fmt.Println("Waktu eksekusi", end.Sub(startTime))
	fmt.Println("Url visited: ", urlVisited)
}
