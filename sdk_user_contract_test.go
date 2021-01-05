/**
 * @Author: jasonruan
 * @Date:   2020-12-02 18:41:47
 **/
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
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
	testUserContractCounterGoCreate(t, client, admin1, admin2, admin3, admin4)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 调用合约（异步）======================")
	testUserContractCounterGoInvoke(t, client, "increase", nil, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 执行合约查询接口 ======================")
	testUserContractCounterGoQuery(t, client, "query", nil)

	fmt.Println("====================== 冻结合约 ======================")
	testUserContractCounterGoFreeze(t, client, admin1, admin2, admin3, admin4, false)
	time.Sleep(5 * time.Second)
	fmt.Println("====================== 执行合约查询接口 ======================")
	testUserContractCounterGoQuery(t, client, "query", nil)

	fmt.Println("====================== 解冻合约 ======================")
	testUserContractCounterGoUnfreeze(t, client, admin1, admin2, admin3, admin4, false)
	time.Sleep(5 * time.Second)
	fmt.Println("====================== 执行合约查询接口 ======================")
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
	admin1, admin2, admin3, admin4 *ChainClient) {
	payloadBytes, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, pb.RuntimeType_WASMER, []*pb.KeyValuePair{})
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
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, false)
	fmt.Printf("resp: %+v\n", resp)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("CREATE counter-go contract resp: %+v\n", resp)
}

// 更新合约
func testUserContractCounterGoUpgrade(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient) {
	payloadBytes, err := client.CreateContractUpgradePayload(contractName, upgradeVersion, upgradeByteCodePath, pb.RuntimeType_GASM, []*pb.KeyValuePair{})
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
	fmt.Printf("resp: %+v\n", resp)
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
	fmt.Printf("resp: %+v\n", resp)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("Unfreeze counter-go contract resp: %+v\n", resp)
}

func testUserContractCounterGoInvoke(t *testing.T, client *ChainClient,
	method string, params map[string]string, withSyncResult bool) {
	resp, err := client.InvokeContract(contractName, method, "", params, -1, withSyncResult)
	require.Nil(t, err)

	if resp.Code != pb.TxStatusCode_SUCCESS {
		fmt.Printf("INVOKE counter-go contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if !withSyncResult {
		fmt.Printf("INVOKE counter-go contract resp, [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	} else {
		txInfo := new(pb.TransactionInfo)
		err := proto.Unmarshal(resp.ContractResult.Result, txInfo)
		require.Nil(t, err)

		bytes, err := json.Marshal(txInfo)
		fmt.Printf("INVOKE counter-go contract resp, [code:%d]/[msg:%s]/[txInfo:%s]\n", resp.Code, resp.Message, string(bytes))
	}
}

func testUserContractCounterGoQuery(t *testing.T, client *ChainClient,
	method string, params map[string]string) {
	resp, err := client.QueryContract(contractName, method, params, -1)
	require.Nil(t, err)
	fmt.Printf("QUERY counter-go contract resp: %+v\n", resp)
}
