package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strings"
)

type Map struct {
	parent string
	child  string
}

// Jika mengandung salah satu dari identifier seperti File: pada URL, return true
func checkIgnoredLink(url string) bool {
	ignoredLinks := [...]string{"/File:", "/Special:", "/Template:", "/Template_page:", "/Help:", "/Category:", "Special:", "/Wikipedia:", "/Portal:", "/Talk:"}
	for _, st := range ignoredLinks {
		if strings.Contains(url, st) {
			return true
		}
	}
	return false
}

func getAllLinks(url string) []string {
	c := colly.NewCollector()

	// Inisialisasi array
	var links []string

	// Cari semua link dan kalau berawalan /wiki/ ditambahkan, dan jika ada yang mengandung ignoredLinks diabaikan
	c.OnHTML("div#mw-content-text a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasPrefix(link, "/wiki/") && !checkIgnoredLink(link) {
			links = append(links, "https://en.wikipedia.org"+link)
		}
	})

	err := c.Visit(url)
	if err != nil {
		return nil
	}
	return links
}

func DLS(currentURL string, targetURL string, limit int) bool {
	if currentURL == targetURL {
		return true
	}

	if limit < 0 {
		return false
	}

	links := getAllLinks(currentURL)

	for _, link := range links {
		fmt.Println("Cek link", link, "di level", limit)
		if DLS(link, targetURL, limit-1) {
			return true
		}
	}
	return false
}

func IDS(startURL string, targetURL string, maxDepth int) bool {
	for i := 0; i <= maxDepth; i++ {
		if DLS(startURL, targetURL, i) {
			return true
		}
	}
	return false
}

func main() {
	startURL := "https://en.wikipedia.org/wiki/Russia"
	targetURL := "https://en.wikipedia.org/wiki/Indonesia"
	i := 1
	for {
		if IDS(startURL, targetURL, i) {
			fmt.Println("Berhasil")
			break
		} else {
			fmt.Println("Belum ada di level", i)
			i++
		}
	}

}
