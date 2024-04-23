package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()
	// Format URL: http://localhost:8080/gopedia/?method=BFS&source=source&target=target
	router.GET("gopedia/", func(c *gin.Context) {
		// Ambil metode dulu apakah BFS atau IDS
		method := c.Query("method")
		source := c.Query("source")
		target := c.Query("target")

		// IDS
		if method == "IDS" {
			// Set maximum depth ke 9
			maxDepth := 9
			var result []string
			startTime := time.Now()
			numOfArticles := 0
			var elapsedTime time.Duration
			fmt.Println("Source", source, "Target", target)
			var startURL, targetURL string
			startURL = "https://en.wikipedia.org/wiki/" + source
			targetURL = "https://en.wikipedia.org/wiki/" + target

			// Panggil IDS
			if IDS(startURL, targetURL, maxDepth, &result, &numOfArticles) {
				end := time.Now()
				elapsedTime = end.Sub(startTime)
				fmt.Println("Waktu eksekusi", elapsedTime)
			}

			// Tampilkan hasil
			c.JSON(http.StatusOK, gin.H{
				"numOfArticles": numOfArticles,
				"result":        result,
				"length":        len(result),
				"elapsedTime":   elapsedTime,
			})

		} else { // BFS
			history := make(map[string]string)
			var elapsed time.Duration
			BFS(source, target, &history, &elapsed)
			c.JSON(http.StatusOK, gin.H{
				"elapsedTime": elapsed,
				"result":      history,
			})
		}
	})

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
