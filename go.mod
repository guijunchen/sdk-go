module chainmaker.org/chainmaker-sdk-go

go 1.15

require (
	chainmaker.org/chainmaker-go/common v0.0.0
	chainmaker.org/chainmaker-sdk-pb v0.0.0
	github.com/Rican7/retry v0.1.0
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/hokaccha/go-prettyjson v0.0.0-20201222001619-a42f9ac2ec8e
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.16.0
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/grpc v1.36.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
)

replace (
	chainmaker.org/chainmaker-go/common => ./common
	chainmaker.org/chainmaker-sdk-pb => ./pb/protogo
	chainmaker.org/chainmaker-go/pb => ./pb/protogo
)
