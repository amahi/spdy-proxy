export GOPATH=$(shell pwd)

all: build

build:
	go get c
	go get p
	go get a

get-deps:
	go get -d

race:
	mkdir -p bin
	go build -race p && mv p bin/
	go build -race c && mv c bin/
	go build -race a && mv a bin/

clean:
	rm -rf bin pkg

update-lib:
	git submodule init
	git submodule update --merge
	(cd src/github.com/amahi/spdy && git pull origin master && git checkout master)
