
package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/SlyMarbo/spdy"
	"io"
	"net/http"
	"net/http/httputil"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

const HOST_PORT = "proxy.example.com:1444"
const DIR_TO_SERVE = "."

func main() {
	// connect to P
	conn, err := tls.Dial("tcp", HOST_PORT, &tls.Config{InsecureSkipVerify: true})
	handle(err)

	// build the request
	buf := new(bytes.Buffer)
	_, err = buf.WriteString("Hello from C")
	handle(err)
	req, err := http.NewRequest("PUT", "https://"+HOST_PORT, buf)
	handle(err)

	// make the client connection
	client := httputil.NewClientConn(conn, nil)
	res, err := client.Do(req)
	handle(err)
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
	handle(err)
	fmt.Println("Ready")
	handle(server.Run())
}
