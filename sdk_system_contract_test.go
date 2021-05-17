/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/consensus"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/discovery"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/store"
	"encoding/hex"
	"fmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/stretchr/testify/require"
	"testing"
)

// [系统合约]
func TestSystemContract(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	blockInfo := testSystemContractGetBlockByHeight(t, client, -1)
	testSystemContractGetTxByTxId(t, client, blockInfo.Block.Txs[0].Header.TxId)
	testSystemContractGetBlockByHash(t, client, hex.EncodeToString(blockInfo.Block.Header.BlockHash))
	testSystemContractGetBlockByTxId(t, client, blockInfo.Block.Txs[0].Header.TxId)
	testSystemContractGetLastConfigBlock(t, client)
	testSystemContractGetLastBlock(t, client)
	//testSystemContractGetFullBlockByHeight(t, client, 10)
	testSystemContractGetCurrentBlockHeight(t, client)
	testSystemContractGetChainInfo(t, client)

	systemChainClient, err := NewChainClient(
		WithChainClientOrgId(orgId1),
		WithChainClientChainId(chainId),
		WithChainClientLogger(getDefaultLogger()),
		WithUserKeyFilePath(fmt.Sprintf(userKeyPath, orgId1)),
		WithUserCrtFilePath(fmt.Sprintf(userCrtPath, orgId1)),
		AddChainClientNodeConfig(node1),
		AddChainClientNodeConfig(node2),
	)
	require.Nil(t, err)

	testSystemContractGetNodeChainList(t, systemChainClient)
}

func testSystemContractGetTxByTxId(t *testing.T, client *ChainClient, txId string) *common.TransactionInfo {
	transactionInfo, err := client.GetTxByTxId(txId)
	require.Nil(t, err)
	return transactionInfo
}

func testSystemContractGetBlockByHeight(t *testing.T, client *ChainClient, blockHeight int64) *common.BlockInfo {
	blockInfo, err := client.GetBlockByHeight(blockHeight, true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGetBlockByHash(t *testing.T, client *ChainClient, blockHash string) *common.BlockInfo {
	blockInfo, err := client.GetBlockByHash(blockHash, true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGetBlockByTxId(t *testing.T, client *ChainClient, txId string) *common.BlockInfo {
	blockInfo, err := client.GetBlockByTxId(txId, true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGetLastConfigBlock(t *testing.T, client *ChainClient) *common.BlockInfo {
	blockInfo, err := client.GetLastConfigBlock(true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGetLastBlock(t *testing.T, client *ChainClient) *common.BlockInfo {
	blockInfo, err := client.GetLastBlock(true)
	require.Nil(t, err)
	fmt.Printf("last block height: %d\n", blockInfo.Block.Header.BlockHeight)
	marshal, err := prettyjson.Marshal(blockInfo)
	require.Nil(t, err)
	fmt.Printf("blockInfo: %s\n", marshal)
	return blockInfo
}

func testSystemContractGetCurrentBlockHeight(t *testing.T, client *ChainClient) int64 {
	height, err := client.GetCurrentBlockHeight()
	require.Nil(t, err)
	fmt.Printf("current block height: %d\n", height)
	return height
}

func testSystemContractGetChainInfo(t *testing.T, client *ChainClient) *discovery.ChainInfo {
	chainConfig := testGetChainConfig(t, client)
	chainInfo := &discovery.ChainInfo{}
	if chainConfig.Consensus.Type != consensus.ConsensusType_SOLO {
		var err error
		chainInfo, err = client.GetChainInfo()
		require.Nil(t, err)
	}
	return chainInfo
}

func testSystemContractGetNodeChainList(t *testing.T, client *ChainClient) *discovery.ChainList {
	chainList, err := client.GetNodeChainList()
	require.Nil(t, err)
	return chainList
}

func testSystemContractGetFullBlockByHeight(t *testing.T, client *ChainClient, blockHeight int64) *store.BlockWithRWSet {
	fullBlockInfo, err := client.GetFullBlockByHeight(blockHeight)
	require.Nil(t, err)
	marshal, err := prettyjson.Marshal(fullBlockInfo)
	require.Nil(t, err)
	fmt.Printf("fullBlockInfo: %s\n", marshal)
	return fullBlockInfo
}
