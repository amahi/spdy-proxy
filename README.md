spdy-proxy
==========

Reference implementation of a barebones SPDY/HTTPS proxying server for semi-permanent back-end server.

	H	Proxy clients
	^
	|
	v
	P	Proxy
	^
	|
	v
	C	Client devices that semi-permanently connect to P, then become Servers for H

There is a client (front-end) API to the proxy for clients, called H.

The C client connect to P semi-permanently. Once it's connected, H calls are proxy'd to C, with TLS and SPDY.

It uses the excellent SlyMarbo/spdy library.

(C) 2013, Amahi
