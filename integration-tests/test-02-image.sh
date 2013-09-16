#!/bin/bash

. functions

curl $debug -D /tmp/test.header -k -s -o /tmp/downloaded \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1444" \
	"$base_url/image.jpg"

if ( valid_sum 'a37afec0825c483a906f32adbb70528b6d5867b4' &&
	valid_sum '6a7aa20315136543c4c34729e8ca533f2bb7a953' /tmp/test.header ); then
	pass;
else
	fail;
fi;
