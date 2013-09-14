export GOPATH=$(shell pwd)

all: build

build:
	go get c
	go get p

clean:
	rm -rf bin pkg
