/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"database/sql"
	"fmt"
	"strings"

	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/pb-go/store"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gogo/protobuf/proto"
)

const (
	mysqlDBNamePrefix     = "cm_archived_chain"
	mysqlTableNamePrefix  = "t_block_info"
	rowsPerBlockInfoTable = 100000
)

func (cc *ChainClient) CreateArchiveBlockPayload(targetBlockHeight int64) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [Archive] to be signed payload")

	payload := &common.ArchiveBlockPayload{
		BlockHeight: targetBlockHeight,
	}

	bytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("construct archive payload failed, %s", err)
	}

	return bytes, nil
}

func (cc *ChainClient) CreateRestoreBlockPayload(fullBlock []byte) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [restore] to be signed payload")

	payload := &common.RestoreBlockPayload{
		FullBlock: fullBlock,
	}

	bytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("construct archive payload failed, %s", err)
	}

	return bytes, nil
}

func (cc *ChainClient) SignArchivePayload(payloadBytes []byte) ([]byte, error) {
	//return cc.signSystemContractPayload(payloadBytes)
	return payloadBytes, nil
}

func (cc *ChainClient) SendArchiveBlockRequest(mergeSignedPayloadBytes []byte, timeout int64) (*common.TxResponse, error) {
	return cc.sendContractRequest(common.TxType_ARCHIVE_FULL_BLOCK, mergeSignedPayloadBytes, timeout, false)
}

func (cc *ChainClient) SendRestoreBlockRequest(mergeSignedPayloadBytes []byte, timeout int64) (*common.TxResponse, error) {
	return cc.sendContractRequest(common.TxType_RESTORE_FULL_BLOCK, mergeSignedPayloadBytes, timeout, false)
}

func (cc *ChainClient) GetArchivedFullBlockByHeight(blockHeight int64) (*store.BlockWithRWSet, error) {
	fullBlock, err := cc.GetFromArchiveStore(blockHeight)
	if err != nil {
		return nil, err
	}

	return fullBlock, nil
}

func (cc *ChainClient) GetArchivedBlockByHeight(blockHeight int64, withRWSet bool) (*common.BlockInfo, error) {
	fullBlock, err := cc.GetFromArchiveStore(blockHeight)
	if err != nil {
		return nil, err
	}

	blockInfo := &common.BlockInfo{
		Block: fullBlock.Block,
	}

	if withRWSet {
		blockInfo.RwsetList = fullBlock.TxRWSets
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetArchivedBlockByTxId(txId string, withRWSet bool) (*common.BlockInfo, error) {
	blockHeight, err := cc.GetBlockHeightByTxId(txId)
	if err != nil {
		return nil, err
	}

	return cc.GetArchivedBlockByHeight(blockHeight, withRWSet)
}

func (cc *ChainClient) GetArchivedBlockByHash(blockHash string, withRWSet bool) (*common.BlockInfo, error) {
	blockHeight, err := cc.GetBlockHeightByHash(blockHash)
	if err != nil {
		return nil, err
	}

	return cc.GetArchivedBlockByHeight(blockHeight, withRWSet)
}

func (cc *ChainClient) GetArchivedTxByTxId(txId string) (*common.TransactionInfo, error) {
	blockHeight, err := cc.GetBlockHeightByTxId(txId)
	if err != nil {
		return nil, err
	}

	blockInfo, err := cc.GetArchivedBlockByHeight(blockHeight, false)
	if err != nil {
		return nil, err
	}

	for idx, tx := range blockInfo.Block.Txs {
		if tx.Header.TxId == txId {
			return &common.TransactionInfo{
				Transaction: tx,
				BlockHeight: uint64(blockInfo.Block.Header.BlockHeight),
				BlockHash:   blockInfo.Block.Header.BlockHash,
				TxIndex:     uint32(idx),
			}, nil
		}
	}

	return nil, fmt.Errorf("CANNOT BE HERE! unknown tx [%s] in archive block [%d]", txId, blockHeight)
}

func (cc *ChainClient) GetFromArchiveStore(blockHeight int64) (*store.BlockWithRWSet, error) {
	archiveType := Config.ChainClientConfig.ArchiveConfig.Type
	if archiveType == "mysql" {
		return cc.GetArchivedBlockFromMySQL(blockHeight)
	}

	return nil, fmt.Errorf("unsupport archive type [%s]", archiveType)
}

func (cc *ChainClient) GetArchivedBlockFromMySQL(blockHeight int64) (*store.BlockWithRWSet, error) {

	var (
		blockWithRWSetBytes []byte
		hmac                string
		blockWithRWSet      store.BlockWithRWSet
	)

	dest := Config.ChainClientConfig.ArchiveConfig.Dest
	destList := strings.Split(dest, ":")
	if len(destList) != 4 {
		return nil, fmt.Errorf("invalid archive dest")
	}

	user, pwd, host, port := destList[0], destList[1], destList[2], destList[3]
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s_%s?charset=utf8mb4",
		user, pwd, host, port, mysqlDBNamePrefix, cc.chainId))
	if err != nil {
		return nil, fmt.Errorf("mysql init failed, %s", err.Error())
	}
	defer db.Close()

	err = db.QueryRow(fmt.Sprintf("SELECT Fblock_with_rwset, Fhmac from %s_%d WHERE Fblock_height=?",
		mysqlTableNamePrefix, blockHeight/rowsPerBlockInfoTable+1), blockHeight).Scan(&blockWithRWSetBytes, &hmac)
	if err != nil {
		return nil, fmt.Errorf("select from mysql failed, %s", err.Error())
	}

	// TODO: hmac校验

	err = proto.Unmarshal(blockWithRWSetBytes, &blockWithRWSet)
	if err != nil {
		return nil, fmt.Errorf("unmarshal store.BlockWithRWSet failed, %s", err.Error())
	}

	return &blockWithRWSet, nil
}
