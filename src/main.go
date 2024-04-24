package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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
			startURL = "https://en.wikipedia.org/wiki/" + source
			targetURL = "https://en.wikipedia.org/wiki/" + target

			// Panggil IDS
			hasil, found := IDS(startURL, targetURL, maxDepth, &numOfArticles)
			result := []string{startURL}
			result = append(result, hasil...)

			if found {
				end := time.Now()
				elapsedTime = end.Sub(startTime)
				fmt.Println("Waktu eksekusi", elapsedTime)
			}

			// Tampilkan hasil
			c.JSON(http.StatusOK, gin.H{
				"numOfArticles": numOfArticles,
				"result":        result,
				"length":        len(result),
				"elapsedTime":   elapsedTime.String(),
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
