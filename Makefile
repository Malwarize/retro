.PHONY: build run clean

all:  clean init build 

build: init
	go build -ldflags="-s -w" -o bin/goplay client/main.go
	go build -ldflags="-s -w" -o bin/goplayer server/main.go
	
clean: 
	rm -rf bin/

init:
	mkdir -p bin/
