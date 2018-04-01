# Go parameters
LDFLAGS="-s -w"
${HOME}=()

all: deps test gosint_build

deps:
	dep ensure

test:
	go test -v ./...

gosint_build:
	go build -o gosint -v -ldflags=${LDFLAGS} cmd/gosint/main.go

