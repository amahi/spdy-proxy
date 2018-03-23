package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
	"golang.org/x/net/http2"
)

const HOST_PORT_API = "0.0.0.0:1443"
const HOST_PORT_SERVERS = "0.0.0.0:1444"

type stats_s struct {
	sync.Mutex
	incoming int
	serving  int
}

var stats stats_s

// Used in sending the response.
// Essentially, this is just adding
// the Close method so that it fulfils
// the io.ReadCloser interface.
type buffer struct {
	bytes.Buffer
}

func (b *buffer) Close() error {
	return nil
}

// placeholder for proper error handling.
func handle(err error) {
	if err != nil {
		panic(err)
	}
}

type Proxy struct {
	clientConn *http2.ClientConn
}

func (p *Proxy) RequestFromC(w http.ResponseWriter, r *http.Request) error {
	if p.clientConn == nil {
		log.Println("Warning: Could not serve request because C is not connected.")
		http.NotFound(w, r)
		return nil
	}

	u := r.URL
	if u.Host == "" {
		u.Host = HOST_PORT_API
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	// err := p.session.NewStreamProxy(r, w)
	res, err := p.clientConn.RoundTrip(r)

	buf := new(buffer)
	_, err = io.Copy(buf, res.Body)
	handle(err)
	handle(res.Body.Close())
	fmt.Printf("Reply from C: %q.\n", buf.String())

	return err
}

func (p *Proxy) ServeC(w http.ResponseWriter, r *http.Request) {
	// clean up the old connection
	if p.clientConn != nil {
		// p.clientConn.Close()
	}

	// Read in the request body.
	buf := new(buffer)
	_, err := io.Copy(buf, r.Body)
	handle(err)
	handle(r.Body.Close())
	fmt.Printf("%q from C: %q.\n", r.Method, buf.String())

	// re-purpose the connection.
	conn, _, err := w.(http.Hijacker).Hijack()
	handle(err)

	// send the 200 to C.
	buf.Reset()
	buf.WriteString("Hello from P")

	res := &http.Response{
		Status:        "200 Connection Established",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          buf,
		ContentLength: int64(buf.Len()),
	}

	handle(res.Write(conn))

	// prepare for serving requests from A.
	transport := new(http2.Transport)
	p.clientConn, err = transport.NewClientConn(conn) 
	fmt.Println("Ready")
}

func (p *Proxy) ServeA(w http.ResponseWriter, r *http.Request) {
	err := p.RequestFromC(w, r)
	if err != nil {
		log.Println("Encountered an error serving API request:", err)
	}
}

func (p *Proxy) DebugURL(w http.ResponseWriter, r *http.Request) {
	stats.Lock()
	fmt.Fprintf(w, "goroutines:  %d\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "incoming: %d\nserving: %d\n", stats.incoming, stats.serving)
	stats.Unlock()
}

func main() {
	certFile := "cert.pem"
	keyFile := "cert.key"

	tls := flag.Bool("t", false, "enable TLS")
	flag.Parse()

	proxy := new(Proxy)
	http.HandleFunc("/", proxy.ServeC)

	go func(){ // Serve C
		err := http.ListenAndServeTLS(HOST_PORT_SERVERS, certFile, keyFile, nil)
		handle(err)
	}()

	hServe := new(http.Server)
	mux := http.NewServeMux()
	mux.HandleFunc("/", proxy.ServeA)
	mux.HandleFunc("/debug", proxy.DebugURL)
	hServe.Handler = mux
	hServe.Addr = HOST_PORT_API
	// hServe.WriteTimeout = 10 * time.Second
	// hServe.ReadTimeout = 10 * time.Second
	// spdy.AddSPDY(hServe)
	if *tls {
		fmt.Println("Serving on", HOST_PORT_API, "with TLS")
		handle(hServe.ListenAndServeTLS(certFile, keyFile)) // Serve H
	} else {
		fmt.Println("Serving on", HOST_PORT_API, "*without* TLS")
		handle(hServe.ListenAndServe()) // Serve H
	}
}
