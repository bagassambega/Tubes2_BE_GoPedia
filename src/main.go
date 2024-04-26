package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// Format URL: http://localhost:8080/gopedia?method=BFS&source=source&target=target

	router.Use(cors.Default()) //Supaya bisa diakses front-end
	router.GET("gopedia/", func(c *gin.Context) {
		// Ambil metode dulu apakah BFS atau IDS
		method := c.Query("method")
		source := c.Query("source")
		target := c.Query("target")

		// IDS
		if method == "IDS" {
			// Set maximum depth ke 9
			maxDepth := 9
			startTime := time.Now()
			numOfArticles := 0
			var elapsedTime time.Duration
			fmt.Println("Source", source, "Target", target)
			var startURL, targetURL string
			startURL = "https://en.wikipedia.org/wiki/" + convertToTitleCase(source)
			targetURL = "https://en.wikipedia.org/wiki/" + convertToTitleCase(target)

			// Panggil IDS
			hasil, numOfArticles, found := IDSGoroutine(startURL, targetURL, maxDepth, &numOfArticles)
			//hasil, found := IDS(startURL, targetURL, maxDepth, &numOfArticles)
			result := []string{startURL}
			result = append(result, hasil...)

			if found {
				end := time.Now()
				elapsedTime = end.Sub(startTime)
				fmt.Println("Waktu eksekusi", elapsedTime)
			}

			// Dapatkan judul artikel dari link
			//for i, link := range result {
			//	result[i] = convertToArticleTitle(link)
			//}

			// Tampilkan hasil
			c.JSON(http.StatusOK, gin.H{
				"numOfArticles": numOfArticles,
				"result":        result,
				"length":        len(result),
				"elapsedTime":   elapsedTime.String(),
			})

		} else { // BFS
			var numOfArticles int

			var elapsedTime time.Duration
			start := convertToTitleCase(source)
			goal := convertToTitleCase(target)

			startTime := time.Now()
			result, found := BFS(start, goal, &numOfArticles)

			if found {
				end := time.Now()
				elapsedTime = end.Sub(startTime)
				fmt.Println("Waktu eksekusi", elapsedTime)

				// reverse result
				length := len(result)
				for i := 0; i < length/2; i++ {
					// Swap elements from both ends
					result[i], result[length-1-i] = result[length-1-i], result[i]
				}
			}

			for i, link := range result {
				result[i] = convertToTitleCase(link)
			}

			c.JSON(http.StatusOK, gin.H{
				"numOfArticles": numOfArticles,
				"result":        result,
				"length":        len(result),
				"elapsedTime":   elapsedTime.String(),
			})
		}
	})

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
