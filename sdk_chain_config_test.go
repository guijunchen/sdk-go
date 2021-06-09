/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/accesscontrol"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/config"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

const (
	testKey     = "key001"
	nodePeerId1 = "QmQVkTSF6aWzRSddT3rro6Ve33jhKpsHFaQoVxHKMWzhuN"
	nodePeerId2 = "QmQVkTSF6aWzRSddT3rro6Ve33jhKpsHFaQoVxHKMWzhuN"
)

func TestChainConfig(t *testing.T) {
	var (
		chainConfig *config.ChainConfig
	)

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
	time.Sleep(5 * time.Second)
	chainConfig = testGetChainConfig(t, client)
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
	trustCount := len(testGetChainConfig(t, client).TrustRoots)
	raw, err := ioutil.ReadFile("testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt")
	require.Nil(t, err)
	trustRootOrgId := orgId5
	trustRootCrt := string(raw)
	testChainConfigTrustRootAdd(t, client, admin1, admin2, admin3, admin4, trustRootOrgId, trustRootCrt)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, trustCount+1, len(chainConfig.TrustRoots))
	require.Equal(t, trustRootOrgId, chainConfig.TrustRoots[trustCount].OrgId)
	require.Equal(t, trustRootCrt, chainConfig.TrustRoots[trustCount].Root)

	// 4) [TrustRootUpdate]
	admin5, err := createAdmin(orgId5)
	require.Nil(t, err)
	raw, err = ioutil.ReadFile("testdata/crypto-config/wx-org6.chainmaker.org/ca/ca.crt")
	require.Nil(t, err)
	trustRootOrgId = orgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootUpdate(t, client, admin1, admin2, admin3, admin5, trustRootOrgId, trustRootCrt)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, trustCount+1, len(chainConfig.TrustRoots))
	require.Equal(t, trustRootOrgId, chainConfig.TrustRoots[trustCount].OrgId)
	require.Equal(t, trustRootCrt, chainConfig.TrustRoots[trustCount].Root)

	// 5) [TrustRootDelete]
	trustRootOrgId = orgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootDelete(t, client, admin1, admin2, admin3, admin5, trustRootOrgId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, trustCount, len(chainConfig.TrustRoots))

	// 6) [PermissionAdd]
	permissionCount := len(testGetChainConfig(t, client).ResourcePolicies)
	permissionResourceName := "TEST_PREMISSION"
	policy := &accesscontrol.Policy{
		Rule: "ANY",
	}
	testChainConfigPermissionAdd(t, client, admin1, admin2, admin3, admin4, permissionResourceName, policy)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, permissionCount+1, len(chainConfig.ResourcePolicies))
	require.Equal(t, true, proto.Equal(policy, chainConfig.ResourcePolicies[permissionCount].Policy))

	// 7) [PermissionUpdate]
	permissionResourceName = "TEST_PREMISSION"
	policy = &accesscontrol.Policy{
		Rule: "ANY",
	}
	testChainConfigPermissionUpdate(t, client, admin1, admin2, admin3, admin4, permissionResourceName, policy)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, permissionCount+1, len(chainConfig.ResourcePolicies))
	require.Equal(t, true, proto.Equal(policy, chainConfig.ResourcePolicies[permissionCount].Policy))

	// 8) [PermissionDelete]
	testChainConfigPermissionDelete(t, client, admin1, admin2, admin3, admin4, permissionResourceName)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, permissionCount, len(chainConfig.ResourcePolicies))

	// 9) [ConsensusNodeAddrAdd]
	nodeOrgId := orgId4
	nodeIds := []string{nodePeerId1}
	testChainConfigConsensusNodeIdAdd(t, client, admin1, admin2, admin3, admin4, nodeOrgId, nodeIds)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, nodeOrgId, chainConfig.Consensus.Nodes[3].OrgId)
	require.Equal(t, 2, len(chainConfig.Consensus.Nodes[3].NodeId))
	require.Equal(t, nodeIds[0], chainConfig.Consensus.Nodes[3].NodeId[1])

	// 10) [ConsensusNodeAddrUpdate]
	nodeOrgId = orgId4
	nodeOldId := nodePeerId1
	nodeNewId := nodePeerId2
	testChainConfigConsensusNodeIdUpdate(t, client, admin1, admin2, admin3, admin4, nodeOrgId, nodeOldId, nodeNewId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, nodeOrgId, chainConfig.Consensus.Nodes[3].OrgId)
	require.Equal(t, 2, len(chainConfig.Consensus.Nodes[3].NodeId))
	require.Equal(t, nodeNewId, chainConfig.Consensus.Nodes[3].NodeId[1])

	// 11) [ConsensusNodeAddrDelete]
	nodeOrgId = orgId4
	testChainConfigConsensusNodeIdDelete(t, client, admin1, admin2, admin3, admin4, nodeOrgId, nodeNewId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, nodeOrgId, chainConfig.Consensus.Nodes[3].OrgId)
	require.Equal(t, 1, len(chainConfig.Consensus.Nodes[3].NodeId))

	// 12) [ConsensusNodeOrgAdd]
	raw, err = ioutil.ReadFile("testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt")
	require.Nil(t, err)
	trustRootOrgId = orgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootAdd(t, client, admin1, admin2, admin3, admin4, trustRootOrgId, trustRootCrt)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, 5, len(chainConfig.TrustRoots))
	require.Equal(t, trustRootOrgId, chainConfig.TrustRoots[4].OrgId)
	require.Equal(t, trustRootCrt, chainConfig.TrustRoots[4].Root)
	nodeOrgId = orgId5
	nodeIds = []string{nodePeerId1}
	testChainConfigConsensusNodeOrgAdd(t, client, admin1, admin2, admin3, admin4, nodeOrgId, nodeIds)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, 5, len(chainConfig.Consensus.Nodes))
	require.Equal(t, nodeOrgId, chainConfig.Consensus.Nodes[4].OrgId)
	require.Equal(t, 1, len(chainConfig.Consensus.Nodes[4].NodeId))
	require.Equal(t, nodeIds[0], chainConfig.Consensus.Nodes[4].NodeId[0])

	// 13) [ConsensusNodeOrgUpdate]
	nodeOrgId = orgId5
	nodeIds = []string{nodePeerId2}
	testChainConfigConsensusNodeOrgUpdate(t, client, admin1, admin2, admin3, admin4, nodeOrgId, nodeIds)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, 5, len(chainConfig.Consensus.Nodes))
	require.Equal(t, nodeOrgId, chainConfig.Consensus.Nodes[4].OrgId)
	require.Equal(t, 1, len(chainConfig.Consensus.Nodes[4].NodeId))
	require.Equal(t, nodeIds[0], chainConfig.Consensus.Nodes[4].NodeId[0])

	// 14) [ConsensusNodeOrgDelete]
	nodeOrgId = orgId5
	testChainConfigConsensusNodeOrgDelete(t, client, admin1, admin2, admin3, admin4, nodeOrgId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, 4, len(chainConfig.Consensus.Nodes))

	// 15) [ConsensusExtAdd]
	kvs := []*common.KeyValuePair{
		{
			Key:   testKey,
			Value: "test_value",
		},
	}
	testChainConfigConsensusExtAdd(t, client, admin1, admin2, admin3, admin4, kvs)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, 2, len(chainConfig.Consensus.ExtConfig))
	require.Equal(t, true, proto.Equal(kvs[0], chainConfig.Consensus.ExtConfig[1]))

	// 16) [ConsensusExtUpdate]
	kvs = []*common.KeyValuePair{
		{
			Key:   testKey,
			Value: "updated_value",
		},
	}
	testChainConfigConsensusExtUpdate(t, client, admin1, admin2, admin3, admin4, kvs)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, 2, len(chainConfig.Consensus.ExtConfig))
	require.Equal(t, true, proto.Equal(kvs[0], chainConfig.Consensus.ExtConfig[1]))

	// 16) [ConsensusExtDelete]
	keys := []string{testKey}
	testChainConfigConsensusExtDelete(t, client, admin1, admin2, admin3, admin4, keys)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(t, client)
	require.Equal(t, 1, len(chainConfig.Consensus.ExtConfig))
}

func testGetChainConfig(t *testing.T, client *ChainClient) *config.ChainConfig {
	resp, err := client.GetChainConfig()
	require.Nil(t, err)
	fmt.Printf("GetChainConfig resp: %+v\n", resp)
	return resp
}

func testGetChainConfigByBlockHeight(t *testing.T, client *ChainClient) {
	resp, err := client.GetChainConfigByBlockHeight(1)
	require.Nil(t, err)
	fmt.Printf("GetChainConfig resp: %+v\n", resp)
}

func testGetChainConfigSeq(t *testing.T, client *ChainClient) {
	seq, err := client.GetChainConfigSequence()
	require.Nil(t, err)
	fmt.Printf("chainconfig seq: %d\n", seq)
}

func testChainConfigCoreUpdate(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	txSchedulerTimeout, txSchedulerValidateTimeout int) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigCoreUpdatePayload(
		txSchedulerTimeout, txSchedulerValidateTimeout)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigBlockUpdate(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	txTimestampVerify bool,
	txTimeout, blockTxCapacity, blockSize, blockInterval int) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigBlockUpdatePayload(
		txTimestampVerify, txTimeout, blockTxCapacity, blockSize, blockInterval)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigTrustRootAdd(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	trustRootOrgId, trustRootCrt string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigTrustRootAddPayload(trustRootOrgId, trustRootCrt)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigTrustRootUpdate(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	trustRootOrgId, trustRootCrt string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigTrustRootUpdatePayload(trustRootOrgId, trustRootCrt)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigTrustRootDelete(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	trustRootOrgId string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigTrustRootDeletePayload(trustRootOrgId)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigPermissionAdd(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	permissionResourceName string, policy *accesscontrol.Policy) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigPermissionAddPayload(permissionResourceName, policy)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigPermissionUpdate(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	permissionResourceName string, policy *accesscontrol.Policy) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigPermissionUpdatePayload(permissionResourceName, policy)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigPermissionDelete(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	permissionResourceName string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigPermissionDeletePayload(permissionResourceName)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeIdAdd(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	nodeAddrOrgId string, nodeIds []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeIdAddPayload(nodeAddrOrgId, nodeIds)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeIdUpdate(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	nodeAddrOrgId, nodeOldIds, nodeNewIds string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeIdUpdatePayload(nodeAddrOrgId, nodeOldIds, nodeNewIds)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeIdDelete(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	nodeAddrOrgId, nodeId string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeIdDeletePayload(nodeAddrOrgId, nodeId)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeOrgAdd(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	nodeAddrOrgId string, nodeIds []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeOrgAddPayload(nodeAddrOrgId, nodeIds)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeOrgUpdate(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	nodeAddrOrgId string, nodeIds []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeOrgUpdatePayload(nodeAddrOrgId, nodeIds)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeOrgDelete(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	nodeAddrOrgId string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeOrgDeletePayload(nodeAddrOrgId)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusExtAdd(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	kvs []*common.KeyValuePair) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusExtAddPayload(kvs)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusExtUpdate(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	kvs []*common.KeyValuePair) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusExtUpdatePayload(kvs)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusExtDelete(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	keys []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusExtDeletePayload(keys)
	require.Nil(t, err)

	signAndSendRequest(t, client, admin1, admin2, admin3, admin4, payloadBytes)
}

func signAndSendRequest(t *testing.T, client,
	admin1, admin2, admin3, admin4 *ChainClient,
	payloadBytes []byte) {
	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignChainConfigPayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes2, err := admin2.SignChainConfigPayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes3, err := admin3.SignChainConfigPayload(payloadBytes)
	require.Nil(t, err)

	signedPayloadBytes4, err := admin4.SignChainConfigPayload(payloadBytes)
	require.Nil(t, err)

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeChainConfigSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	require.Nil(t, err)

	// 发送配置更新请求
	resp, err := client.SendChainConfigUpdateRequest(mergeSignedPayloadBytes)
	require.Nil(t, err)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)

	fmt.Printf("chain config [CoreUpdate] resp: %+v", resp)
}
