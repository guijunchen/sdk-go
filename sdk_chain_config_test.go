/**
 * @Author: jasonruan
 * @Date:   2020-12-02 18:40:10
 **/
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestChainConfig(t *testing.T) {
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

	// 1) [CoreUpdate]
	rand.Seed(time.Now().UnixNano())
	txSchedulerTimeout := rand.Intn(61)
	txSchedulerValidateTimeout := rand.Intn(61)
	testChainConfigCoreUpdate(t, client, admin1, admin2, admin3, admin4, txSchedulerTimeout, txSchedulerValidateTimeout)
	time.Sleep(2 * time.Second)
	chainConfig := testGetChainConfig(t, client)
	fmt.Printf("txSchedulerTimeout: %d, txSchedulerValidateTimeout: %d\n", txSchedulerTimeout, txSchedulerValidateTimeout)
	fmt.Printf("chainConfig txSchedulerTimeout: %d, txSchedulerValidateTimeout: %d\n",
		chainConfig.Core.TxSchedulerTimeout, chainConfig.Core.TxSchedulerValidateTimeout)
	require.Equal(t, int(chainConfig.Core.TxSchedulerTimeout), txSchedulerTimeout)
	require.Equal(t, int(chainConfig.Core.TxSchedulerValidateTimeout), txSchedulerValidateTimeout)
}

func testGetChainConfig(t *testing.T, client *ChainClient) *pb.ChainConfig {
	resp, err := client.ChainConfigGet()
	require.Nil(t, err)
	fmt.Printf("GetChainConfig resp: %+v\n", resp)
	return resp
}

func testGetChainConfigByBlockHeight(t *testing.T, client *ChainClient) {
	resp, err := client.ChainConfigGetByBlockHeight(1)
	require.Nil(t, err)
	fmt.Printf("GetChainConfig resp: %+v\n", resp)
}

func testGetChainConfigSeq(t *testing.T, client *ChainClient) {
	seq, err := client.ChainConfigGetSeq()
	require.Nil(t, err)
	fmt.Printf("chainconfig seq: %d\n", seq)
}

func testChainConfigCoreUpdate(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	txSchedulerTimeout, txSchedulerValidateTimeout int) {

	// 配置块更新payload生成
	payloadBytes, err := client.ChainConfigCreateCoreUpdatePayload(
		txSchedulerTimeout, txSchedulerValidateTimeout)
	require.Nil(t, err)

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.ChainConfigPayloadCollectSign(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes2, err := admin2.ChainConfigPayloadCollectSign(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes3, err := admin3.ChainConfigPayloadCollectSign(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes4, err := admin4.ChainConfigPayloadCollectSign(payloadBytes)
	require.Nil(t, err)

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.ChainConfigPayloadMergeSign([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	require.Nil(t, err)

	// 发送配置更新请求
	resp, err := client.SendChainConfigUpdateRequest(mergeSignedPayloadBytes)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("chain config [CoreUpdate] resp: %+v", resp)
}
