.PHONY: all test clean build release

GOFLAGS ?= $(GOFLAGS:)
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0
REPO = github.com/abulimov/db-checker
BINARY = db-checker

all: test

get:
	@go get $(GOFLAGS) -t ./...

release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o $(BINARY) $(REPO)
	zip -6 release/$(BINARY).linux-amd64.zip $(BINARY)
	GOOS=linux GOARCH=arm go build -o $(BINARY) $(REPO)
	zip -6 release/$(BINARY).linux-arm.zip $(BINARY)
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY) $(REPO)
	zip -6 release/$(BINARY).darwin-amd64.zip $(BINARY)

build: get
	go build $(GOFLAGS) $(REPO)

test: get
	@go test -v $(GOFLAGS) `go list ./... | grep -v /vendor/`

clean:
	@go clean $(GOFLAGS) -i $(REPO)
