all: pb grpc

pb:
	cd pb && protoc -I=proto/access --go_out=:./  --go_opt=paths=source_relative proto/access/*.proto
	cd pb && protoc -I=proto/common -I=proto/idmgmt --go_out=:./ --go_opt=paths=source_relative proto/common/*.proto
	cd pb && protoc -I=proto/consensus -I=proto/common -I=proto/idmgmt --go_out=:./ --go_opt=paths=source_relative proto/consensus/*.proto
	cd pb && mkdir -p tbft mbft
	cd pb/tbft && protoc -I=../proto/consensus/tbft -I=../proto/common -I=../proto/idmgmt --go_out=:./ --go_opt=paths=source_relative ../proto/consensus/tbft/*.proto
	cd pb/mbft && protoc -I=../proto/consensus/mbft -I=../proto/common -I=../proto/idmgmt --go_out=:./ --go_opt=paths=source_relative ../proto/consensus/mbft/*.proto
	cd pb && protoc -I=proto/config -I=proto/consensus -I=proto/common -I=proto/idmgmt -I=proto/access --go_out=:./ --go_opt=paths=source_relative proto/config/*.proto
	cd pb && protoc -I=proto/idmgmt --go_out=:./ --go_opt=paths=source_relative proto/idmgmt/*.proto
	cd pb && protoc -I=proto/net -I=proto/common -I=proto/idmgmt --go_out=:./ --go_opt=paths=source_relative proto/net/*.proto
	cd pb && protoc -I=proto/store -I=proto/common -I=proto/idmgmt --go_out=:./ --go_opt=paths=source_relative proto/store/*.proto
	cd pb && protoc -I=proto/txpool -I=proto/common -I=proto/idmgmt --go_out=:./ --go_opt=paths=source_relative proto/txpool/*.proto
	cd pb && protoc -I=proto/common -I=proto/sync -I=proto/idmgmt --go_out=:./ --go_opt=paths=source_relative proto/sync/*.proto
	cd pb && protoc -I=proto/discovery --go_out=:./ --go_opt=paths=source_relative proto/discovery/*.proto

grpc:
	cd pb && protoc -I=proto/api -I=proto/common -I=proto/idmgmt -I=proto/config --go-grpc_out==plugins=grpc:. --go-grpc_opt=paths=source_relative proto/api/rpc_node.proto

dep:
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0
	go get -u google.golang.org/grpc

clean:
	rm -f pb/*.pb.go pb/*/*.pb.go

.PHONY: pb
