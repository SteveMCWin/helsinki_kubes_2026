package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func fetchImage(path string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		// port = "8080"
		panic("PORT is not set")
	}

	imagePath := os.Getenv("IMAGE_FILE_PATH")
	if imagePath == "" {
		// imagePath = "/data/image.jpg"
		panic("IMAGE_FILE_PATH is not set")
	}

	todoBackendUrl := os.Getenv("TODO_BACKEND_URL")
	if todoBackendUrl == "" {
		// todoBackendUrl = "http://todoapp-svc:1235/todos"
		panic("TODO_BACKEND_URL is not set")
	}

	picsumUrl := os.Getenv("PICSUM_URL")
	if picsumUrl == "" {
		panic("PICSUM_URL is not set")
	}

	refreshMinutes := os.Getenv("IMAGE_REFRESH_MINUTES")
	if refreshMinutes == "" {
		panic("IMAGE_REFRESH_MINUTES is not set")
	}
	refreshInterval, err := strconv.Atoi(refreshMinutes)
	if err != nil {
		panic(fmt.Sprintf("IMAGE_REFRESH_MINUTES is not a number: %q", refreshMinutes))
	}

	htmlPath := os.Getenv("HTML_PATH")
	if htmlPath == "" {
		panic("HTML_PATH is not set")
	}

	cssPath := os.Getenv("CSS_PATH")
	if cssPath == "" {
		panic("CSS_PATH is not set")
	}

	go func() {
		for {
			if err := fetchImage(imagePath, picsumUrl); err != nil {
				fmt.Printf("failed to fetch image: %v\n", err)
			}
			time.Sleep(time.Duration(refreshInterval) * time.Minute)
		}
	}()

	r := gin.Default()
	r.LoadHTMLFiles(htmlPath)
	r.StaticFile("/style.css", cssPath)
	r.StaticFile("/image.jpg", imagePath)

	r.GET("/", func(c *gin.Context) {
		resp, err := http.Get(todoBackendUrl)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		var todos []Todo
		if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.HTML(200, "index.html", gin.H{"Todos": todos})
	})

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
