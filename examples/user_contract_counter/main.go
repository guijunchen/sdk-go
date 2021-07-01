/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"time"

	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
)

const (
	createContractTimeout = 5
	contractName          = "counter-go-1"
	version               = "1.0.0"
	upgradeVersion        = "2.0.0"
	byteCodePath          = "../../testdata/counter-go-demo/counter-rust-0.7.2.wasm"
	upgradeByteCodePath   = "../../testdata/counter-go-demo/counter-go-upgrade.wasm"
)

func main() {
	testUserContractCounterGo()
}

func testUserContractCounterGo() {
	client, err := examples.CreateClientWithCertBytes()
	if err != nil {
		panic(err)
	}

	admin1, err := examples.CreateAdmin(examples.OrgId1)
	if err != nil {
		panic(err)
	}
	admin2, err := examples.CreateAdmin(examples.OrgId2)
	if err != nil {
		panic(err)
	}
	admin3, err := examples.CreateAdmin(examples.OrgId3)
	if err != nil {
		panic(err)
	}
	admin4, err := examples.CreateAdmin(examples.OrgId4)
	if err != nil {
		panic(err)
	}

	fmt.Println("====================== 创建合约（异步）======================")
	testUserContractCounterGoCreate(client, admin1, admin2, admin3, admin4, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 调用合约（异步）======================")
	testUserContractCounterGoInvoke(client, "increase", nil, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 执行合约查询接口1 ======================")
	testUserContractCounterGoQuery(client, "query", nil)

	fmt.Println("====================== 冻结合约 ======================")
	testUserContractCounterGoFreeze(client, admin1, admin2, admin3, admin4, false)
	time.Sleep(5 * time.Second)
	fmt.Println("====================== 执行合约查询接口2 ======================")
	testUserContractCounterGoQuery(client, "query", nil)

	fmt.Println("====================== 解冻合约 ======================")
	testUserContractCounterGoUnfreeze(client, admin1, admin2, admin3, admin4, false)
	time.Sleep(5 * time.Second)
	fmt.Println("====================== 执行合约查询接口3 ======================")
	testUserContractCounterGoQuery(client, "query", nil)

	//fmt.Println("====================== 吊销合约 ======================")
	//testUserContractCounterGoRevoke(t, client, admin1, admin2, admin3, admin4, false)
	//time.Sleep(5 * time.Second)
	//fmt.Println("====================== 执行合约查询接口 ======================")
	//testUserContractCounterGoQuery(t, client, "query", nil)

	fmt.Println("====================== 调用合约（同步）======================")
	testUserContractCounterGoInvoke(client, "increase", nil, true)

	fmt.Println("====================== 执行合约查询接口 ======================")
	testUserContractCounterGoQuery(client, "query", nil)

	fmt.Println("====================== 升级合约（异步）======================")
	testUserContractCounterGoUpgrade(client, admin1, admin2, admin3, admin4)
	time.Sleep(5 * time.Second)

	params := map[string]string{
		"key":   "key001",
		"name":  "name001",
		"value": "value001",
	}
	testUserContractCounterGoInvoke(client, "upgrade_set_store", params, false)
	time.Sleep(5 * time.Second)

	testUserContractCounterGoQuery(client, "upgrade_get_store", params)
}

// [用户合约]
func testUserContractCounterGoCreate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, withSyncResult bool) {

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		contractName, version, byteCodePath, common.RuntimeType_WASMER, []*common.KeyValuePair{}, withSyncResult)
	if err != nil {
		panic(err)
	}

	fmt.Printf("CREATE counter-go contract resp: %+v\n", resp)
}

// 更新合约
func testUserContractCounterGoUpgrade(client, admin1, admin2, admin3, admin4 *sdk.ChainClient) {
	payloadBytes, err := client.CreateContractUpgradePayload(contractName, upgradeVersion, upgradeByteCodePath, common.RuntimeType_GASM, []*common.KeyValuePair{})
	if err != nil {
		panic(err)
	}

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	if err != nil {
		panic(err)
	}

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, -1, false)
	if err != nil {
		panic(err)
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("UPGRADE counter-go contract resp: %+v\n", resp)
}

// 冻结合约
func testUserContractCounterGoFreeze(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	withSyncResult bool) {
	payloadBytes, err := client.CreateContractFreezePayload(contractName)
	if err != nil {
		panic(err)
	}

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	if err != nil {
		panic(err)
	}

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, withSyncResult)
	fmt.Printf("resp: %+v\n", resp)
	if err != nil {
		panic(err)
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Freeze counter-go contract resp: %+v\n", resp)
}

// 解冻合约
func testUserContractCounterGoUnfreeze(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	withSyncResult bool) {
	payloadBytes, err := client.CreateContractUnfreezePayload(contractName)
	if err != nil {
		panic(err)
	}

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	if err != nil {
		panic(err)
	}

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, withSyncResult)
	fmt.Printf("unfreeze resp: %+v\n", resp)
	if err != nil {
		panic(err)
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Unfreeze counter-go contract resp: %+v\n", resp)
}

// 吊销合约
func testUserContractCounterGoRevoke(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	withSyncResult bool) {
	payloadBytes, err := client.CreateContractRevokePayload(contractName)
	if err != nil {
		panic(err)
	}

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	if err != nil {
		panic(err)
	}

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, withSyncResult)
	fmt.Printf("revoke resp: %+v\n", resp)
	if err != nil {
		panic(err)
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("revoke counter-go contract resp: %+v\n", resp)
}

func testUserContractCounterGoInvoke(client *sdk.ChainClient, method string, params map[string]string,
	withSyncResult bool) {

	err := invokeUserContract(client, contractName, method, "", params, withSyncResult)
	if err != nil {
		panic(err)
	}
}

func testUserContractCounterGoQuery(client *sdk.ChainClient, method string, params map[string]string) {
	resp, err := client.QueryContract(contractName, method, params, -1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("QUERY counter-go contract resp: %+v\n", resp)
}

func createUserContract(client *sdk.ChainClient, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	contractName, version, byteCodePath string, runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {

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

func invokeUserContractWithContractResult(client *sdk.ChainClient, contractName, method, txId string,
	params map[string]string, withSyncResult bool) (*common.ContractResult, error) {

	resp, err := client.InvokeContract(contractName, method, txId, params, -1, withSyncResult)
	if err != nil {
		return nil, err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return nil, fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	return resp.ContractResult, nil
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

func invokeUserContractStepByStep(client *sdk.ChainClient, contractName, method, txId string,
	params map[string]string, withSyncResult bool) error {
	req, err := client.GetTxRequest(contractName, method, "", params)
	if err != nil {
		return err
	}

	resp, err := client.SendTxRequest(req, -1, withSyncResult)
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
