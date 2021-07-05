/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
)

const (
	createContractTimeout = 5
	storageContractName = "storage001"
	storageVersion      = "1.0.0"
	storageByteCodePath = "../../testdata/storage-evm-demo/storage.bin"
	storageABIPath      = "../../testdata/storage-evm-demo/storage.abi"

	sdkConfigOrg1Admin1Path  = "../sdk_configs/sdk_config_org1_admin1.yml"
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
	sdkConfigOrg2Admin1Path  = "../sdk_configs/sdk_config_org2_admin1.yml"
	sdkConfigOrg3Admin1Path  = "../sdk_configs/sdk_config_org3_admin1.yml"
	sdkConfigOrg4Admin1Path  = "../sdk_configs/sdk_config_org4_admin1.yml"
)

func main() {
	testUserContractStorageEVM()
}

func testUserContractStorageEVM() {
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

	fmt.Println("====================== 创建Storage合约 ======================")
	testUserContractStorageEVMCreate(client, admin1, admin2, admin3, admin4, true, true)

	fmt.Println("====================== 设置数值 ======================")
	testUserContractStorageEVMSet(client, 123, true)

	fmt.Println("====================== 查看数值 ======================")
	testUserContractStorageEVMGet(client, true)
}

func testUserContractStorageEVMCreate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, 
	withSyncResult bool, isIgnoreSameContract bool) {

	codeBytes, err := ioutil.ReadFile(storageByteCodePath)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		storageContractName, storageVersion, string(codeBytes), common.RuntimeType_EVM, nil, withSyncResult)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE EVM storage contract resp: %+v\n", resp)
}

func createUserContract(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, contractName, version,
	byteCodePath string, runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {

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

func testUserContractStorageEVMSet(client *sdk.ChainClient, data int64, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(storageABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	dataByte, err := myAbi.Pack("set", big.NewInt(data))
	if err != nil {
		log.Fatalln(err)
	}

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	pairs := map[string]string{
		"data": dataString,
	}

	err = invokeUserContract(client, storageContractName, method, "", pairs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string, params map[string]string, withSyncResult bool) error {

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

func testUserContractStorageEVMGet(client *sdk.ChainClient, withSyncResult bool) {

	abiJson, err := ioutil.ReadFile(storageABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	dataByte, err := myAbi.Pack("get")
	if err != nil {
		log.Fatalln(err)
	}

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	pairs := map[string]string{
		"data": dataString,
	}

	result, err := invokeUserContractWithResult(client, storageContractName, method, "", pairs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}

	val, err := myAbi.Unpack("get", result)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("val: %d\n", val)
}

func invokeUserContractWithResult(client *sdk.ChainClient, contractName, method, txId string,
	params map[string]string, withSyncResult bool) ([]byte, error) {

	resp, err := client.InvokeContract(contractName, method, txId, params, -1, withSyncResult)
	if err != nil {
		return nil, err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return nil, fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	return resp.ContractResult.Result, nil
}
