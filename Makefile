.PHONY: build clean install test fmt vet

BINARY_NAME=qpass
BINARY_PATH=./$(BINARY_NAME)

build:
	go build -o $(BINARY_NAME) .

clean:
	rm -f $(BINARY_NAME)

install: build
	cp $(BINARY_NAME) /usr/local/bin/

test:
	go test -v ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint: fmt vet

all: lint build

.DEFAULT_GOAL := build
