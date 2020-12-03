/**
 * @Author: jasonruan
 * @Date:   2020-12-02 18:40:10
 **/
package chainmaker_sdk_go

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/stretchr/testify/require"
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

	// 2) [BlockUpdate]
	tx_timestamp_verify := rand.Intn(2) == 0
	txTimeout := rand.Intn(1000) + 600
	blockTxCapacity := rand.Intn(1000) + 1
	blockSize := rand.Intn(10) + 1
	blockInterval := rand.Intn(10000) + 10
	testChainConfigBlockUpdate(t, client, admin1, admin2, admin3, admin4, tx_timestamp_verify, txTimeout, blockTxCapacity, blockSize, blockInterval)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	fmt.Printf("tx_timestamp_verify: %s, txTimeout: %d, blockTxCapacity: %d, blockSize: %d, blockInterval: %d\n", strconv.FormatBool(tx_timestamp_verify), txTimeout, blockTxCapacity, blockSize, blockInterval)
	fmt.Printf("chainConfig txSchedulerTimeout: tx_timestamp_verify: %s, txTimeout: %d, blockTxCapacity: %d, blockSize: %d, blockInterval: %d\n",
		strconv.FormatBool(chainConfig.Block.TxTimestampVerify), chainConfig.Block.TxTimeout, chainConfig.Block.BlockTxCapacity, chainConfig.Block.BlockSize, chainConfig.Block.BlockInterval)
	require.Equal(t, tx_timestamp_verify, chainConfig.Block.TxTimestampVerify)
	require.Equal(t, txTimeout, int(chainConfig.Block.TxTimeout))
	require.Equal(t, blockTxCapacity, int(chainConfig.Block.BlockTxCapacity))
	require.Equal(t, blockSize, int(chainConfig.Block.BlockSize))
	require.Equal(t, blockInterval, int(chainConfig.Block.BlockInterval))

	// 3) [TrustRootAdd]
	raw, err := ioutil.ReadFile("testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt")
	require.Nil(t, err)
	trustRootOrgId := "wx-org5.chainmaker.org"
	trustRootCrt := string(raw)
	testChainConfigTrustRootAdd(t, client, admin1, admin2, admin3, admin4, trustRootOrgId, trustRootCrt)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, len(chainConfig.TrustRoots), 5)
	require.Equal(t, chainConfig.TrustRoots[4].OrgId, trustRootOrgId)
	require.Equal(t, chainConfig.TrustRoots[4].Root, trustRootCrt)

	// 4) [TrustRootUpdate]
	admin5, err := createAdmin(orgId5)
	require.Nil(t, err)
	raw, err = ioutil.ReadFile("testdata/crypto-config/wx-org6.chainmaker.org/ca/ca.crt")
	require.Nil(t, err)
	trustRootOrgId = "wx-org5.chainmaker.org"
	trustRootCrt = string(raw)
	testChainConfigTrustRootUpdate(t, client, admin1, admin2, admin3, admin5, trustRootOrgId, trustRootCrt)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, len(chainConfig.TrustRoots), 5)
	require.Equal(t, chainConfig.TrustRoots[4].OrgId, trustRootOrgId)
	require.Equal(t, chainConfig.TrustRoots[4].Root, trustRootCrt)

	// 5) [TrustRootDelete]
	trustRootOrgId = "wx-org5.chainmaker.org"
	trustRootCrt = string(raw)
	testChainConfigTrustRootDelete(t, client, admin1, admin2, admin3, admin5, trustRootOrgId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, len(chainConfig.TrustRoots), 4)
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

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigBlockUpdate(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	txTimestampVerify bool,
	txTimeout, blockTxCapacity, blockSize, blockInterval int) {

	// 配置块更新payload生成
	payloadBytes, err := client.ChainConfigCreateBlockUpdatePayload(
		txTimestampVerify, txTimeout, blockTxCapacity, blockSize, blockInterval)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigTrustRootAdd(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	trustRootOrgId, trustRootCrt string) {

	// 配置块更新payload生成
	payloadBytes, err := client.ChainConfigCreateTrustRootAddPayload(trustRootOrgId, trustRootCrt)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigTrustRootUpdate(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	trustRootOrgId, trustRootCrt string) {

	// 配置块更新payload生成
	payloadBytes, err := client.ChainConfigCreateTrustRootUpdatePayload(trustRootOrgId, trustRootCrt)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigTrustRootDelete(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	trustRootOrgId string) {

	// 配置块更新payload生成
	payloadBytes, err := client.ChainConfigCreateTrustRootDeletePayload(trustRootOrgId)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func signAndSendRequest(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	payloadBytes []byte) {
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
