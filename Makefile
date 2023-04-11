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
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags \
		"-X main.Version=$(GIT_SHA)" \
		-o $(BINARY_NAME) -v ./...

vet:
	$(GOVET) ./...

test: vet
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME) $(BINARY_UNIX)
	rm -r dist

run:
	LOG_LEVEL=INFO $(GORUN) -ldflags "-X main.Version=$(GIT_SHA)" \
		./$(BINARY_NAME).go \
		-loglevel "INFO" -server "ldap://localhost:8389" \
		-user "cn=admin,dc=planetexpress,dc=com" \
		-password "GoodNewsEveryone" \
		-basedn "dc=planetexpress,dc=com" \
		-searchfilter "(&(objectclass=Group))"

dev:
	docker-compose up && docker-compose down

release:
	goreleaser release --snapshot --clean
	goreleaser check
