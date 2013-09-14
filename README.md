SPDY Proxy
==========

Reference implementation of a barebones SPDY/HTTPS proxying server for semi-permanent back-end server.

Thanks to Jamie Hall, David Anderson and Derrick McKee for their help and contributions to this proxy server implementation.

Intro
=====

There is a client (front-end) API to the proxy for clients, called H.

The C (server) is a client that connects to P semi-permanently. Once it's connected, H calls are proxy'd to C, with TLS and SPDY, with C becoming a server for H. 

Once H ends the connection, C's SPDY connection is kept and reused for any other calls with H. Each call from H is funnelled to C via a SPDY stream.

	H	Proxy clients
	^
	|
	v
	P	Proxy
	^
	|
	v
	C	Client device that semi-permanently connect to P, then become Servers for H

It uses the excellent SlyMarbo/spdy library.

Both the server to H and the interface with C support SPDY.

(C) 2013, Amahi
