all: build gomod

build:
	go build ./...

gomod:
	go get -u chainmaker.org/chainmaker/pb-go@v2.0.0_dev
	go get -u chainmaker.org/chainmaker/common@v2.0.0_dev

.PHONY: all
