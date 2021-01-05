/**
 * @Author: jasonruan
 * @Date:   2020-12-31 11:31:41
 **/
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
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

	fmt.Println("====================== 调用合约（异步）======================")
	testUserContractCounterGoInvoke(t, client, "increase", nil, false)
	time.Sleep(5 * time.Second)

	fmt.Println("====================== 执行合约查询接口 ======================")
	testUserContractCounterGoQuery(t, client, "query", nil)
}

// [用户合约]
func testMultiSignUserContractCreate(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient) {

	var (
		err error
		payloadBytes []byte
		entry *pb.EndorsementEntry
		resp *pb.TxResponse
	)

	payloadBytes, err = client.CreateContractCreatePayload(contractName, version, byteCodePath, pb.RuntimeType_WASMER, []*pb.KeyValuePair{})
	require.Nil(t, err)

	entry, err = admin1.SignMultiSignPayload(payloadBytes)
	require.Nil(t, err)
	resp, err = admin1.SendMultiSignReq(pb.TxType_MANAGE_USER_CONTRACT, payloadBytes, entry, 100000, -1)
	require.Nil(t, err)
	fmt.Printf("send multi sign req resp: code:%d, msg:%s, payload:%+v\n", resp.Code, resp.Message, resp.ContractResult)

	txId := string(resp.ContractResult.Result)

	// 休眠，等待多签请求完成
	time.Sleep(5 * time.Second)

	entry, err = admin2.SignMultiSignPayload(payloadBytes)
	require.Nil(t, err)
	resp, err = admin2.SendMultiSignVote(pb.VoteStatus_AGREE, txId, "", entry, -1)
	require.Nil(t, err)
	fmt.Printf("send multi sign vote1 resp: code:%d, msg:%s, payload:%+v\n", resp.Code, resp.Message, resp.ContractResult)

	entry, err = admin3.SignMultiSignPayload(payloadBytes)
	require.Nil(t, err)
	resp, err = admin3.SendMultiSignVote(pb.VoteStatus_AGREE, txId, "", entry, -1)
	require.Nil(t, err)
	fmt.Printf("send multi sign vote2 resp: code:%d, msg:%s, payload:%+v\n", resp.Code, resp.Message, resp.ContractResult)

	entry, err = admin4.SignMultiSignPayload(payloadBytes)
	require.Nil(t, err)
	resp, err = admin4.SendMultiSignVote(pb.VoteStatus_AGREE, txId, "", entry, -1)
	require.Nil(t, err)
	fmt.Printf("send multi sign vote3 resp: code:%d, msg:%s, payload:%+v\n", resp.Code, resp.Message, resp.ContractResult)
}
