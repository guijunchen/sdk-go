package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/store"
)

func (cc *ChainClient) ArchiveBlock(targetBlockHeight int64) (*common.TxResponse, error) {
	panic("implement me")
}

func (cc *ChainClient) RestoreBlocks(startBlockHeight int64) (*common.TxResponse, error) {
	panic("implement me")
}

func (cc *ChainClient) GetArchivedTxByTxId(txId string) (*common.TransactionInfo, error) {
	panic("implement me")
}

func (cc *ChainClient) GetArchivedBlockByHeight(blockHeight int64, withRWSet bool) (*common.BlockInfo, error) {
	panic("implement me")
}

func (cc *ChainClient) GetArchivedFullBlockByHeight(blockHeight int64) (*store.BlockWithRWSet, error) {
	panic("implement me")
}

func (cc *ChainClient) GetArchivedBlockByHash(blockHash string, withRWSet bool) (*common.BlockInfo, error) {
	panic("implement me")
}

func (cc *ChainClient) GetArchivedBlockByTxId(txId string, withRWSet bool) (*common.BlockInfo, error) {
	panic("implement me")
}

