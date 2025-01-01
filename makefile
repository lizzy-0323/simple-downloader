.PHONY: all build clean run

BUILD_DIR = bin

all: build

build:
	go build -o $(BUILD_DIR)/downloader main.go

clean:
	rm -f $(BUILD_DIR)/downloader

run: build
	$(BUILD_DIR)/downloader