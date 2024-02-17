.PHONY: build run clean

all: clean build 

build: 
	go build -ldflags="-s -w" -o bin/goplay client/main.go
	go build -ldflags="-s -w" -o bin/goplayer server/main.go
	goupx bin/goplay bin/goplayer

completion:
	go build -o bin/goplay client/main.go
	go build -o bin/goplayer server/main.go
	

clean:
	rm -rf bin/
