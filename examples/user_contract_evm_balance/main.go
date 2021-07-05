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
	balanceContractName   = "balance007"
	balanceVersion        = "1.0.0"
	balanceByteCodePath   = "../../testdata/balance-evm-demo/ledger_balance.bin"
	balanceABIPath        = "../../testdata/balance-evm-demo/ledger_balance.abi"

	client1AddrInt = "1087848554046178479107522336262214072175637027873"
	client2AddrInt = "944104665674401770091203869615921096651560803325"
	amount         = 200

	sdkConfigOrg1Admin1Path  = "../sdk_configs/sdk_config_org1_admin1.yml"
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
	sdkConfigOrg2Admin1Path  = "../sdk_configs/sdk_config_org2_admin1.yml"
	sdkConfigOrg3Admin1Path  = "../sdk_configs/sdk_config_org3_admin1.yml"
	sdkConfigOrg4Admin1Path  = "../sdk_configs/sdk_config_org4_admin1.yml"
)

func main() {
	testUserContractBalanceEVM()
}

func testUserContractBalanceEVM() {
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

	fmt.Println("====================== 创建Balance合约 ======================")
	testUserContractBalanceEVMCreate(client, admin1, admin2, admin3, admin4, true, true)

	fmt.Println("====================== 设置addr2余额 ======================")
	testUserContractBalanceEVMUpdateBalance(client, client2AddrInt, 1234, true)
	fmt.Println("====================== 查看addr2余额 ======================")
	testUserContractBalanceEVMGetBalance(client, client2AddrInt, true)

	fmt.Println("====================== 设置自己余额 ======================")
	testUserContractBalanceEVMUpdateMyBalance(client, 1178, true)
	fmt.Println("====================== 查看自己余额 ======================")
	testUserContractBalanceEVMGetMyBalance(client, client1AddrInt, true)

	fmt.Println("====================== my给addr2地址转账 ======================")
	testUserContractBalanceEVMTransfer(client, true)

	fmt.Println("====================== 查看addr2余额 ======================")
	testUserContractBalanceEVMGetBalance(client, client2AddrInt, true)
	fmt.Println("====================== 查看自己余额 ======================")
	testUserContractBalanceEVMGetMyBalance(client, client1AddrInt, true)
}

func testUserContractBalanceEVMCreate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	withSyncResult bool, isIgnoreSameContract bool) {

	byteCode, err := ioutil.ReadFile(balanceByteCodePath)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		balanceContractName, balanceVersion, string(byteCode), common.RuntimeType_EVM, nil, withSyncResult)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE EVM balance contract resp: %+v\n", resp)
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

func testUserContractBalanceEVMTransfer(client *sdk.ChainClient, withSyncResult bool) {

	abiJson, err := ioutil.ReadFile(balanceABIPath)
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

	err = invokeUserContract(client, balanceContractName, method, "", pairs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func testUserContractBalanceEVMUpdateBalance(client *sdk.ChainClient, address string, data int64, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(balanceABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	//addr := evmutils.StringToAddress(address)
	addr := evmutils.BigToAddress(evmutils.FromDecimalString(address))

	dataByte, err := myAbi.Pack("updateBalance", big.NewInt(data), addr)
	if err != nil {
		log.Fatalln(err)
	}

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	pairs := map[string]string{
		"data": dataString,
	}

	err = invokeUserContract(client, balanceContractName, method, "", pairs, withSyncResult)
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

func testUserContractBalanceEVMGetBalance(client *sdk.ChainClient, address string, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(balanceABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	addr := evmutils.BigToAddress(evmutils.FromDecimalString(address))

	dataByte, err := myAbi.Pack("balances", addr)
	if err != nil {
		log.Fatalln(err)
	}

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	pairs := map[string]string{
		"data": dataString,
	}

	result, err := invokeUserContractWithResult(client, balanceContractName, method, "", pairs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}

	balance, err := myAbi.Unpack("balances", result)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("addr [%s] => %d\n", address, balance)
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

func testUserContractBalanceEVMUpdateMyBalance(client *sdk.ChainClient, data int64, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(balanceABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	dataByte, err := myAbi.Pack("updateMyBalance", big.NewInt(data))
	if err != nil {
		log.Fatalln(err)
	}

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	pairs := map[string]string{
		"data": dataString,
	}

	err = invokeUserContract(client, balanceContractName, method, "", pairs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}
}

func testUserContractBalanceEVMGetMyBalance(client *sdk.ChainClient, address string, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(balanceABIPath)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}

	addr := evmutils.BigToAddress(evmutils.FromDecimalString(address))
	//dataByte, err := myAbi.Pack("balances", evmutils.BigToAddress(addr))
	dataByte, err := myAbi.Pack("balances", addr)
	if err != nil {
		log.Fatalln(err)
	}

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	pairs := map[string]string{
		"data": dataString,
	}

	result, err := invokeUserContractWithResult(client, balanceContractName, method, "", pairs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}

	balance, err := myAbi.Unpack("balances", result)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("addr [%s] => %d\n", address, balance)
}
