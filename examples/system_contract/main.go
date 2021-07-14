/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/hokaccha/go-prettyjson"

	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/pb-go/consensus"
	"chainmaker.org/chainmaker/pb-go/discovery"
	"chainmaker.org/chainmaker/pb-go/store"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
	sdkutils "chainmaker.org/chainmaker/sdk-go/utils"
)

const (
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	testSystemContract()
	testSystemContractArchive()
}

// [系统合约]
func testSystemContract() {
	//client, err := createClientWithConfig()
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	genesisBlockInfo := testSystemContractGetBlockByHeight(client, 1)
	testSystemContractGetTxByTxId(client, genesisBlockInfo.Block.Txs[0].Payload.TxId)
	testSystemContractGetBlockByHash(client, hex.EncodeToString(genesisBlockInfo.Block.Header.BlockHash))
	testSystemContractGetBlockByTxId(client, genesisBlockInfo.Block.Txs[0].Payload.TxId)
	testSystemContractGetLastConfigBlock(client)
	testSystemContractGetLastBlock(client)
	testSystemContractGetChainInfo(client)

	systemChainClient, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	testSystemContractGetNodeChainList(systemChainClient)
}

func testSystemContractArchive() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	var blockHeight uint64 = 4
	fullBlock := testSystemContractGetFullBlockByHeight(client, blockHeight)
	heightByTxId := testSystemContractGetBlockHeightByTxId(client, fullBlock.Block.Txs[0].Payload.TxId)
	if blockHeight != heightByTxId {
		log.Fatalln("blockHeight != heightByTxId")
	}
	heightByHash := testSystemContractGetBlockHeightByHash(client, hex.EncodeToString(fullBlock.Block.Header.BlockHash))
	if blockHeight != heightByHash {
		log.Fatalln("blockHeight != heightByHash")
	}

	testSystemContractGetCurrentBlockHeight(client)
	testSystemContractGetArchivedBlockHeight(client)
	testSystemContractGetBlockHeaderByHeight(client)
}

func testSystemContractGetTxByTxId(client *sdk.ChainClient, txId string) *common.TransactionInfo {
	transactionInfo, err := client.GetTxByTxId(txId)
	if err != nil {
		log.Fatalln(err)
	}
	return transactionInfo
}

func testSystemContractGetBlockByHeight(client *sdk.ChainClient, blockHeight uint64) *common.BlockInfo {
	blockInfo, err := client.GetBlockByHeight(blockHeight, true)
	if err != nil {
		log.Fatalln(err)
	}
	return blockInfo
}

func testSystemContractGetBlockByHash(client *sdk.ChainClient, blockHash string) *common.BlockInfo {
	blockInfo, err := client.GetBlockByHash(blockHash, true)
	if err != nil {
		log.Fatalln(err)
	}
	return blockInfo
}

func testSystemContractGetBlockByTxId(client *sdk.ChainClient, txId string) *common.BlockInfo {
	blockInfo, err := client.GetBlockByTxId(txId, true)
	if err != nil {
		log.Fatalln(err)
	}
	return blockInfo
}

func testSystemContractGetLastConfigBlock(client *sdk.ChainClient) *common.BlockInfo {
	blockInfo, err := client.GetLastConfigBlock(true)
	if err != nil {
		log.Fatalln(err)
	}
	return blockInfo
}

func testSystemContractGetLastBlock(client *sdk.ChainClient) *common.BlockInfo {
	blockInfo, err := client.GetLastBlock(true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("last block height: %d\n", blockInfo.Block.Header.BlockHeight)
	marshal, err := prettyjson.Marshal(blockInfo)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("blockInfo: %s\n", marshal)
	return blockInfo
}

func testSystemContractGetCurrentBlockHeight(client *sdk.ChainClient) uint64 {
	height, err := client.GetCurrentBlockHeight()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("current block height: %d\n", height)
	return height
}

func testSystemContractGetArchivedBlockHeight(client *sdk.ChainClient) uint64 {
	height, err := client.GetArchivedBlockHeight()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("archived block height: %d\n", height)
	return height
}

func testSystemContractGetBlockHeightByTxId(client *sdk.ChainClient, txId string) uint64 {
	height, err := client.GetBlockHeightByTxId(txId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("txId [%s] => block height: %d\n", txId, height)
	return height
}

func testSystemContractGetBlockHeightByHash(client *sdk.ChainClient, blockHash string) uint64 {
	height, err := client.GetBlockHeightByHash(blockHash)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("blockHash [%s] => block height: %d\n", blockHash, height)
	return height
}

func testSystemContractGetChainInfo(client *sdk.ChainClient) *discovery.ChainInfo {
	chainConfig, err := client.GetChainConfig()
	if err != nil {
		log.Fatalln(err)
	}
	chainInfo := &discovery.ChainInfo{}
	if chainConfig.Consensus.Type != consensus.ConsensusType_SOLO {
		var err error
		chainInfo, err = client.GetChainInfo()
		if err != nil {
			log.Fatalln(err)
		}
	}
	return chainInfo
}

func testSystemContractGetNodeChainList(client *sdk.ChainClient) *discovery.ChainList {
	chainList, err := client.GetNodeChainList()
	if err != nil {
		log.Fatalln(err)
	}
	return chainList
}

func testSystemContractGetFullBlockByHeight(client *sdk.ChainClient, blockHeight uint64) *store.BlockWithRWSet {
	fullBlockInfo, err := client.GetFullBlockByHeight(blockHeight)
	if err != nil {
		if sdkutils.IsArchivedString(err.Error()) {
			fmt.Println("Is archived...")
		}
	}
	if err != nil {
		log.Fatalln(err)
	}
	marshal, err := prettyjson.Marshal(fullBlockInfo)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("fullBlockInfo: %s\n", marshal)
	return fullBlockInfo
}

func testSystemContractGetBlockHeaderByHeight(client *sdk.ChainClient) {
	_, err := client.GetBlockHeaderByHeight(0)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = client.GetBlockHeaderByHeight(5)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = client.GetBlockHeaderByHeight(-2)
	if err == nil {
		log.Fatalln("require err not nil")
	}
}
