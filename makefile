.PHONY: all build clean run

BUILD_DIR = bin
BINARY_NAME = go-downloader

all: build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

run: build
	$(BUILD_DIR)/$(BINARY_NAME)