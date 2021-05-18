package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/store"
	"fmt"
	"github.com/golang/protobuf/proto"
)

func (cc *ChainClient) CreateArchiveBlockPayload(targetBlockHeight int64) ([]byte, error) {
	pairs := []*common.KeyValuePair{
		{
			Key:   "targetBlockHeight",
			Value: fmt.Sprintf("%d", targetBlockHeight),
		},
	}

	return cc.CreateArchivePayload(common.ArchiveStoreContractFunction_ARCHIVE_BLOCK.String(), pairs)
}

func (cc *ChainClient) CreateRestoreBlocksPayload(startBlockHeight int64) ([]byte, error) {
	panic("implement me")
}

func (cc *ChainClient) CreateArchivePayload(method string, kvs []*common.KeyValuePair) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [Archive] to be signed payload")

	payload := &common.SystemContractPayload{
		ChainId:      cc.chainId,
		ContractName: common.ContractName_SYSTEM_CONTRACT_ARCHIVE_STORE.String(),
		Method:       method,
		Parameters:   kvs,
	}

	bytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("construct archive payload failed, %s", err)
	}

	return bytes, nil
}

func (cc *ChainClient) SignArchivePayload(payloadBytes []byte) ([]byte, error) {
	return cc.signSystemContractPayload(payloadBytes)
}

func (cc *ChainClient) MergeArchivePayload(signedPayloadBytes [][]byte) ([]byte, error) {
	return mergeContractManageSignedPayload(signedPayloadBytes)
}

func (cc *ChainClient) SendArchiveBlockRequest(mergeSignedPayloadBytes []byte, timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	return cc.sendContractRequest(common.TxType_INVOKE_SYSTEM_CONTRACT, mergeSignedPayloadBytes, timeout, withSyncResult)
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

