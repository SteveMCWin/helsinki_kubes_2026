package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func newUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func timestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

func main() {
	id := newUUID()

	go func() {
		for {
			fmt.Printf("%s: %s\n", timestamp(), id)
			time.Sleep(5 * time.Second)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "5607"
	}

	r := gin.Default()
	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"timestamp": timestamp(),
			"string":    id,
		})
	})

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
