/**
 * @Author: jasonruan
 * @Date:   2020-12-02 18:41:47
 **/
package chainmaker_sdk_go

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"testing"
	"time"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/stretchr/testify/require"
)

const (
	createContractTimeout = 7
)

func TestUserContractCounterGo(t *testing.T) {
	client, err := createClient()
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

	fmt.Println("====================== 调用合约（同步）======================")
	testUserContractCounterGoInvoke(t, client, "increase", nil, true)
	time.Sleep(5 * time.Second)

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
	payloadBytes, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, pb.RuntimeType_GASM, []*pb.KeyValuePair{})
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
	resp, err := client.SendContractCreateRequest(mergeSignedPayloadBytes, createContractTimeout, false)
	fmt.Printf("resp: %+v\n", resp)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("CREATE counter-go contract resp: %+v\n", resp)
}

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
	resp, err := client.SendContractUpgradeRequest(mergeSignedPayloadBytes, -1, false)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("UPGRADE counter-go contract resp: %+v\n", resp)
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
