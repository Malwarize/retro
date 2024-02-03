.PHONY: build run clean

all: clean build run

build: 
	go build -o build/goplay

run:
	./build/goplay

clean:
	rm -rf build/

