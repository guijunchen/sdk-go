module chainmaker.org/chainmaker-sdk-pb

go 1.15

require (
	chainmaker.org/chainmaker-go/pb v0.0.0-00010101000000-000000000000
	github.com/gogo/protobuf v1.3.2
	google.golang.org/grpc v1.36.0
)

replace chainmaker.org/chainmaker-go/pb => ./
