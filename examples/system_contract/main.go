/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"fmt"

	"github.com/hokaccha/go-prettyjson"

	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/pb-go/consensus"
	"chainmaker.org/chainmaker/pb-go/discovery"
	"chainmaker.org/chainmaker/pb-go/store"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
)

func main() {
	testSystemContract()
	testSystemContractArchive()
}

// [系统合约]
func testSystemContract() {
	//client, err := createClientWithConfig()
	client, err := examples.CreateClient()
	if err != nil {
		panic(err)
	}

	blockInfo := testSystemContractGetBlockByHeight(client, -1)
	testSystemContractGetTxByTxId(client, blockInfo.Block.Txs[0].Header.TxId)
	testSystemContractGetBlockByHash(client, hex.EncodeToString(blockInfo.Block.Header.BlockHash))
	testSystemContractGetBlockByTxId(client, blockInfo.Block.Txs[0].Header.TxId)
	testSystemContractGetLastConfigBlock(client)
	testSystemContractGetLastBlock(client)
	testSystemContractGetChainInfo(client)

	systemChainClient, err := examples.CreateClient()
	if err != nil {
		panic(err)
	}

	testSystemContractGetNodeChainList(systemChainClient)
}

func testSystemContractArchive() {
	client, err := examples.CreateClient()
	if err != nil {
		panic(err)
	}

	var blockHeight int64 = 4
	fullBlock := testSystemContractGetFullBlockByHeight(client, blockHeight)
	heightByTxId := testSystemContractGetBlockHeightByTxId(client, fullBlock.Block.Txs[0].Header.TxId)
	if blockHeight != heightByTxId {
		panic("blockHeight != heightByTxId")
	}
	heightByHash := testSystemContractGetBlockHeightByHash(client, hex.EncodeToString(fullBlock.Block.Header.BlockHash))
	if blockHeight != heightByHash {
		panic("blockHeight != heightByHash")
	}

	testSystemContractGetCurrentBlockHeight(client)
	testSystemContractGetArchivedBlockHeight(client)
	testSystemContractGetBlockHeaderByHeight(client)
}

func testSystemContractGetTxByTxId(client *sdk.ChainClient, txId string) *common.TransactionInfo {
	transactionInfo, err := client.GetTxByTxId(txId)
	if err != nil {
		panic(err)
	}
	return transactionInfo
}

func testSystemContractGetBlockByHeight(client *sdk.ChainClient, blockHeight int64) *common.BlockInfo {
	blockInfo, err := client.GetBlockByHeight(blockHeight, true)
	if err != nil {
		panic(err)
	}
	return blockInfo
}

func testSystemContractGetBlockByHash(client *sdk.ChainClient, blockHash string) *common.BlockInfo {
	blockInfo, err := client.GetBlockByHash(blockHash, true)
	if err != nil {
		panic(err)
	}
	return blockInfo
}

func testSystemContractGetBlockByTxId(client *sdk.ChainClient, txId string) *common.BlockInfo {
	blockInfo, err := client.GetBlockByTxId(txId, true)
	if err != nil {
		panic(err)
	}
	return blockInfo
}

func testSystemContractGetLastConfigBlock(client *sdk.ChainClient) *common.BlockInfo {
	blockInfo, err := client.GetLastConfigBlock(true)
	if err != nil {
		panic(err)
	}
	return blockInfo
}

func testSystemContractGetLastBlock(client *sdk.ChainClient) *common.BlockInfo {
	blockInfo, err := client.GetLastBlock(true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("last block height: %d\n", blockInfo.Block.Header.BlockHeight)
	marshal, err := prettyjson.Marshal(blockInfo)
	if err != nil {
		panic(err)
	}
	fmt.Printf("blockInfo: %s\n", marshal)
	return blockInfo
}

func testSystemContractGetCurrentBlockHeight(client *sdk.ChainClient) int64 {
	height, err := client.GetCurrentBlockHeight()
	if err != nil {
		panic(err)
	}
	fmt.Printf("current block height: %d\n", height)
	return height
}

func testSystemContractGetArchivedBlockHeight(client *sdk.ChainClient) int64 {
	height, err := client.GetArchivedBlockHeight()
	if err != nil {
		panic(err)
	}
	fmt.Printf("archived block height: %d\n", height)
	return height
}

func testSystemContractGetBlockHeightByTxId(client *sdk.ChainClient, txId string) int64 {
	height, err := client.GetBlockHeightByTxId(txId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("txId [%s] => block height: %d\n", txId, height)
	return height
}

func testSystemContractGetBlockHeightByHash(client *sdk.ChainClient, blockHash string) int64 {
	height, err := client.GetBlockHeightByHash(blockHash)
	if err != nil {
		panic(err)
	}
	fmt.Printf("blockHash [%s] => block height: %d\n", blockHash, height)
	return height
}

func testSystemContractGetChainInfo(client *sdk.ChainClient) *discovery.ChainInfo {
	chainConfig, err := client.GetChainConfig()
	if err != nil {
		panic(err)
	}
	chainInfo := &discovery.ChainInfo{}
	if chainConfig.Consensus.Type != consensus.ConsensusType_SOLO {
		var err error
		chainInfo, err = client.GetChainInfo()
		if err != nil {
			panic(err)
		}
	}
	return chainInfo
}

func testSystemContractGetNodeChainList(client *sdk.ChainClient) *discovery.ChainList {
	chainList, err := client.GetNodeChainList()
	if err != nil {
		panic(err)
	}
	return chainList
}

func testSystemContractGetFullBlockByHeight(client *sdk.ChainClient, blockHeight int64) *store.BlockWithRWSet {
	fullBlockInfo, err := client.GetFullBlockByHeight(blockHeight)
	if err != nil {
		if sdk.IsArchivedString(err.Error()) {
			fmt.Println("Is archived...")
		}
	}
	if err != nil {
		panic(err)
	}
	marshal, err := prettyjson.Marshal(fullBlockInfo)
	if err != nil {
		panic(err)
	}
	fmt.Printf("fullBlockInfo: %s\n", marshal)
	return fullBlockInfo
}

func testSystemContractGetBlockHeaderByHeight(client *sdk.ChainClient) {
	_, err := client.GetBlockHeaderByHeight(0)
	if err != nil {
		panic(err)
	}

	_, err = client.GetBlockHeaderByHeight(5)
	if err != nil {
		panic(err)
	}

	_, err = client.GetBlockHeaderByHeight(-2)
	if err == nil {
		panic("require err not nil")
	}
}
