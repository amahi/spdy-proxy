#!/bin/bash

. functions

curl $debug -D /tmp/test.header -k -s -o /tmp/downloaded \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1444" \
	"$base_url/sample.mp4"

if ( valid_sum 'dcff5a900a60c7615be5f74c145a6487bda02489' &&
	valid_sum 'e4659604854987ecfba1ba741042628ba291c516' /tmp/test.header ); then
	pass;
else
	fail;
fi;
