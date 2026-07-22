package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"
)

func newUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func main() {
	path := os.Getenv("LOG_FILE_PATH")
	if path == "" {
		path = "/data/status.log"
	}

	id := newUUID()

	for {
		line := fmt.Sprintf("%s: %s\n", time.Now().UTC().Format("2006-01-02T15:04:05.000Z"), id)

		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		if _, err := f.WriteString(line); err != nil {
			panic(err)
		}
		f.Close()

		fmt.Print(line)
		time.Sleep(5 * time.Second)
	}
}
