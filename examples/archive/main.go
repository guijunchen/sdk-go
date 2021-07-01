/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hokaccha/go-prettyjson"

	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
)

const (
	sdkConfigOrg1Admin1Path  = "../sdk_configs/sdk_config_org1_admin1.yml"
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	testArchive()
	testRestore()
	testGetFromArchiveStore()
}

func testArchive() {
	admin1, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Admin1Path)
	if err != nil {
		panic(err)
	}

	fmt.Println("====================== 数据归档 ======================")
	var targetBlockHeight int64 = 20
	testArchiveBlock(admin1, targetBlockHeight)
}

func testRestore() {
	admin1, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Admin1Path)
	if err != nil {
		panic(err)
	}

	client1, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		panic(err)
	}

	fmt.Println("====================== 归档恢复 ======================")
	var blockHeight int64 = 2

	fullBlock, err := client1.GetArchivedFullBlockByHeight(blockHeight)
	if err != nil {
		panic(err)
	}
	prettyJsonShow("GetArchivedFullBlockByHeight fullBlock", fullBlock)

	fullBlockBytes, err := proto.Marshal(fullBlock)
	if err != nil {
		panic(err)
	}

	testRestoreBlock(admin1, fullBlockBytes)
}

func testGetFromArchiveStore() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		panic(err)
	}

	fmt.Println("====================== 归档查询 ======================")
	var blockHeight int64 = 8
	fullBlockInfo, err := client.GetFromArchiveStore(blockHeight)
	if err != nil {
		panic(err)
	}
	prettyJsonShow("GetFromArchiveStore fullBlockInfo", fullBlockInfo)

	fullBlockInfo, err = client.GetArchivedFullBlockByHeight(blockHeight)
	if err != nil {
		panic(err)
	}
	prettyJsonShow("GetArchivedFullBlockByHeight fullBlockInfo", fullBlockInfo)

	blockInfo, err := client.GetArchivedBlockByHeight(blockHeight, true)
	if err != nil {
		panic(err)
	}
	prettyJsonShow("GetArchivedBlockByHeight with rwset", blockInfo)

	blockInfo, err = client.GetArchivedBlockByHeight(blockHeight, false)
	if err != nil {
		panic(err)
	}
	prettyJsonShow("GetArchivedBlockByHeight without rwset", blockInfo)

	blockInfo, err = client.GetArchivedBlockByHash(hex.EncodeToString(blockInfo.Block.Header.BlockHash), true)
	if err != nil {
		panic(err)
	}
	prettyJsonShow("GetArchivedBlockByHash with rwset", blockInfo)

	blockInfo, err = client.GetArchivedBlockByHash(hex.EncodeToString(blockInfo.Block.Header.BlockHash), false)
	if err != nil {
		panic(err)
	}
	prettyJsonShow("GetArchivedBlockByHash without rwset", blockInfo)

	txId := blockInfo.Block.Txs[0].Header.TxId
	txInfo, err := client.GetArchivedTxByTxId(txId)
	if err != nil {
		panic(err)
	}
	prettyJsonShow("GetArchivedTxByTxId", txInfo)
}

func prettyJsonShow(name string, v interface{}) {
	marshal, err := prettyjson.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n\n====== %s ======\n%s\n==========================\n", name, marshal)
}

func testArchiveBlock(admin1 *sdk.ChainClient, targetBlockHeight int64) {
	var (
		err                error
		payload            []byte
		signedPayloadBytes []byte
		resp               *common.TxResponse
		result             string
	)

	payload, err = admin1.CreateArchiveBlockPayload(targetBlockHeight)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes, err = admin1.SignArchivePayload(payload)
	if err != nil {
		panic(err)
	}

	resp, err = admin1.SendArchiveBlockRequest(signedPayloadBytes, -1)
	if err != nil {
		panic(err)
	}

	err = examples.CheckProposalRequestResp(resp, false)
	if err != nil {
		panic(err)
	}

	result = string(resp.ContractResult.Result)

	fmt.Printf("resp: %+v, result:%s\n", resp, result)
}

func testRestoreBlock(admin1 *sdk.ChainClient, fullBlock []byte) {
	var (
		err                error
		payload            []byte
		signedPayloadBytes []byte
		resp               *common.TxResponse
		result             string
	)

	payload, err = admin1.CreateRestoreBlockPayload(fullBlock)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes, err = admin1.SignArchivePayload(payload)
	if err != nil {
		panic(err)
	}

	resp, err = admin1.SendRestoreBlockRequest(signedPayloadBytes, -1)
	if err != nil {
		panic(err)
	}

	err = examples.CheckProposalRequestResp(resp, false)
	if err != nil {
		panic(err)
	}

	result = string(resp.ContractResult.Result)

	fmt.Printf("resp: %+v, result:%s\n", resp, result)
}
