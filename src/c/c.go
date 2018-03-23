package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
	"golang.org/x/net/http2"
	"io/ioutil"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

type handler struct {
	rt string
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {

	if body, err := ioutil.ReadAll(rq.Body); err == nil && len(body) > 0 {
		filename := "/tmp/postdat"
		err := ioutil.WriteFile(filename, body, 0666)
		if err != nil {
			fmt.Println(err)
		}
	}
	fileServer := http.FileServer(http.Dir(h.rt))
	fileServer.ServeHTTP(rw, rq)
}

const HOST_PORT = "localhost:1444"

func main() {

	root := flag.String("r", "./testdata", "root of the directory to serve")
	flag.Parse()

	for {
		const SLEEP_RETRY = 5
		var conn *tls.Conn
		var err error
		for i := 0; i < 10; i++ {
			// connect to P.
			conn, err = tls.Dial("tcp", HOST_PORT, &tls.Config{InsecureSkipVerify: true})
			if err != nil {
				time.Sleep(100 * time.Millisecond)
			} else {
				break
			}
		}
		if conn == nil {
			log.Println("Failed to connect. Waiting", SLEEP_RETRY, "seconds.")
			time.Sleep(SLEEP_RETRY * time.Second)
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
			log.Println("Error: Failed to make connection to P:", err)
			continue
		}
		buf.Reset()
		_, err = io.Copy(buf, res.Body)
		handle(err)
		fmt.Printf("%q from P: %q.\n", res.Status, buf.String())

		c, _ := client.Hijack()
		conn = c.(*tls.Conn)

		serverConnOpts := new(http2.ServeConnOpts)
		serverConnOpts.Handler = &handler{rt: *root}
		server := new(http2.Server)
		server.ServeConn(conn, serverConnOpts)
	}
}
