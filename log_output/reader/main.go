package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func readStatus(path string) (timestamp string, id string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}

	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	last := lines[len(lines)-1]

	parts := strings.SplitN(last, ": ", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected line format: %q", last)
	}

	return parts[0], parts[1], nil
}

func main() {
	path := os.Getenv("LOG_FILE_PATH")
	if path == "" {
		path = "/data/status.log"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5607"
	}

	r := gin.Default()
	r.GET("/status", func(c *gin.Context) {
		timestamp, id, err := readStatus(path)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"timestamp": timestamp,
			"string":    id,
		})
	})

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
