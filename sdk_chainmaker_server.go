package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"context"
	"fmt"
)

func (cc ChainClient) GetChainMakerServerVersion() (string, error) {
	cc.logger.Debug("[SDK] begin to get chainmaker server version")
	req := &pb.ChainMakerVersionRequest{}
	client, err := cc.pool.getClient()
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	res, err := client.rpcNode.GetChainMakerVersion(ctx, req)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", fmt.Errorf("get chainmaker server version failed, %s", res.Message)
	}
	return res.Version, nil
}
