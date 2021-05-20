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
	//client, err := createClientWithConfig()
	client, err := createClient()
	require.Nil(t, err)

	//blockInfo := testSystemContractGetBlockByHeight(t, client, -1)
	//testSystemContractGetTxByTxId(t, client, blockInfo.Block.Txs[0].Header.TxId)
	//testSystemContractGetBlockByHash(t, client, hex.EncodeToString(blockInfo.Block.Header.BlockHash))
	//testSystemContractGetBlockByTxId(t, client, blockInfo.Block.Txs[0].Header.TxId)
	//testSystemContractGetLastConfigBlock(t, client)
	//testSystemContractGetLastBlock(t, client)
	//testSystemContractGetChainInfo(t, client)
	//
	//systemChainClient, err := NewChainClient(
	//	WithChainClientOrgId(orgId1),
	//	WithChainClientChainId(chainId),
	//	WithChainClientLogger(getDefaultLogger()),
	//	WithUserKeyFilePath(fmt.Sprintf(userKeyPath, orgId1)),
	//	WithUserCrtFilePath(fmt.Sprintf(userCrtPath, orgId1)),
	//	AddChainClientNodeConfig(node1),
	//	AddChainClientNodeConfig(node2),
	//)
	//require.Nil(t, err)
	//
	//testSystemContractGetNodeChainList(t, systemChainClient)

	// Archive test
	var blockHeight int64 = 4
	fullBlock := testSystemContractGetFullBlockByHeight(t, client, blockHeight)
	heightByTxId := testSystemContractGetBlockHeightByTxId(t, client, fullBlock.Block.Txs[0].Header.TxId)
	require.Equal(t, blockHeight, heightByTxId)
	heightByHash := testSystemContractGetBlockHeightByHash(t, client, hex.EncodeToString(fullBlock.Block.Header.BlockHash))
	require.Equal(t, blockHeight, heightByHash)

	testSystemContractGetCurrentBlockHeight(t, client)
	testSystemContractGetArchivedBlockHeight(t, client)
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

func testSystemContractGetArchivedBlockHeight(t *testing.T, client *ChainClient) int64 {
	height, err := client.GetArchivedBlockHeight()
	require.Nil(t, err)
	fmt.Printf("archived block height: %d\n", height)
	return height
}

func testSystemContractGetBlockHeightByTxId(t *testing.T, client *ChainClient, txId string) int64 {
	height, err := client.GetBlockHeightByTxId(txId)
	require.Nil(t, err)
	fmt.Printf("txId [%s] => block height: %d\n", txId, height)
	return height
}

func testSystemContractGetBlockHeightByHash(t *testing.T, client *ChainClient, blockHash string) int64 {
	height, err := client.GetBlockHeightByHash(blockHash)
	require.Nil(t, err)
	fmt.Printf("blockHash [%s] => block height: %d\n", blockHash, height)
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
	if err != nil {
		if IsArchivedString(err.Error()) {
			fmt.Println("Is archived...")
		}
	}
	require.Nil(t, err)
	marshal, err := prettyjson.Marshal(fullBlockInfo)
	require.Nil(t, err)
	fmt.Printf("fullBlockInfo: %s\n", marshal)
	return fullBlockInfo
}
