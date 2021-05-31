all: build gomod

build:
	go build ./...

gomod:
	go get -u chainmaker.org/chainmaker/pb-go
	go get -u chainmaker.org/chainmaker/common

.PHONY: all
