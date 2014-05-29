package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/amahi/spdy"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

type handler struct {
        data []byte
        rt string
}
func (h *handler) ServeHTTP(rw http.ResponseWriter,rq *http.Request) {
        if rq.Body!=nil {
                h.data = make([]byte, int(rq.ContentLength))
                _,err := rq.Body.(io.Reader).Read(h.data)
                if err != nil {
                        fmt.Println(err)
                }
                filename := "/tmp/postdat"
                f, err := os.Create(filename)
                if err != nil {
                        fmt.Println(err)
                }
                n, err := f.Write(h.data)
                if err != nil {
                        fmt.Println(n, err)
                }
                f.Close()
        }
        fileserver := http.FileServer(http.Dir(h.rt))
        fileserver.ServeHTTP(rw, rq)
}

const HOST_PORT = "localhost:1444"

func main() {

	root := flag.String("r", "./testdata", "root of the directory to serve")
	spdy_debug := flag.Bool("s", false, "enable SPDY debug output")
	flag.Parse()

	if *spdy_debug {
		// enable spdy debug messages
		spdy.EnableDebug()
	}

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
		server := new(http.Server)
		server.Handler = &handler{data:nil,rt:*root}
		//http.FileServer(http.Dir(*root))
		session := spdy.NewServerSession(conn, server)
		session.Serve()
	}
}
