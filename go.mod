module chainmaker.org/chainmaker-sdk-go

go 1.15

require (
	chainmaker.org/chainmaker-go/common v0.0.0
	github.com/Rican7/retry v0.1.0
	github.com/ethereum/go-ethereum v1.10.2
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.4.3
	github.com/hokaccha/go-prettyjson v0.0.0-20201222001619-a42f9ac2ec8e
	github.com/samkumar/hibe v0.0.0-20171013061409-c1cd171b6178
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.16.0
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/grpc v1.36.0
)

replace chainmaker.org/chainmaker-go/common => ./common
