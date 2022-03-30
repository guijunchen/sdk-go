/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/evmutils"
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
	CreateContractTimeout = 5
	contractVersion       = "1.0.0"

	calleeName   	      = "callee"
	calleeBin   	  	  = "../../testdata/cross-call-evm-demo/Callee.bin"
	callerName            = "caller"
	callerBin             = "../../testdata/cross-call-evm-demo/Caller.bin"
	callerABI             = "../../testdata/cross-call-evm-demo/Caller.abi"

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	testCrossCall(sdkConfigOrg1Client1Path)
}

func testCrossCall(sdkPath string) {
	fmt.Println("====================== create client ======================")
	client, err := examples.CreateChainClientWithSDKConf(sdkPath)
	if err != nil {
		log.Fatalln(err)
	}

	usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1, examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	fmt.Println("====================== 创建callee(被调用者)合约 ======================")
	testCreateCallee(client, true, true, usernames...)

	fmt.Println("====================== 创建caller(调用者)合约 ======================")
	testCreateCaller(client, true, true, usernames...)

	fmt.Println("====================== caller跨合约调用callee的Adder方法 ======================")
	testCallerCrossCallCalleeAdder(client, true)
}

func testCreateCallee(client *sdk.ChainClient, withSyncResult bool, isIgnoreSameContract bool, usernames ...string) {

	codeBytes, err := ioutil.ReadFile(calleeBin)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := createUserContract(client, calleeName, contractVersion, string(codeBytes), common.RuntimeType_EVM,
		nil, withSyncResult, usernames...)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE EVM callee contract resp: [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	fmt.Printf("contract result: [code:%d]/[msg:%s]/[contractResult:%+X]\n",  resp.ContractResult.Code,
		resp.ContractResult.Message, resp.ContractResult.Result)
}

func testCreateCaller(client *sdk.ChainClient, withSyncResult bool, isIgnoreSameContract bool, usernames ...string) {

	codeBytes, err := ioutil.ReadFile(callerBin)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := createUserContract(client, callerName, contractVersion, string(codeBytes), common.RuntimeType_EVM,
		nil, withSyncResult, usernames...)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE EVM caller contract resp: [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	fmt.Printf("contract result: [code:%d]/[msg:%s]/[contractResult:%+X]\n",  resp.ContractResult.Code,
		resp.ContractResult.Message, resp.ContractResult.Result)
}

func testCallerCrossCallCalleeAdder(client *sdk.ChainClient, withSyncResult bool) {

	abiJson, err := ioutil.ReadFile(callerABI)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	addr := evmutils.MakeAddress([]byte(calleeName))
	callee := evmutils.BigToAddress(addr)
	dataByte, err := myAbi.Pack("crossCall", callee, big.NewInt(40), big.NewInt(60))
	if err != nil {
		log.Fatalln(err)
	}

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	kvs := []*common.KeyValuePair{
		{
			Key:   "data",//protocol.ContractEvmParamKey
			Value: []byte(dataString),
		},
	}

	err = invokeUserContract(client, callerName, method, "", kvs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
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

	resp, err := client.SendContractManageRequest(payload, endorsers, CreateContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
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
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message,
			resp.ContractResult.Result)
	} else {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
		fmt.Printf("contract result: [code:%d]/[msg:%s]/[contractResult:%+X]\n",  resp.ContractResult.Code,
			resp.ContractResult.Message, resp.ContractResult.Result)
	}

	return nil
}
