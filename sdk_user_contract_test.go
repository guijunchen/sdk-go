/**
 * @Author: jasonruan
 * @Date:   2020-12-02 18:41:47
 **/
package chainmaker_sdk_go

import (
	"fmt"
	"testing"
	"time"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/stretchr/testify/require"
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

	testUserContractCounterGoCreate(t, client, admin1, admin2, admin3, admin4)
	time.Sleep(5 * time.Second)

	testUserContractCounterGoInvoke(t, client, "increase", nil)
	time.Sleep(5 * time.Second)

	testUserContractCounterGoQuery(t, client, "query", nil)

	testUserContractCounterGoUpgrade(t, client, admin1, admin2, admin3, admin4)
	time.Sleep(5 * time.Second)

	params := map[string]string{
		"key":   "jasonruan",
		"name":  "jasonruan",
		"value": "jasonruan",
	}
	testUserContractCounterGoInvoke(t, client, "upgrade_set_store", params)
	time.Sleep(5 * time.Second)

	testUserContractCounterGoQuery(t, client, "upgrade_get_store", params)
}

// [用户合约]
func testUserContractCounterGoCreate(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient) {
	payloadBytes, err := client.CreateContractManagePayload(TYPE_CREATE, contractName, version, byteCodePath, pb.RuntimeType_GASM_CPP, []*pb.KeyValuePair{})
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
	resp, err := client.SendContractManageRequest(TYPE_CREATE, mergeSignedPayloadBytes, -1)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("CREATE counter-go contract resp: %+v\n", resp)
}

func testUserContractCounterGoUpgrade(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient) {
	payloadBytes, err := client.CreateContractManagePayload(TYPE_UPGRADE, contractName, upgradeVersion, upgradeByteCodePath, pb.RuntimeType_GASM_CPP, []*pb.KeyValuePair{})
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
	resp, err := client.SendContractManageRequest(TYPE_UPGRADE, mergeSignedPayloadBytes, -1)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("CREATE counter-go upgrade resp: %+v\n", resp)
}

func testUserContractCounterGoInvoke(t *testing.T, client *ChainClient,
	method string, params map[string]string) {
	resp, err := client.InvokeContract(contractName, method, "", params, -1)
	require.Nil(t, err)
	fmt.Printf("INVOKE counter-go contract resp: %+v\n", resp)
}

func testUserContractCounterGoQuery(t *testing.T, client *ChainClient,
	method string, params map[string]string) {
	resp, err := client.QueryContract(contractName, method, params, -1)
	require.Nil(t, err)
	fmt.Printf("QUERY counter-go contract resp: %+v\n", resp)
}
