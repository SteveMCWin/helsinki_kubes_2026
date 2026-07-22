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

	r := gin.Default()
	r.StaticFile("/", "./front/index.html")

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
