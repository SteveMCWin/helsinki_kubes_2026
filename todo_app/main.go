package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func fetchImage(path string) error {
	resp, err := http.Get("https://picsum.photos/1200")
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

	r := gin.Default()
	r.StaticFile("/", "./front/index.html")
	r.StaticFile("/image.jpg", imagePath)

	fmt.Printf("Server started in port %s\n", port)
	r.Run(":" + port)
}
