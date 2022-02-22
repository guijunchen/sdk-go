/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
)

const (
	createContractTimeout = 5
	ContractCName   	  = "CreatorC"
	ContractDName   	  = "beCreated"
	storageVersion        = "1.0.0"
	CBinPath   			  = "../../testdata/inner-create-evm-demo/C.bin"
	CABIPath        	  = "../../testdata/inner-create-evm-demo/C.abi"
	DABIPath        	  = "../../testdata/inner-create-evm-demo/D.abi"

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	testContractCreateContract(sdkConfigOrg1Client1Path)
}

func testContractCreateContract(sdkPath string) {
	fmt.Println("====================== create client ======================")
	client, err := examples.CreateChainClientWithSDKConf(sdkPath)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 创建Creator C合约 ======================")
	usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	testCCreate(client, true, true, usernames...)

	fmt.Println("====================== 调用创建者C合约 ======================")
	testInvokeC(client, true)

	fmt.Println("====================== 调用被创建者D合约 ======================")
	testInvokeD(client, true)
}

func testCCreate(client *sdk.ChainClient, withSyncResult bool, isIgnoreSameContract bool, usernames ...string) {

	codeBytes, err := ioutil.ReadFile(CBinPath)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := createUserContract(client, ContractCName, storageVersion, string(codeBytes), common.RuntimeType_EVM,
		nil, withSyncResult, usernames...)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE EVM factory contract resp: %+v\n", resp)
}

func createUserContract(client *sdk.ChainClient, contractName, version, byteCodePath string, runtime common.RuntimeType,
	kvs []*common.KeyValuePair, withSyncResult bool, usernames ...string) (*common.TxResponse, error) {

	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	//endorsers, err := examples.GetEndorsers(payload, usernames...)
	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		return nil, err
	}

	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func testInvokeC(client *sdk.ChainClient, withSyncResult bool) {

	abiJson, err := ioutil.ReadFile(CABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	dataByte, err := myAbi.Pack("createDSalted", big.NewInt(10000), ContractDName)
	if err != nil {
		log.Fatalln(err)
	}

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	err = invokeUserContract(client, ContractCName, method, "", kvs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func testInvokeD(client *sdk.ChainClient, withSyncResult bool) {

	abiJson, err := ioutil.ReadFile(DABIPath)
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

	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	err = invokeUserContract(client, ContractDName, method, "", kvs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string, kvs []*common.KeyValuePair, withSyncResult bool) error {

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
