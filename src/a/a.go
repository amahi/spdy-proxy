package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const BUFFER_SIZE = 256
const SLEEP_TIME = 200

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	resp, err := http.Get("http://localhost:1443/sample.avi")
	handle(err)
	defer resp.Body.Close()

	f, err := os.OpenFile("/tmp/outfile", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	handle(err)

	copied := int64(0)
	for {
		written, err := io.CopyN(f, resp.Body, BUFFER_SIZE)
		copied += written
		if err != nil {
			// EOF -- ignore the rest
			break
		}
		time.Sleep(SLEEP_TIME * time.Millisecond)
	}

	fmt.Printf("Received: %d bytes\n", copied)
}

