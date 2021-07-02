/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/binary"
	"fmt"
	"log"

	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
)

const (
	queryAddr             = "query_address"
	createContractTimeout = 5

	sdkConfigOrg1Admin1Path  = "../sdk_configs/sdk_config_org1_admin1.yml"
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
	sdkConfigOrg2Admin1Path  = "../sdk_configs/sdk_config_org2_admin1.yml"
	sdkConfigOrg2Client1Path = "../sdk_configs/sdk_config_org2_client1.yml"
	sdkConfigOrg3Admin1Path  = "../sdk_configs/sdk_config_org3_admin1.yml"
	sdkConfigOrg4Admin1Path  = "../sdk_configs/sdk_config_org4_admin1.yml"
)

var (
	assetContractName = "asset001"
	assetVersion      = "1.0.0"
	assetByteCodePath = "../../testdata/asset-wasm-demo/rust-asset-management-1.0.0.wasm"
)

func main() {
	testUserContractAsset()
	testUserContractAssetBalanceOf()
}

func testUserContractAssetBalanceOf() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	client2, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg2Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 1)查询钱包地址 ======================")
	addr1 := testUserContractAssetQuery(client, queryAddr, nil)
	fmt.Printf("client1 address: %s\n", addr1)
	addr2 := testUserContractAssetQuery(client2, queryAddr, nil)
	fmt.Printf("client2 address: %s\n", addr2)

	fmt.Println("====================== 2)查询钱包余额 ======================")
	getBalance(client, addr1)
	getBalance(client, addr2)
}

func testUserContractAsset() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	client2, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg2Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	admin1, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Admin1Path)
	if err != nil {
		log.Fatalln(err)
	}
	admin2, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg2Admin1Path)
	if err != nil {
		log.Fatalln(err)
	}
	admin3, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg3Admin1Path)
	if err != nil {
		log.Fatalln(err)
	}
	admin4, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg4Admin1Path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 1)安装钱包合约 ======================")
	pairs := []*common.KeyValuePair{
		{
			Key:   "issue_limit",
			Value: "100000000",
		},
		{
			Key:   "total_supply",
			Value: "100000000",
		},
	}
	testUserContractAssetCreate(client, admin1, admin2, admin3, admin4, pairs, true, false)

	fmt.Println("====================== 2)注册另一个用户 ======================")
	testUserContractAssetInvokeRegister(client2, "register", true)

	fmt.Println("====================== 3)查询钱包地址 ======================")
	addr1 := testUserContractAssetQuery(client, queryAddr, nil)
	fmt.Printf("client1 address: %s\n", addr1)
	addr2 := testUserContractAssetQuery(client2, queryAddr, nil)
	fmt.Printf("client2 address: %s\n", addr2)

	fmt.Println("====================== 4)给用户分别发币100000 ======================")
	amount := "100000"
	testUserContractAssetInvoke(client, "issue_amount", amount, addr1, true)
	testUserContractAssetInvoke(client, "issue_amount", amount, addr2, true)

	fmt.Println("====================== 5)分别查看余额 ======================")
	getBalance(client, addr1)
	getBalance(client, addr2)

	fmt.Println("====================== 6)A给B转账100 ======================")
	amount = "100"
	testUserContractAssetInvoke(client, "transfer", amount, addr2, true)

	fmt.Println("====================== 7)再次分别查看余额 ======================")
	getBalance(client, addr1)
	getBalance(client, addr2)
}

func testUserContractAssetCreate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	kvs []*common.KeyValuePair, withSyncResult bool, isIgnoreSameContract bool) {

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		assetContractName, assetVersion, assetByteCodePath, common.RuntimeType_WASMER, kvs, withSyncResult)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE asset contract resp: %+v\n", resp)
}

func testUserContractAssetInvokeRegister(client *sdk.ChainClient, method string, withSyncResult bool) {
	err := invokeUserContract(client, assetContractName, method, "", nil, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func testUserContractAssetQuery(client *sdk.ChainClient, method string, params map[string]string) string {
	resp, err := client.QueryContract(assetContractName, method, params, -1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("QUERY asset contract [%s] resp: %+v\n", method, resp)

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		log.Fatalln(err)
	}
	return string(resp.ContractResult.Result)
}

func testUserContractAssetInvoke(client *sdk.ChainClient, method string, amount, addr string, withSyncResult bool) {
	params := map[string]string{
		"amount": amount,
		"to":     addr,
	}
	err := invokeUserContract(client, assetContractName, method, "", params, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func getBalance(client *sdk.ChainClient, addr string) {
	params := map[string]string{
		"owner": addr,
	}
	balance := testUserContractAssetQuery(client, "balance_of", params)
	val, err := sdk.BytesToInt([]byte(balance), binary.LittleEndian)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("client [%s] balance: %d\n", addr, val)
}

func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string,
	params map[string]string, withSyncResult bool) error {

	resp, err := client.InvokeContract(contractName, method, txId, params, -1, withSyncResult)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if !withSyncResult {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	} else {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	}

	return nil
}

func createUserContract(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, contractName, version,
	byteCodePath string, runtime common.RuntimeType, kvs []*common.KeyValuePair,
	withSyncResult bool) (*common.TxResponse, error) {

	payloadBytes, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	if err != nil {
		return nil, err
	}

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
