package BFS

import (
	"fmt"
	"strings"
	"time"
	"github.com/gocolly/colly"
)

type Pair struct {
    First  string
    Second bool
}

type node struct {
	childSize int
	children  []*node
	parent    *node
	value     Pair
}

func newNode(Value string, Parent *node) *node {
	return &node{value: Pair{First: Value, Second: false}, childSize: 0, parent: Parent}
}

func BFS(c *colly.Collector, start string, end string, linkTree *node) bool {
	find := true
	if (!linkTree.value.Second) {
		c.Visit("https://en.wikipedia.org/wiki/" + linkTree.value.First)
	}
	for i := 0; i < linkTree.childSize && find; i++ {
		tempTree := linkTree.children[i]
		if tempTree.value.First == end {
			linkTree = tempTree
			fmt.Println(linkTree.value)
			find = false
		}
	}
	return find
}

func main() {
	var start string
	var end string
	// visited := map[*node]bool{}

	linkTree := newNode("", nil)
	currentTree := linkTree

	var find bool = true
	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	fmt.Scan(&start)
	linkTree.value.First = start
	fmt.Scan(&end)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if href != "" && len(href) > 6 && href[:6] == "/wiki/" && href != "/wiki/Main_Page" && 
			href[:10] == "/wiki/File" && href != "/wiki/Special" &&
			strings.ToLower(href[(len(href)-4):]) != ".jpg" && strings.ToLower(href[(len(href)-4):]) != ".png" && 
			strings.ToLower(href[(len(href)-4):]) != ".jpeg" && strings.ToLower(href[(len(href)-4):]) != ".svg" {
			currentTree.children = append(currentTree.children, newNode(href[6:], currentTree))
			currentTree.childSize++
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL)
		currentTree.value.Second = true
	})

	for find {
		if currentTree.value.First == end {
			fmt.Println((currentTree.value))
			find = false
		} else {
			find = BFS(c, start, end, linkTree)
			for i := 0; i < linkTree.childSize && find; i++ {
				currentTree = linkTree.children[i]
				find = BFS(c, start, end, currentTree)
			}
		}
	}

	for currentTree.value.First != start {
		println(currentTree.value)
		currentTree = currentTree.parent
	}
	println(currentTree.value)
}
