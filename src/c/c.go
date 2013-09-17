package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/SlyMarbo/spdy"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

const HOST_PORT = "localhost:1444"
const DIR_TO_SERVE = "./media"

func main() {
	for {
		var conn *tls.Conn
		var err error
		for i := 0; i < 10; i++ {
			// connect to P.
			conn, err := tls.Dial("tcp", HOST_PORT, &tls.Config{InsecureSkipVerify: true})
			if err != nil {
				time.Sleep(100 * time.Millisecond)
			} else {
				break
			}
		}
		if conn == nil {
			fmt.Println("Failed to connect. Waiting 30 seconds.")
			time.Sleep(30 * time.Second)
			continue
		}
	
		// build the request
		buf := new(bytes.Buffer)
		_, err = buf.WriteString("Hello from C")
		handle(err)
		req, err := http.NewRequest("PUT", "https://"+HOST_PORT, buf)
		handle(err)
	
		// make the client connection
		client := httputil.NewClientConn(conn, nil)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: Failed to make connection to P:", err)
			continue
		}
		buf.Reset()
		_, err = io.Copy(buf, res.Body)
		handle(err)
		fmt.Printf("%q from P: %q.\n", res.Status, buf.String())
	
		// swap
		//spdy.EnableDebugOutput()
		c, _ := client.Hijack()
		conn = c.(*tls.Conn)
		srv := &http.Server{Handler: http.FileServer(http.Dir(DIR_TO_SERVE))}
		spdy.AddSPDY(srv)
		server, err := spdy.NewServerConn(conn, srv, 3)
		if err != nil {
			fmt.Println("Encountered error creating SPDY server connection:", err)
			continue
		}
		fmt.Println("Ready")
		err = server.Run()
		if err != nil {
			fmt.Println("Encountered error serving:", err, "\nReconnecting to proxy...")
		}
	}
}
