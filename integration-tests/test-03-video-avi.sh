#!/bin/bash

. ./functions

curl $debug -D /tmp/test.header -k -s -o /tmp/downloaded \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1444" \
	"$base_url/sample.avi"

if ( valid_sum '56c8ca5e53defe7bac756bffb9893099a95a67fd' &&
	valid_sum '60b28f3ec6826ba9aa6a62661e8b9c410f8afcd6' /tmp/test.header ); then
	pass;
else
	fail;
fi;
