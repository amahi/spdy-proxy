#!/bin/bash

. ./functions

curl $debug -D /tmp/test.header -k -s -o /tmp/downloaded \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1444" \
	"$base_url/sample.mkv"

if ( valid_sum '48c9a4e3a24324e3d3ac6b284f7502878ace1909' &&
	valid_sum 'ce151f300431922dd8e57d983bbddede3a9211d7' /tmp/test.header ); then
	pass;
else
	fail;
fi;
