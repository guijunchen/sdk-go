/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/pb/common"
	"chainmaker.org/chainmaker-go/pb/consensus"
	"chainmaker.org/chainmaker-go/pb/discovery"
	"encoding/hex"
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
	testSystemContractGetChainInfo(t, client)

	systemChainClient, err := NewChainClient(
		WithChainClientOrgId(orgId1),
		WithChainClientChainId(chainId),
		WithChainClientLogger(getDefaultLogger()),
		WithUserKeyFilePath(userKeyPath),
		WithUserCrtFilePath(userCrtPath),
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
