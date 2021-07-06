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

	"chainmaker.org/chainmaker/common/evmutils"
	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
)

const (
	createContractTimeout = 5
	tokenContractName     = "token001"
	tokenVersion          = "1.0.0"
	tokenByteCodePath     = "../../testdata/token-evm-demo/token.bin"
	tokenABIPath          = "../../testdata/token-evm-demo/token.abi"

	client1AddrInt = "1087848554046178479107522336262214072175637027873"
	client2AddrInt = "944104665674401770091203869615921096651560803325"
	client1AddrSki = "7081212378e72d4ecf406c30384f82a74c2a0c8d9e91ccfa94c245023942240f"
	amount         = 200

	sdkConfigOrg1Admin1Path  = "../sdk_configs/sdk_config_org1_admin1.yml"
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
	sdkConfigOrg2Admin1Path  = "../sdk_configs/sdk_config_org2_admin1.yml"
	sdkConfigOrg3Admin1Path  = "../sdk_configs/sdk_config_org3_admin1.yml"
	sdkConfigOrg4Admin1Path  = "../sdk_configs/sdk_config_org4_admin1.yml"
)

func main() {
	testUserContractTokenEVM()
}

func testUserContractTokenEVM() {
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

	fmt.Println("====================== 创建Token合约,给client1地址分配初始代币 ======================")
	testUserContractTokenEVMCreate(client, admin1, admin2, admin3, admin4, true, true)

	fmt.Println("====================== 查看余额 ======================")
	testUserContractTokenEVMBalanceOf(client, client1AddrInt, true)
	testUserContractTokenEVMBalanceOf(client, client2AddrInt, true)

	fmt.Println("====================== client1给client2地址转账 ======================")
	testUserContractTokenEVMTransfer(client, amount, true)

	fmt.Println("====================== 查看余额 ======================")
	testUserContractTokenEVMBalanceOf(client, client1AddrInt, true)
	testUserContractTokenEVMBalanceOf(client, client2AddrInt, true)
}

func testUserContractTokenEVMCreate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	withSyncResult bool, isIgnoreSameContract bool) {

	abiJson, err := ioutil.ReadFile(tokenABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	// 方式1: 16进制
	//addr := evmutils.BigToAddress(evmutils.FromHexString(client1Addr[2:]))
	// 方式2: Int
	//addr := evmutils.BigToAddress(evmutils.FromDecimalString(client1AddrInt))
	// 方式3: ski
	addrInt, err := evmutils.MakeAddressFromHex(client1AddrSki)
	if err != nil {
		log.Fatalln(err)
	}
	addr := evmutils.BigToAddress(addrInt)

	dataByte, err := myAbi.Pack("", addr)
	if err != nil {
		log.Fatalln(err)
	}

	data := hex.EncodeToString(dataByte)
	pairs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: data,
		},
	}

	byteCode, err := ioutil.ReadFile(tokenByteCodePath)
	if err != nil {
		log.Fatalln(err)
	}

	//bc := string(byteCode)
	//bc = strings.TrimSpace(bc)

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		//tokenContractName, tokenVersion, bc + data, common.RuntimeType_EVM, pairs, withSyncResult)
		tokenContractName, tokenVersion, string(byteCode), common.RuntimeType_EVM, pairs, withSyncResult)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE EVM token contract resp: %+v\n", resp)
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

func testUserContractTokenEVMTransfer(client *sdk.ChainClient, amount int64, withSyncResult bool) {

	abiJson, err := ioutil.ReadFile(tokenABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	//addr := evmutils.StringToAddress(client2Addr)
	addr := evmutils.BigToAddress(evmutils.FromDecimalString(client2AddrInt))

	dataByte, err := myAbi.Pack("transfer", addr, big.NewInt(amount))
	if err != nil {
		log.Fatalln(err)
	}

	data := hex.EncodeToString(dataByte)
	method := data[0:8]

	pairs := map[string]string{
		"data": data,
	}

	err = invokeUserContract(client, tokenContractName, method, "", pairs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func testUserContractTokenEVMBalanceOf(client *sdk.ChainClient, address string, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(tokenABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	//addr := evmutils.StringToAddress(address)
	addr := evmutils.BigToAddress(evmutils.FromDecimalString(address))

	methodName := "balanceOf"
	dataByte, err := myAbi.Pack(methodName, addr)
	if err != nil {
		log.Fatalln(err)
	}

	data := hex.EncodeToString(dataByte)
	method := data[0:8]

	pairs := map[string]string{
		"data": data,
	}

	result, err := invokeUserContractWithResult(client, tokenContractName, method, "", pairs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}

	balance, err := myAbi.Unpack(methodName, result)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("addr [%s] => %d\n", address, balance)
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
