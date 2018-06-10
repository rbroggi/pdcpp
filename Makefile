SHELL=/bin/bash
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=gofmt
BINARY_NAME=$(shell echo $${PWD\#\#*/})
BINARY_UNIX=$(BINARY_NAME)_unix
PROTO_COMPILER=protoc
VERSION=1.0.0

# The name of the executable (default is current directory name)
TARGET=./target/$(BINARY_NAME)
TARGET_UNIX=./target/$(BINARY_UNIX)
PROTO_MESSAGE_DIR_FILE_PATH=./proto/
PROTO_MESSAGE_FILE_PATH=./proto/grpc_pdc_price_provider.proto
PROTO_OUT=pdc_trade/
GRPC_API_CONFIG_PATH=./config/PdcTradePricesService.yaml
SWAGGER_DOC_OUT=swagger/
BUILD := $(git rev-parse HEAD)
# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY: all build clean install uninstall fmt simplify check run

all: regenerate gofmt test build

build:
	$(GOBUILD) -o $(TARGET) -v

goget:
	$(GOGET) -u google.golang.org/grpc
	$(GOGET) -u github.com/golang/protobuf/protoc-gen-go
	$(GOGET) -u github.com/spf13/cobra/cobra
	$(GOGET) -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	$(GOGET) -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger


regenerate:
	mkdir -p $(PROTO_OUT)
	$(PROTO_COMPILER) -I $(PROTO_MESSAGE_DIR_FILE_PATH) --go_out=plugins=grpc:$(PROTO_OUT) $(PROTO_MESSAGE_FILE_PATH)
	$(PROTO_COMPILER) -I $(PROTO_MESSAGE_DIR_FILE_PATH) --grpc-gateway_out=logtostderr=true,grpc_api_configuration=$(GRPC_API_CONFIG_PATH):$(PROTO_OUT) $(PROTO_MESSAGE_FILE_PATH)
	$(PROTO_COMPILER) -I $(PROTO_MESSAGE_DIR_FILE_PATH) --swagger_out=logtostderr=true,grpc_api_configuration=$(GRPC_API_CONFIG_PATH):$(SWAGGER_DOC_OUT) $(PROTO_MESSAGE_FILE_PATH)

test:
	$(GOTEST) -v ./...

gofmt:
	$(GOFMT) -l -s -w .

clean:
	$(GOCLEAN)
	rm -f $(TARGET)
	rm -f $(TARGET_UNIX)
	rm -rf $(PROTO_OUT)


# Cross compilation
build-linux:
				CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
				docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v

