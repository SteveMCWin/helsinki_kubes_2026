package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	filePath := os.Getenv("PINGPONG_FILE_PATH")
	if filePath == "" {
		filePath = "/data/pingpong_count.txt"
	}

	var mu sync.Mutex
	counter := 0

	r := gin.Default()
	r.GET("/pingpong", func(c *gin.Context) {
		mu.Lock()
		current := counter
		counter++
		err := os.WriteFile(filePath, []byte(strconv.Itoa(current)), 0644)
		mu.Unlock()
		if err != nil {
			fmt.Printf("failed to persist pingpong count: %v\n", err)
		}

		c.String(200, "pong %d", current)
	})

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
