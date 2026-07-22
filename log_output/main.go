package main

import (
	"crypto/rand"
	"fmt"
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
	id := newUUID()

	for {
		fmt.Printf("%s: %s\n", time.Now().UTC().Format("2006-01-02T15:04:05.000Z"), id)
		time.Sleep(5 * time.Second)
	}
}
