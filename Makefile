# Heavily based off https://sohlich.github.io/post/go_makefile/
# pared down from the above and added deps download via go mod

# Go parameters
GOCMD=/usr/local/go/bin/go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GETDEPS=$(GOCMD) mod download
BINARY_NAME=ipdetails.out

all: deps test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v ./example
		@echo "build done run with: ./$(BINARY_NAME)"
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./example
		./$(BINARY_NAME)
deps:
		$(GOGETGETDEPS)
