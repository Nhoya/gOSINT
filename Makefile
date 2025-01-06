# Go parameters
LDFLAGS="-s -w"

all: test build

build:
	go build -ldflags $(LDFLAGS) ./...

test:
	go test -v ./...

.PHONY: all build test
