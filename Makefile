# Go parameters
LDFLAGS="-s -w"

all: deps test gosint_build

deps:
	dep ensure

test:
	go test -v ./...

gosint_build:
	go build -v -ldflags=${LDFLAGS} cmd/gosint.go
