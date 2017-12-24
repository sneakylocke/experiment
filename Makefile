GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get
BINARY_NAME=mybinary
BINARY_UNIX=$(BINARY_NAME)_unix

BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

build: lint
	$(GOBUILD) -v

test:
	$(GOTEST) -v -cover -coverprofile=coverage/coverage.out

cover: test
	$(GOTOOL) cover -html=coverage/coverage.out

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

lint: $(GOMETALINTER)
	gometalinter ./... --vendor
