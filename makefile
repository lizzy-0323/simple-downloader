.PHONY: all build clean run

all: build

build:
	go build -o downloader main.go

clean:
	rm -f downloader

run: build
	./downloader