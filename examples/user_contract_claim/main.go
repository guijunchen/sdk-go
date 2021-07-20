/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"
	"time"

	"chainmaker.org/chainmaker/common/random/uuid"
	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
)

const (
	createContractTimeout = 5
	claimContractName     = "claim001"
	claimVersion          = "2.0.0"
	claimByteCodePath     = "../../testdata/claim-wasm-demo/rust-fact-2.0.0.wasm"

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"

	sdkConfigOrg1Admin1Path  = "../sdk_configs/sdk_config_org1_admin1.yml"
	sdkConfigOrg2Admin1Path  = "../sdk_configs/sdk_config_org2_admin1.yml"
	sdkConfigOrg3Admin1Path  = "../sdk_configs/sdk_config_org3_admin1.yml"
	sdkConfigOrg4Admin1Path  = "../sdk_configs/sdk_config_org4_admin1.yml"
)

func main() {
	testUserContractClaim()
}

func testUserContractClaim() {
	fmt.Println("====================== create client ======================")
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
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

	fmt.Println("====================== 创建合约 ======================")
	testUserContractClaimCreate(client, admin1, admin2, admin3, admin4, true, true)

	fmt.Println("====================== 调用合约 ======================")
	fileHash, err := testUserContractClaimInvoke(client, "save", true)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 执行合约查询接口 ======================")
	//txId := "1cbdbe6106cc4132b464185ea8275d0a53c0261b7b1a470fb0c3f10bd4a57ba6"
	//fileHash = txId[len(txId)/2:]
	kvs := []*common.KeyValuePair {
		{
			Key: "file_hash",
			Value: []byte(fileHash),
		},
	}
	testUserContractClaimQuery(client, "find_by_file_hash", kvs)
}

func testUserContractClaimCreate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	withSyncResult bool, isIgnoreSameContract bool) {

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		claimContractName, claimVersion, claimByteCodePath, common.RuntimeType_WASMER, []*common.KeyValuePair{}, withSyncResult)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE claim contract resp: %+v\n", resp)
}

func createUserContract(client *sdk.ChainClient, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	contractName, version, byteCodePath string, runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {

	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	// 各组织Admin权限用户签名
	signedPayload1, err := admin1.SignContractManagePayload(payload)
	if err != nil {
		return nil, err
	}

	signedPayload2, err := admin2.SignContractManagePayload(payload)
	if err != nil {
		return nil, err
	}

	signedPayload3, err := admin3.SignContractManagePayload(payload)
	if err != nil {
		return nil, err
	}

	signedPayload4, err := admin4.SignContractManagePayload(payload)
	if err != nil {
		return nil, err
	}

	var endosers []*common.EndorsementEntry
	endosers = append(endosers, signedPayload1)
	endosers = append(endosers, signedPayload2)
	endosers = append(endosers, signedPayload3)
	endosers = append(endosers, signedPayload4)

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(payload, endosers, createContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func testUserContractClaimInvoke(client *sdk.ChainClient,
	method string, withSyncResult bool) (string, error) {

	//curTime := fmt.Sprintf("%d", CurrentTimeMillisSeconds())
	curTime := time.Now().Format("2006-01-02 15:04:05")

	fileHash := uuid.GetUUID()
	kvs := []*common.KeyValuePair {
		{
			Key: "time",
			Value: []byte(curTime),
		},
		{
			Key: "file_hash",
			Value: []byte(fileHash),
		},
		{
			Key: "file_name",
			Value: []byte(fmt.Sprintf("file_%s", curTime)),
		},
	}

	err := invokeUserContract(client, claimContractName, method, "", kvs, withSyncResult)
	if err != nil {
		return "", err
	}

	return fileHash, nil
}

func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string,
	kvs []*common.KeyValuePair, withSyncResult bool) error {

	resp, err := client.InvokeContract(contractName, method, txId, kvs, -1, withSyncResult)
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

func testUserContractClaimQuery(client *sdk.ChainClient, method string, kvs []*common.KeyValuePair) {
	resp, err := client.QueryContract(claimContractName, method, kvs, -1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("QUERY claim contract resp: %+v\n", resp)
}
