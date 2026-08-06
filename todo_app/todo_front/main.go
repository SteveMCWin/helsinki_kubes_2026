package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func fetchImage(path string) error {
	resp, err := http.Get("https://picsum.photos/600")
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
		port = "8080"
	}

	imagePath := os.Getenv("IMAGE_FILE_PATH")
	if imagePath == "" {
		imagePath = "/data/image.jpg"
	}

	todoBackendUrl := os.Getenv("TODO_BACKEND_URL")
	if todoBackendUrl == "" {
		todoBackendUrl = "http://todoapp-svc:1235/todos"
	}

	go func() {
		for {
			if err := fetchImage(imagePath); err != nil {
				fmt.Printf("failed to fetch image: %v\n", err)
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	r := gin.Default()
	r.LoadHTMLFiles("./front/index.html")
	r.StaticFile("/style.css", "./front/style.css")
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
