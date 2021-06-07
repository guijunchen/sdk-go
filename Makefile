SED = sed -i
ifeq ($(shell uname -s),Darwin)
	SED = sed -i ''
endif

all: replace_go_package create_proto_dir pb grpc reset_go_package

replace_go_package:
	${SED} "s%chainmaker.org/chainmaker-go/pb/protogo%chainmaker.org/chainmaker-sdk-go/pb/protogo%g" pb/proto/*/*.proto

reset_go_package:
	${SED} "s%chainmaker.org/chainmaker-sdk-go/pb/protogo%chainmaker.org/chainmaker-go/pb/protogo%g" pb/proto/*/*.proto

create_proto_dir:
	mkdir -p pb/protogo

pb:
	cd pb/proto && protoc -I=. --gogofaster_out=:../protogo --gogofaster_opt=paths=source_relative accesscontrol/*.proto
	cd pb/proto && protoc -I=. --gogofaster_out=:../protogo --gogofaster_opt=paths=source_relative common/*.proto
	cd pb/proto && protoc -I=. --gogofaster_out=:../protogo --gogofaster_opt=paths=source_relative consensus/*.proto
	cd pb/proto && protoc -I=. --gogofaster_out=:../protogo --gogofaster_opt=paths=source_relative consensus/tbft/*.proto
	#cd pb/proto && protoc -I=. --gogofaster_out=:../protogo --gogofaster_opt=paths=source_relative consensus/mbft/*.proto
	cd pb/proto && protoc -I=. --gogofaster_out=:../protogo --gogofaster_opt=paths=source_relative config/*.proto
	cd pb/proto && protoc -I=. --gogofaster_out=:../protogo --gogofaster_opt=paths=source_relative net/*.proto
	cd pb/proto && protoc -I=. --gogofaster_out=:../protogo --gogofaster_opt=paths=source_relative store/*.proto
	cd pb/proto && protoc -I=. --gogofaster_out=:../protogo --gogofaster_opt=paths=source_relative txpool/*.proto
	cd pb/proto && protoc -I=. --gogofaster_out=:../protogo --gogofaster_opt=paths=source_relative sync/*.proto
	cd pb/proto && protoc -I=. --gogofaster_out=:../protogo --gogofaster_opt=paths=source_relative discovery/*.proto
grpc:
	cd pb/proto && protoc -I=. --go-grpc_out==plugins=grpc:../protogo --go-grpc_opt=paths=source_relative api/rpc_node.proto

dep:
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0
	go get -u google.golang.org/grpc
	go get -u github.com/gogo/protobuf/protoc-gen-gogofaster

clean:
	rm -rf pb/protogo/*

.PHONY: pb
