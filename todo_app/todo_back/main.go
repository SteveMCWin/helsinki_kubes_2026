package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID   int    `json:"id"`
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

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS todos (id SERIAL PRIMARY KEY, text TEXT NOT NULL, done BOOLEAN NOT NULL DEFAULT false)"); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Status(200)
	})

	r.GET("/todos", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, text, done FROM todos ORDER BY id")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		todos := []Todo{}
		for rows.Next() {
			var todo Todo
			if err := rows.Scan(&todo.ID, &todo.Text, &todo.Done); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			todos = append(todos, todo)
		}

		c.JSON(200, todos)
	})

	r.POST("/todos", func(c *gin.Context) {
		var body newTodoRequest
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		text := strings.TrimSpace(body.Text)
		log.Printf("received todo: %q", text)

		if text == "" {
			log.Printf("rejected todo: text must not be empty")
			c.JSON(400, gin.H{"error": "todo text must not be empty"})
			return
		}
		if len(text) > maxTodoLength {
			msg := fmt.Sprintf("todo text must be %d characters or fewer", maxTodoLength)
			log.Printf("rejected todo: %s", msg)
			c.JSON(400, gin.H{"error": msg})
			return
		}

		todo := Todo{Text: text, Done: false}
		if err := db.QueryRow("INSERT INTO todos (text, done) VALUES ($1, $2) RETURNING id", todo.Text, todo.Done).Scan(&todo.ID); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, todo)
	})

	r.DELETE("/todos/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid todo id"})
			return
		}

		result, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if rowsAffected == 0 {
			c.JSON(404, gin.H{"error": "todo not found"})
			return
		}

		c.Status(204)
	})

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
