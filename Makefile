.PHONY: build run clean

all:  clean init build 

build: init build-client build-server

build-client: init
	go build -o bin/goplay client/main.go

build-server: init
	go build -o bin/goplayer server/main.go

clean: 
	rm -rf bin/

init:
	mkdir -p bin/
