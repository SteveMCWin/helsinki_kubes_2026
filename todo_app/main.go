package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const maxTodoLength = 140

type Todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type newTodoRequest struct {
	Text string `json:"text"`
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

	go func() {
		for {
			if err := fetchImage(imagePath); err != nil {
				fmt.Printf("failed to fetch image: %v\n", err)
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	var mu sync.Mutex
	todos := []Todo{
		{Text: "Learn Kubernetes", Done: false},
		{Text: "Deploy the todo app", Done: false},
		{Text: "Take a break", Done: true},
	}

	r := gin.Default()
	r.StaticFile("/", "./front/index.html")
	r.StaticFile("/style.css", "./front/style.css")
	r.StaticFile("/image.jpg", imagePath)

	r.GET("/todos", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()
		c.JSON(200, todos)
	})

	r.POST("/todos", func(c *gin.Context) {
		var body newTodoRequest
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		text := strings.TrimSpace(body.Text)
		if text == "" {
			c.JSON(400, gin.H{"error": "todo text must not be empty"})
			return
		}
		if len(text) > maxTodoLength {
			c.JSON(400, gin.H{"error": fmt.Sprintf("todo text must be %d characters or fewer", maxTodoLength)})
			return
		}

		todo := Todo{Text: text, Done: false}
		mu.Lock()
		todos = append(todos, todo)
		mu.Unlock()

		c.JSON(201, todo)
	})

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
