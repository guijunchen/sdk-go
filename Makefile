all: build gomod

build:
	go build ./...

gomod:
	go get -u chainmaker.org/chainmaker/pb-go@v2.0.0_qc
	go get -u chainmaker.org/chainmaker/common@v2.0.0_qc

lint:
	golangci-lint run .

.PHONY: all
