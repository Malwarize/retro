.PHONY: build run clean

all:  clean init build 

build: init build-client build-server

build-client: init
	go build -ldflags "-w -s" -o bin/retro client/main.go

build-server: init
	go build -ldflags "-w -s" -o bin/retroPlayer server/main.go

clean: 
	rm -rf bin/

init:
	mkdir -p bin/
