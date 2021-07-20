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
	balanceContractName   = "balance001"
	balanceVersion        = "1.0.0"
	balanceByteCodePath   = "../../testdata/balance-evm-demo/ledger_balance.bin"
	balanceABIPath        = "../../testdata/balance-evm-demo/ledger_balance.abi"

	// use cmc to calulate this addr, eg: `./cmc cert addr --cert-path xxx.tls.crt`
	client1AddrInt = "1018109374098032500766612781247089211099623418384"
	client2AddrInt = "1317892642413437150535769048733130623036570974971"
	amount         = 200

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"

	sdkConfigOrg1Admin1Path  = "../sdk_configs/sdk_config_org1_admin1.yml"
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

	//====================== create client ======================
	//====================== 创建Balance合约 ======================
	//	CREATE EVM balance contract resp: <nil>
	//====================== 设置addr2余额 ======================
	//	invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:gas_used:5888 ]
	//====================== 查看addr2余额 ======================
	//addr [1317892642413437150535769048733130623036570974971] => [1234]
	//====================== 设置自己余额 ======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:gas_used:5867 ]
	//====================== 查看自己余额 ======================
	//addr [1018109374098032500766612781247089211099623418384] => [1178]
	//====================== my给addr2地址转账 ======================
	//invoke contract success, resp: [code:0]/[msg:OK]/[contractResult:result:"\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\000\001" gas_used:18923 contract_event:<topic:"ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" tx_id:"23e83b60a9b1463e9dbde388aedcfb455fc0d28aa05640dda7e49c88ef563350" contract_name:"532c238cec7071ce8655aba07e50f9fb16f72ca1" contract_version:"1.0.0" event_data:"000000000000000000000000b2559a706dbab942671d5978ea97a13c82af6a10" event_data:"000000000000000000000000e6d859965e3dd4ef88c99130350b44b2ae9fbafb" event_data:"00000000000000000000000000000000000000000000000000000000000000c8" > ]
	//====================== 查看addr2余额 ======================
	//addr [1317892642413437150535769048733130623036570974971] => [1434]
	//====================== 查看自己余额 ======================
	//addr [1018109374098032500766612781247089211099623418384] => [978]
}

func calcContractName(contractName string) string {
	return hex.EncodeToString(evmutils.Keccak256([]byte(contractName)))[24:]
}

func testUserContractBalanceEVMCreate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	withSyncResult bool, isIgnoreSameContract bool) {

	byteCode, err := ioutil.ReadFile(balanceByteCodePath)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		calcContractName(balanceContractName), balanceVersion, string(byteCode), common.RuntimeType_EVM, nil, withSyncResult)
	if !isIgnoreSameContract {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("CREATE EVM balance contract resp: %+v\n", resp)
}

func createUserContract(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, contractName, version,
	byteCodePath string, runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {

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

	var endorsers []*common.EndorsementEntry
	endorsers = append(endorsers, signedPayload1)
	endorsers = append(endorsers, signedPayload2)
	endorsers = append(endorsers, signedPayload3)
	endorsers = append(endorsers, signedPayload4)

	// 发送创建合约请求
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

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	kvs := []*common.KeyValuePair {
		{
			Key: "data",
			Value: []byte(dataString),
		},
	}

	err = invokeUserContract(client, calcContractName(balanceContractName), method, "", kvs, withSyncResult)
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

	kvs := []*common.KeyValuePair {
		{
			Key: "data",
			Value: []byte(dataString),
		},
	}

	err = invokeUserContract(client, calcContractName(balanceContractName), method, "", kvs, withSyncResult)
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

	kvs := []*common.KeyValuePair {
		{
			Key: "data",
			Value: []byte(dataString),
		},
	}

	result, err := invokeUserContractWithResult(client, calcContractName(balanceContractName), method, "", kvs, withSyncResult)
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
	kvs []*common.KeyValuePair, withSyncResult bool) ([]byte, error) {

	resp, err := client.InvokeContract(contractName, method, txId, kvs, -1, withSyncResult)
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

	kvs := []*common.KeyValuePair {
		{
			Key: "data",
			Value: []byte(dataString),
		},
	}

	err = invokeUserContract(client, calcContractName(balanceContractName), method, "", kvs, withSyncResult)
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

	kvs := []*common.KeyValuePair {
		{
			Key: "data",
			Value: []byte(dataString),
		},
	}

	result, err := invokeUserContractWithResult(client, calcContractName(balanceContractName), method, "", kvs, withSyncResult)
	if err != nil {
		log.Fatalln(err)
	}

	balance, err := myAbi.Unpack("balances", result)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("addr [%s] => %d\n", address, balance)
}
