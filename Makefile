
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get

BINARY_NAME := btwarch
BINARY_UNIX := $(BINARY_NAME)_unix
BINARY_WIN := $(BINARY_NAME).exe
BINARY_MAC := $(BINARY_NAME)_mac
BINARY_DIR := $(shell pwd)/bin

LDFLAGS := -ldflags="-s -w"

.PHONY: deps build build-linux build-windows build-mac build-all run clean db-up db-down db-logs migrate migrate-status

deps:
	$(GOGET) -v ./...

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v ./cmd/app/main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_UNIX) -v ./cmd/app/main.go

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_WIN) -v ./cmd/app/main.go

build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_MAC) -v ./cmd/app/main.go

build-all: build-linux build-windows build-mac

run:
	$(GOCMD) run cmd/app/main.go

db-up:
	docker-compose up -d postgres

db-down:
	docker-compose down

db-logs:
	docker-compose logs -f postgres

migrate:
	$(GOCMD) run cmd/migrate/main.go

migrate-status:
	$(GOCMD) run cmd/migrate/main.go status

clean:
	$(GOCLEAN)
	rm -f $(BINARY_DIR)/$(BINARY_NAME)
	rm -f $(BINARY_DIR)/$(BINARY_UNIX)
	rm -f $(BINARY_DIR)/$(BINARY_WIN)
	rm -f $(BINARY_DIR)/$(BINARY_MAC)
