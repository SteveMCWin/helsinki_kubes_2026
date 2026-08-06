package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	counter := 0

	r := gin.Default()
	r.GET("/pingpong", func(c *gin.Context) {
		current := counter
		counter++

		c.String(200, "pong %d", current)
	})
	r.GET("/pings", func(c *gin.Context) {
		c.String(200, "%d", counter)
	})

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
