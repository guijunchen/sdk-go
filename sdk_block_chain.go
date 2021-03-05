package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-pb/config"
	"context"
	"fmt"
)

func (cc ChainClient) CheckNewBlockChainConfig() error {
	cc.logger.Debug("[SDK] begin to send check new block chain config command")
	req := &config.CheckNewBlockChainConfigRequest{}
	client, err := cc.pool.getClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	res, err := client.rpcNode.CheckNewBlockChainConfig(ctx, req)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf("check new block chain config failed, %s", res.Message)
	}
	return nil
}
