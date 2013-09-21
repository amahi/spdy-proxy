export GOPATH=$(shell pwd)

all: build

build:
	go get c
	go get p

clean:
	rm -rf bin pkg

update-lib:
	git submodule init
	git submodule update --merge
	(cd src/github.com/SlyMarbo/spdy && git pull origin master && git checkout master)
