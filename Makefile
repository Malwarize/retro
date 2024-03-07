.PHONY: build run clean

all:  clean init build 

build: init
	go build -o bin/goplay client/main.go
	go build -o bin/goplayer server/main.go
	
clean: 
	rm -rf bin/

init:
	mkdir -p bin/
