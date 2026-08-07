package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

func readPingpongCountHTTP(url string) (int, error) {
	resp, err := http.Get(url)

	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return -1, fmt.Errorf("unexpected pingpong count format: %q", data)
	}

	return count, nil
}

func main() {
	path := os.Getenv("LOG_FILE_PATH")
	if path == "" {
		path = "/data/status.log"
	}

	pingpongUrl := "http://pingpong-svc.exercises:3456/pings"

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

		pingpongs, err := readPingpongCountHTTP(pingpongUrl)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"timestamp": timestamp,
			"string":    id,
			"pingpongs": pingpongs,
		})
	})

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
