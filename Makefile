GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GETDEPS=$(GOCMD) mod download
BINARY_NAME=ipd

all: deps test build

install: $(BINARY_NAME)
		mv $(BINARY_NAME) /usr/local/bin/ipd
build:
		$(GOBUILD) -o $(BINARY_NAME) -v ./app/main.go
		@echo "build done run with: ./$(BINARY_NAME)"
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		./$(BINARY_NAME)
deps:
		$(GOGETGETDEPS)
