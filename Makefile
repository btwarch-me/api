
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get

BINARY_NAME := btwarch_api
BINARY_UNIX := $(BINARY_NAME)_unix
BINARY_WIN := $(BINARY_NAME).exe
BINARY_MAC := $(BINARY_NAME)_mac

LDFLAGS := -ldflags="-s -w"

.PHONY: deps build build-linux build-windows build-mac build-all run clean db-up db-down db-logs

deps:
	$(GOGET) -v ./...

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v ./cmd/main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) -v ./cmd/main.go

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_WIN) -v ./cmd/main.go

build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_MAC) -v ./cmd/main.go

build-all: build-linux build-windows build-mac

run:
	$(GOCMD) run cmd/main.go

db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down

db-logs:
	docker-compose logs -f postgres

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WIN)
	rm -f $(BINARY_MAC)
