/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/pb/common"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	createContractTimeout = 5
)

func TestUserContractCounterGo(t *testing.T) {
	client, err := createClientWithConfig()
	require.Nil(t, err)

	admin1, err := createAdmin(orgId1)
	require.Nil(t, err)
	admin2, err := createAdmin(orgId2)
	require.Nil(t, err)
	admin3, err := createAdmin(orgId3)
	require.Nil(t, err)
	admin4, err := createAdmin(orgId4)
	require.Nil(t, err)

	fmt.Println("====================== 创建合约（异步）======================")
	testUserContractCounterGoCreate(t, client, admin1, admin2, admin3, admin4, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 调用合约（异步）======================")
	testUserContractCounterGoInvoke(t, client, "increase", nil, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 执行合约查询接口1 ======================")
	testUserContractCounterGoQuery(t, client, "query", nil)

	fmt.Println("====================== 冻结合约 ======================")
	testUserContractCounterGoFreeze(t, client, admin1, admin2, admin3, admin4, false)
	time.Sleep(5 * time.Second)
	fmt.Println("====================== 执行合约查询接口2 ======================")
	testUserContractCounterGoQuery(t, client, "query", nil)

	fmt.Println("====================== 解冻合约 ======================")
	testUserContractCounterGoUnfreeze(t, client, admin1, admin2, admin3, admin4, false)
	time.Sleep(5 * time.Second)
	fmt.Println("====================== 执行合约查询接口3 ======================")
	testUserContractCounterGoQuery(t, client, "query", nil)

	//fmt.Println("====================== 吊销合约 ======================")
	//testUserContractCounterGoRevoke(t, client, admin1, admin2, admin3, admin4, false)
	//time.Sleep(5 * time.Second)
	//fmt.Println("====================== 执行合约查询接口 ======================")
	//testUserContractCounterGoQuery(t, client, "query", nil)

	fmt.Println("====================== 调用合约（同步）======================")
	testUserContractCounterGoInvoke(t, client, "increase", nil, true)

	fmt.Println("====================== 执行合约查询接口 ======================")
	testUserContractCounterGoQuery(t, client, "query", nil)

	fmt.Println("====================== 升级合约（异步）======================")
	testUserContractCounterGoUpgrade(t, client, admin1, admin2, admin3, admin4)
	time.Sleep(5 * time.Second)

	params := map[string]string{
		"key":   "key001",
		"name":  "name001",
		"value": "value001",
	}
	testUserContractCounterGoInvoke(t, client, "upgrade_set_store", params, false)
	time.Sleep(5 * time.Second)

	testUserContractCounterGoQuery(t, client, "upgrade_get_store", params)
}

// [用户合约]
func testUserContractCounterGoCreate(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient, withSyncResult bool) {

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		contractName, version, byteCodePath, common.RuntimeType_WASMER, []*common.KeyValuePair{}, withSyncResult)
	require.Nil(t, err)

	fmt.Printf("CREATE counter-go contract resp: %+v\n", resp)
}

// 更新合约
func testUserContractCounterGoUpgrade(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient) {
	payloadBytes, err := client.CreateContractUpgradePayload(contractName, upgradeVersion, upgradeByteCodePath, common.RuntimeType_GASM, []*common.KeyValuePair{})
	require.Nil(t, err)

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	require.Nil(t, err)

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, -1, false)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("UPGRADE counter-go contract resp: %+v\n", resp)
}

// 冻结合约
func testUserContractCounterGoFreeze(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient, withSyncResult bool) {
	payloadBytes, err := client.CreateContractFreezePayload(contractName)
	require.Nil(t, err)

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	require.Nil(t, err)

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, withSyncResult)
	fmt.Printf("resp: %+v\n", resp)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("Freeze counter-go contract resp: %+v\n", resp)
}

// 解冻合约
func testUserContractCounterGoUnfreeze(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient, withSyncResult bool) {
	payloadBytes, err := client.CreateContractUnfreezePayload(contractName)
	require.Nil(t, err)

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	require.Nil(t, err)

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, withSyncResult)
	fmt.Printf("unfreeze resp: %+v\n", resp)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("Unfreeze counter-go contract resp: %+v\n", resp)
}

// 吊销合约
func testUserContractCounterGoRevoke(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient, withSyncResult bool) {
	payloadBytes, err := client.CreateContractRevokePayload(contractName)
	require.Nil(t, err)

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	require.Nil(t, err)

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	require.Nil(t, err)

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, withSyncResult)
	fmt.Printf("revoke resp: %+v\n", resp)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("revoke counter-go contract resp: %+v\n", resp)
}

func testUserContractCounterGoInvoke(t *testing.T, client *ChainClient,
	method string, params map[string]string, withSyncResult bool) {

	err := invokeUserContract(client, contractName, method, "", params, withSyncResult)
	require.Nil(t, err)
}

func testUserContractCounterGoQuery(t *testing.T, client *ChainClient,
	method string, params map[string]string) {
	resp, err := client.QueryContract(contractName, method, params, -1)
	require.Nil(t, err)
	fmt.Printf("QUERY counter-go contract resp: %+v\n", resp)
}

func createUserContract(client *ChainClient, admin1, admin2, admin3, admin4 *ChainClient,
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

	err = checkProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func invokeUserContract(client *ChainClient, contractName, method, txId string, params map[string]string, withSyncResult bool) error {

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
