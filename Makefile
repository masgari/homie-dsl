.PHONY: build
all: build

BUILD_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
#FALGS :=-ldflags="-s -w"
FALGS := ""

build:
	GO111MODULE=on go build $(FLAGS) -o $(GOPATH)/bin/homie-dsl main.go

run:
	GO111MODULE=on go run main.go

test:
	GO111MODULE=on go test -v -timeout 10s -coverprofile=/tmp/homie-dsl-test-coverage ./...

clean:	
	rm $(GOPATH)/bin/homie-dsl