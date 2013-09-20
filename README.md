SPDY Proxy
==========

Reference implementation in [Go](http://golang.org/) of a proxy server for (semi-)permanently connected back-end servers, supporting SPDY and HTTPS.

Thanks to Jamie Hall, David Anderson and Derrick McKee for their help and contributions to how this proxy server idea and implementation.

This project is an extracion from a larger Amahi project and Jamie did this code based on a description of how the Amahi project was architected. It uses Jamie's excellent [SlyMarbo/spdy](https://github.com/SlyMarbo/spdy/) SPDY library.

Intro
=====

This is a proxy server process P and a back-end client for it C.

The proxy P serves (on two ports) a client (front-end) API to serve requests from clients, called A.

The C (server) is a client that connects to P semi-permanently. Once it's connected, A calls are proxy'd to C, with TLS and SPDY, with C becoming a server for A. 

Once A ends the connection, C's SPDY connection is kept and reused for any other calls with A. Each call from A is funnelled to C via a SPDY stream.

	A	Proxy clients
	^
	|
	v
	P	Proxy
	^
	|
	v
	C	Client device that semi-permanently connect to P, then become Servers for A

Both the server to A and the interface with C support SPDY.

Testing
=======

To test, run, e.g. in two separate windows, the client and the server:

    $ ./bin/p
    
    $ ./bin/c
    
You should see messages that they are connected. Then, in a third window, run the tests:

    $ cd integration-tests
    $ make
    Testing test-01-basic-root-dir-listing.sh
    PASS
    Testing test-02-image.sh
    PASS
    Testing test-03-video-avi.sh
    PASS
    Testing test-04-video-mkv.sh
    PASS
    Testing test-05-video-mp4.sh
    PASS
    
To add tests, see how these are added and follow the pattern. For control flow and concurrency tests, there may be more complex things that need to be done, for example using expect, killing connections, simultaneously starting connections, etc.

Design
======

The system works by creating a gateway between 1 connection to C and n connections to As.
This uses Go's `net/http` package design to simplify the interfaces. The system starts by opening
a port with which servers may connect. When a C connects, it does so as a client, with the proxy
(P) fulfilling the role of server. At this point, a brief authentication phasetakes place, using
an HTTPS PUT request. Once this request has completed, the underlying TLSconnection is hijacked
by the application at both ends. This is then repurposed as a SPDY connection where C is the server

(C) 2013, Amahi
