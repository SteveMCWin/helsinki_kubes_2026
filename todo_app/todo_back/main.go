package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type newTodoRequest struct {
	Text string `json:"text"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		// port = "8081"
		panic("PORT is not set")
	}

	maxLength := os.Getenv("MAX_TODO_LENGTH")
	if maxLength == "" {
		panic("MAX_TODO_LENGTH is not set")
	}
	maxTodoLength, err := strconv.Atoi(maxLength)
	if err != nil {
		panic(fmt.Sprintf("MAX_TODO_LENGTH is not a number: %q", maxLength))
	}

	var mu sync.Mutex
	todos := []Todo{}

	r := gin.Default()

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
