package chainmaker_sdk_go

import "chainmaker.org/chainmaker-sdk-go/pb/protogo/common"

func (cc *ChainClient) ArchiveBlock(targetBlockHeight int64) (*common.TxResponse, error) {
	panic("implement me")
}

func (cc *ChainClient) RestoreBlocks(startBlockHeight int64) (*common.TxResponse, error) {
	panic("implement me")
}

