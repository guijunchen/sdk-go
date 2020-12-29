module chainmaker.org/chainmaker-go/chainmaker-sdk-go

go 1.13

require (
	chainmaker.org/chainmaker-go/common v0.0.0
	github.com/Rican7/retry v0.1.0
	github.com/envoyproxy/go-control-plane v0.9.4
	github.com/golang/protobuf v1.4.1
	github.com/hokaccha/go-prettyjson v0.0.0-20201222001619-a42f9ac2ec8e
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
)

replace chainmaker.org/chainmaker-go/common => ./common
