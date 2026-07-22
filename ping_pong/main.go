package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var mu sync.Mutex
	counter := 0

	r := gin.Default()
	r.GET("/pingpong", func(c *gin.Context) {
		mu.Lock()
		current := counter
		counter++
		mu.Unlock()

		c.String(200, "pong %d", current)
	})

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
