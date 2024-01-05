GOBUILD=go build
GOCLEAN=go clean
GOTEST=go test
BINARY_NAME=shredder

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/shredder

test:
	$(GOTEST) -v ./shred

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
