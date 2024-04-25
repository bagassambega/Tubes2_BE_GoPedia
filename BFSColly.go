package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strings"
	// "sync"
	"time"
)

type TreeNode struct {
	Value string
	Children []*TreeNode
}

type Pair struct {
	First  string
	Second bool
}

func NewTreeNode(value string) *TreeNode {
	return &TreeNode{Value: value}
}

func (node *TreeNode) AddChild(child *TreeNode) {
	node.Children = append(node.Children, child)
}

func BuatAntrian(queue *[]*TreeNode, start string) {
	*queue = append(*queue, NewTreeNode(start))
}

func MasukAntrian(queue *[]*TreeNode, link string) {
	*queue = append(*queue, NewTreeNode(link))
}

func AntrianKosong(queue []*TreeNode) bool {
	return len(queue) == 0
}

func HapusAntrian(queue []*TreeNode, parent *string) []*TreeNode {
	if len(queue) <= 1 {
		return []*TreeNode{}
	} else {
		*parent = queue[0].Value
	}
	return queue [1:]
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
	var queue []*TreeNode
	var parent string
	urlVisited := 0
	found := false
	visited := make(map[string]bool)
	history := make(map[*TreeNode]*TreeNode)
	// var mutex sync.Mutex

	fmt.Print("Awal: ")
	fmt.Scan(&start)
	fmt.Print("Akhir: ")
	fmt.Scan(&goal)

	startTime := time.Now()
	root := NewTreeNode(" ")
	root.AddChild(NewTreeNode(start))
	queue = append(queue, root) //Masuk Antrian

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL)
		urlVisited++
	})

	c.OnHTML("div#mw-content-text a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.HasPrefix(href, "/wiki/") && !checkIgnoredLink(href) {
			if href == "/wiki/"+goal {
				found = true
				history[NewTreeNode(href[6:])] = queue[0]
				e.Request.Abort()
				// mutex.Lock()
				// mutex.Unlock()
			} else {
				queue = append(queue, NewTreeNode(href[6:]))
				// mutex.Lock()
				history[queue[len(queue) - 1]] = queue[0]
				// fmt.Println(len(queue), "asd")
				visited[href[6:]] = false
				// mutex.Unlock()
			}
		}
	})
	
	// c.Visit("https://en.wikipedia.org/wiki/" + start)
	// queue = HapusAntrian(queue, &parent)
	
	// limiter := make(chan int, 200)
	for !found {
		fmt.Println(len(queue))
		// fmt.Println(queue[0].Value)
		Node := queue[0]
		queue = HapusAntrian(queue, &parent)
		// mutex.Lock()
		visited[parent] = true
		// mutex.Unlock()
		for _, TreeNode := range Node.Children {
			// limiter <- 1
			currLink := TreeNode.Value
			// go func(link string) {
				// defer func() {
				// 	<-limiter // Release the limiter token
				// }()
				// mutex.Lock()
			if !visited[currLink] {
				// mutex.Unlock()
				// fmt.Println("yoo")
				// mutex.Lock()
				c.Visit("https://en.wikipedia.org/wiki/" + currLink)
				// mutex.Unlock()
				queue = HapusAntrian(queue, &parent)
			}
			// } else {
			// 	// mutex.Unlock()
			// }
				// <-limiter
			// }(currLink)
			if (found) {
				break
			}
			// wg.Wait()
		}
	}

	if found {
		// key := goal
		// for key != start {
			// mutex.Lock()
		// 	fmt.Println((*history)[key])
			// mutex.Unlock()
			// mutex.Lock()
		// 	key = (*history)[key]
			// mutex.Unlock()
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
