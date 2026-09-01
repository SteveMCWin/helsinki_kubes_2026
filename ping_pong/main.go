package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS counter (count INTEGER)"); err != nil {
		panic(err)
	}
	if _, err := db.Exec("INSERT INTO counter (count) SELECT 0 WHERE NOT EXISTS (SELECT 1 FROM counter)"); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		var current int
		if err := db.QueryRow("UPDATE counter SET count = count + 1 RETURNING count - 1").Scan(&current); err != nil {
			c.String(500, "db error: %v", err)
			return
		}

		c.String(200, "pong %d", current)
	})
	r.GET("/pings", func(c *gin.Context) {
		var count int
		if err := db.QueryRow("SELECT count FROM counter").Scan(&count); err != nil {
			c.String(500, "db error: %v", err)
			return
		}

		c.String(200, "%d", count)
	})

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
