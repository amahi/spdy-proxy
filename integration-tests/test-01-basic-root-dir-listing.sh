#!/bin/bash

. functions

curl $debug -D /tmp/test.header -k -s -o /tmp/downloaded \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1444" \
	"$base_url/"

if ( valid_sum 'a8d4b2f1cab59e48ca02cf700e0f94155d421271' &&
	valid_sum '58442a6f911eedae4033c190e91a476400bb6867' /tmp/test.header ); then
	pass;
else
	fail;
fi;
