/**
 * @Author: jasonruan
 * @Date:   2020-12-31 11:31:41
 **/
package chainmaker_sdk_go

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMultiSignUserContract(t *testing.T) {
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
	testMultiSignUserContractCreate(t, client, admin1, admin2, admin3, admin4)
	time.Sleep(5 * time.Second)
}

// [用户合约]
func testMultiSignUserContractCreate(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient) {
	//payloadBytes, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, pb.RuntimeType_WASMER, []*pb.KeyValuePair{})
	//require.Nil(t, err)
	//
	//entry1, err := admin1.SignMultiSignPayload(payloadBytes)
	//require.Nil(t, err)

	//resp, err := admin1.SendMultiSignReq(pb.TxType_CREATE_USER_CONTRACT, payloadBytes, entry1, 100000, -1)
	//require.Nil(t, err)
	//fmt.Printf("send multi sign req resp: code:%d, msg:%s, payload:%+v\n", resp.Code, resp.Message, resp.ContractResult)


	//// 各组织Admin权限用户签名
	//signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	//require.Nil(t, err)
	//
	//signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	//require.Nil(t, err)
	//
	//signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	//require.Nil(t, err)
	//
	//signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	//require.Nil(t, err)
	//
	//// 收集并合并签名
	//mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
	//	signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	//require.Nil(t, err)
	//
	//// 发送创建合约请求
	//resp, err := client.SendContractCreateRequest(mergeSignedPayloadBytes, createContractTimeout, false)
	//fmt.Printf("resp: %+v\n", resp)
	//require.Nil(t, err)
	//
	//err = checkProposalRequestResp(resp, true)
	//require.Nil(t, err)
	//
	//fmt.Printf("CREATE counter-go contract resp: %+v\n", resp)
}
