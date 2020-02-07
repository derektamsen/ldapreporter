 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOVET=$(GOCMD) vet
GOTEST=$(GOCMD) test -race
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
BINARY_NAME=ldapreporter
BINARY_UNIX=$(BINARY_NAME)_unix
GIT_SHA=$(shell git rev-list -1 HEAD)

all: test build

build:
	$(GOBUILD) -ldflags \
		"-X main.Version=$(GIT_SHA)" \
		-o $(BINARY_NAME) -v ./...

vet:
	$(GOVET) ./...

test: vet
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	LOG_LEVEL=INFO $(GORUN) -ldflags "-X main.Version=$(GIT_SHA)" ./$(BINARY_NAME).go -loglevel "INFO"

dev:
	docker-compose up && docker-compose down
