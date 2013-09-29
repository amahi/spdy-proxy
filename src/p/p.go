package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/SlyMarbo/spdy"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const HOST_PORT_API = "localhost:1443"
const HOST_PORT_SERVERS = "localhost:1444"

var incoming, serving int

// responseCopier does the copying of the request
// from H to C and the response from C to H.
type responseCopier struct {
	w http.ResponseWriter
	s spdy.Stream
}

func (r *responseCopier) ReceiveData(_ *http.Request, data []byte, final bool) {
	if data == nil || len(data) == 0 {
		return
	}
	_, err := r.w.Write(data)
	if err != nil {
		r.s.Close()
	}
}

func (r *responseCopier) ReceiveHeader(_ *http.Request, header http.Header) {
	h := r.w.Header()
	status := -1
	for key, values := range header {
		for _, value := range values {
			h.Add(key, value)
			if key == ":status" {
				if i := strings.Index(value, " "); i > 0 {
					value = value[:i]
				}
				s, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("Warning: Failed to parse status code %q.\n", value)
					continue
				}
				status = s
			}
		}
	}
	if status > 0 {
		r.w.WriteHeader(status)
	}
}

func (r *responseCopier) ReceiveRequest(_ *http.Request) bool {
	return false
}

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
	conn spdy.Conn
}

func (p *Proxy) RequestFromC(w http.ResponseWriter, r *http.Request) error {
	if p.conn == nil {
		fmt.Println("Warning: Could not serve request because C is not connected.")
		http.NotFound(w, r)
		return nil
	}

	copier := new(responseCopier)
	copier.w = w
	u := r.URL
	if u.Host == "" {
		u.Host = HOST_PORT_API
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	incoming++
	stream, err := p.conn.Request(r, copier, spdy.DefaultPriority(r.URL))
	if err != nil {
		return err
	}
	copier.s = stream
	serving++
	ret := stream.Run()
	incoming--
	serving--
	return ret
}

func (p *Proxy) ServeC(w http.ResponseWriter, r *http.Request) {
	// clean up the old connection
	if p.conn != nil {
		handle(p.conn.Close())
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
	res := new(http.Response)
	res.Status = "200 Connection Established"
	res.StatusCode = 200
	res.Proto = "HTTP/1.1"
	res.ProtoMajor = 1
	res.ProtoMinor = 1
	buf.Reset()
	message := "Hello from P"
	buf.WriteString(message)
	res.Body = buf
	res.ContentLength = int64(len(message))
	handle(res.Write(conn))

	// prepare for serving requests from H.
	client, err := spdy.NewClientConn(conn, nil, 3)
	handle(err)
	p.conn = client
	fmt.Println("Ready")
	client.Run()
}

func (p *Proxy) ServeA(w http.ResponseWriter, r *http.Request) {
	err := p.RequestFromC(w, r)
	if err != nil {
		fmt.Println("Encountered an error serving API request:", err)
	}
}

func (p *Proxy) DebugURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "incoming: %d\nserving: %d\n", incoming, serving)
}

func main() {
	certFile := "cert.pem"
	keyFile := "cert.key"

	tls := flag.Bool("t", false, "enable TLS")
	spdy_debug := flag.Bool("s", false, "enable SPDY debug output")
	flag.Parse()

	if *spdy_debug {
		// enable spdy debug messages
		spdy.EnableDebugOutput()
	}

	proxy := new(Proxy)
	http.HandleFunc("/", proxy.ServeC)
	go http.ListenAndServeTLS(HOST_PORT_SERVERS, certFile, keyFile, nil) // Serve C

	hServe := new(http.Server)
	mux := http.NewServeMux()
	mux.HandleFunc("/", proxy.ServeA)
	mux.HandleFunc("/debug", proxy.DebugURL)
	hServe.Handler = mux
	hServe.Addr = HOST_PORT_API
	spdy.AddSPDY(hServe)
	if *tls {
		fmt.Println("Serving on", HOST_PORT_API, "with TLS")
		handle(hServe.ListenAndServeTLS(certFile, keyFile)) // Serve H
	} else {
		fmt.Println("Serving on", HOST_PORT_API, "*without* TLS")
		handle(hServe.ListenAndServe()) // Serve H
	}
}
