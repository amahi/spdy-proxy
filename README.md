SPDY Proxy [![Build Status](https://travis-ci.org/amahi/spdy-proxy.png?branch=master)](https://travis-ci.org/amahi/spdy-proxy)
==========

Reference implementation in [Go](http://golang.org/) of a proxy server for (semi-)permanently connected back-end servers, via SPDY.

Thanks to Nilesh Jagnik, Derrick McKee and Jamie Hall for their help and contributions to how this proxy server idea and implementation.

This project is an extracion from a larger Amahi project. It uses Amahi's [amahi/spdy](https://github.com/amahi/spdy/) SPDY library. Amahi's service is now live and the [Amahi mobile iOS app](https://www.amahi.org/ios) has been released using the service.

Intro
=====

This is a proxy server process P and a back-end origin server for it C (which for historic reasons is called a "client", since it connects to P before becoming a server).

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
	C	Client "origin server"
	
The origin server semi-permanently connects to P (it reconnects in case of failures), then becomes a server for A


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

The system works by creating a gateway between `1` connection to `C` and `n` connections to `A`s.
This uses Go's `net/http` package design to simplify the interfaces. The system starts by opening
a port with which servers may connect. When a `C` connects, it does so as a client, with the proxy
(`P`) fulfilling the role of server. At this point, a brief authentication phasetakes place, using
an HTTPS PUT request. Once this request has completed, the underlying TLS connection is hijacked
by the application at both ends. This is then repurposed as a SPDY connection where `C` is the
server and `P` is the client. This enables us to connect the two without opening a port at `C`.

At this point, we have a SPDY connection between `C` and `P`, so `P` is now ready to pass requests to `C`.
When requests arrive at `P`, irrespective of their origin, the request and a callback are prepared
and submitted to the SPDY connection. When response data arrives from `C`, this is passed on to the
requesting client `A`n by the created callback structure.

This separation of inputs and outputs should result in equal treatment of all requests (except where
SPDY request priorities are used), which should give good usage performance.

(C) 2014, Amahi
