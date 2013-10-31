#!/bin/bash

. ./functions

curl $debug -D /tmp/test.header -k -s -o /tmp/test.body \
	-H "Accept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3" \
	-H "Connection: keep-alive" \
	-H "Host: localhost:1444" \
	"$base_url/sample.mp4"

valid_head 'ad86b41f562eb170832434b8ed0b7af09812f40c';
valid_body 'dcff5a900a60c7615be5f74c145a6487bda02489';
