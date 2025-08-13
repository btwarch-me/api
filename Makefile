GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=btwarch

deps:
	$(GOGET) -v ./...

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go
	./$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
