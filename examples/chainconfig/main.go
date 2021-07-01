/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"

	"chainmaker.org/chainmaker/pb-go/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/pb-go/config"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
)

const (
	testKey     = "key001"
	nodePeerId1 = "QmQVkTSF6aWzRSddT3rro6Ve33jhKpsHFaQoVxHKMWzhuN"
	nodePeerId2 = "QmQVkTSF6aWzRSddT3rro6Ve33jhKpsHFaQoVxHKMWzhuN"

	sdkConfigOrg1Admin1Path  = "../sdk_configs/sdk_config_org1_admin1.yml"
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
	sdkConfigOrg2Admin1Path  = "../sdk_configs/sdk_config_org2_admin1.yml"
	sdkConfigOrg3Admin1Path  = "../sdk_configs/sdk_config_org3_admin1.yml"
	sdkConfigOrg4Admin1Path  = "../sdk_configs/sdk_config_org4_admin1.yml"
	sdkConfigOrg5Admin1Path  = "../sdk_configs/sdk_config_org5_admin1.yml"
)

func main() {
	testChainConfig()
}

func testChainConfig() {
	var (
		chainConfig *config.ChainConfig
	)

	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		panic(err)
	}

	admin1, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Admin1Path)
	if err != nil {
		panic(err)
	}
	admin2, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg2Admin1Path)
	if err != nil {
		panic(err)
	}
	admin3, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg3Admin1Path)
	if err != nil {
		panic(err)
	}
	admin4, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg4Admin1Path)
	if err != nil {
		panic(err)
	}

	// 1) [CoreUpdate]
	rand.Seed(time.Now().UnixNano())
	txSchedulerTimeout := rand.Intn(61)
	txSchedulerValidateTimeout := rand.Intn(61)
	testChainConfigCoreUpdate(client, admin1, admin2, admin3, admin4, txSchedulerTimeout, txSchedulerValidateTimeout)
	time.Sleep(5 * time.Second)
	chainConfig = testGetChainConfig(client)
	fmt.Printf("txSchedulerTimeout: %d, txSchedulerValidateTimeout: %d\n", txSchedulerTimeout, txSchedulerValidateTimeout)
	fmt.Printf("chainConfig txSchedulerTimeout: %d, txSchedulerValidateTimeout: %d\n",
		chainConfig.Core.TxSchedulerTimeout, chainConfig.Core.TxSchedulerValidateTimeout)
	if txSchedulerTimeout != int(chainConfig.Core.TxSchedulerTimeout) {
		panic("require txSchedulerTimeout == int(chainConfig.Core.TxSchedulerTimeout)")
	}
	if txSchedulerValidateTimeout != int(chainConfig.Core.TxSchedulerValidateTimeout) {
		panic("require txSchedulerValidateTimeout == int(chainConfig.Core.TxSchedulerValidateTimeout)")
	}

	// 2) [BlockUpdate]
	txTimestampVerify := rand.Intn(2) == 0
	txTimeout := rand.Intn(1000) + 600
	blockTxCapacity := rand.Intn(1000) + 1
	blockSize := rand.Intn(10) + 1
	blockInterval := rand.Intn(10000) + 10
	testChainConfigBlockUpdate(client, admin1, admin2, admin3, admin4, txTimestampVerify, txTimeout, blockTxCapacity, blockSize, blockInterval)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	fmt.Printf("tx_timestamp_verify: %s, txTimeout: %d, blockTxCapacity: %d, blockSize: %d, blockInterval: %d\n", strconv.FormatBool(txTimestampVerify), txTimeout, blockTxCapacity, blockSize, blockInterval)
	fmt.Printf("chainConfig txSchedulerTimeout: tx_timestamp_verify: %s, txTimeout: %d, blockTxCapacity: %d, blockSize: %d, blockInterval: %d\n",
		strconv.FormatBool(chainConfig.Block.TxTimestampVerify), chainConfig.Block.TxTimeout, chainConfig.Block.BlockTxCapacity, chainConfig.Block.BlockSize, chainConfig.Block.BlockInterval)
	if chainConfig.Block.TxTimestampVerify != txTimestampVerify {
		panic("require chainConfig.Block.TxTimestampVerify == txTimestampVerify")
	}
	if txTimeout != int(chainConfig.Block.TxTimeout) {
		panic("require txTimeout == int(chainConfig.Block.TxTimeout)")
	}
	if blockTxCapacity != int(chainConfig.Block.BlockTxCapacity) {
		panic("require equal")
	}
	if blockSize != int(chainConfig.Block.BlockSize) {
		panic("require equal")
	}
	if blockInterval != int(chainConfig.Block.BlockInterval) {
		panic("require equal")
	}

	// 3) [TrustRootAdd]
	trustCount := len(testGetChainConfig(client).TrustRoots)
	raw, err := ioutil.ReadFile("../../testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt")
	if err != nil {
		panic(err)
	}
	trustRootOrgId := examples.OrgId5
	trustRootCrt := string(raw)
	testChainConfigTrustRootAdd(client, admin1, admin2, admin3, admin4, trustRootOrgId, trustRootCrt)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if trustCount+1 != len(chainConfig.TrustRoots) {
		panic("require equal")
	}
	if trustRootOrgId != chainConfig.TrustRoots[trustCount].OrgId {
		panic("require equal")
	}
	if trustRootCrt != chainConfig.TrustRoots[trustCount].Root {
		panic("require equal")
	}

	// 4) [TrustRootUpdate]
	admin5, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg5Admin1Path)
	if err != nil {
		panic(err)
	}
	raw, err = ioutil.ReadFile("../../testdata/crypto-config/wx-org6.chainmaker.org/ca/ca.crt")
	if err != nil {
		panic(err)
	}
	trustRootOrgId = examples.OrgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootUpdate(client, admin1, admin2, admin3, admin5, trustRootOrgId, trustRootCrt)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if trustCount+1 != len(chainConfig.TrustRoots) {
		panic("require equal")
	}
	if trustRootOrgId != chainConfig.TrustRoots[trustCount].OrgId {
		panic("require equal")
	}
	if trustRootCrt != chainConfig.TrustRoots[trustCount].Root {
		panic("require equal")
	}

	// 5) [TrustRootDelete]
	trustRootOrgId = examples.OrgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootDelete(client, admin1, admin2, admin3, admin5, trustRootOrgId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if trustCount != len(chainConfig.TrustRoots) {
		panic("require equal")
	}

	// 6) [PermissionAdd]
	permissionCount := len(testGetChainConfig(client).ResourcePolicies)
	permissionResourceName := "TEST_PREMISSION"
	policy := &accesscontrol.Policy{
		Rule: "ANY",
	}
	testChainConfigPermissionAdd(client, admin1, admin2, admin3, admin4, permissionResourceName, policy)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if permissionCount+1 != len(chainConfig.ResourcePolicies) {
		panic("require equal")
	}
	if !proto.Equal(policy, chainConfig.ResourcePolicies[permissionCount].Policy) {
		panic("require true")
	}

	// 7) [PermissionUpdate]
	permissionResourceName = "TEST_PREMISSION"
	policy = &accesscontrol.Policy{
		Rule: "ANY",
	}
	testChainConfigPermissionUpdate(client, admin1, admin2, admin3, admin4, permissionResourceName, policy)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if permissionCount+1 != len(chainConfig.ResourcePolicies) {
		panic("require equal")
	}
	if !proto.Equal(policy, chainConfig.ResourcePolicies[permissionCount].Policy) {
		panic("require true")
	}

	// 8) [PermissionDelete]
	testChainConfigPermissionDelete(client, admin1, admin2, admin3, admin4, permissionResourceName)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if permissionCount != len(chainConfig.ResourcePolicies) {
		panic("require equal")
	}

	// 9) [ConsensusNodeAddrAdd]
	nodeOrgId := examples.OrgId4
	nodeIds := []string{nodePeerId1}
	testChainConfigConsensusNodeIdAdd(client, admin1, admin2, admin3, admin4, nodeOrgId, nodeIds)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if nodeOrgId != chainConfig.Consensus.Nodes[3].OrgId {
		panic("require equal")
	}
	if 2 != len(chainConfig.Consensus.Nodes[3].NodeId) {
		panic("require equal")
	}
	if nodeIds[0] != chainConfig.Consensus.Nodes[3].NodeId[1] {
		panic("require equal")
	}

	// 10) [ConsensusNodeAddrUpdate]
	nodeOrgId = examples.OrgId4
	nodeOldId := nodePeerId1
	nodeNewId := nodePeerId2
	testChainConfigConsensusNodeIdUpdate(client, admin1, admin2, admin3, admin4, nodeOrgId, nodeOldId, nodeNewId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if nodeOrgId != chainConfig.Consensus.Nodes[3].OrgId {
		panic("require equal")
	}
	if 2 != len(chainConfig.Consensus.Nodes[3].NodeId) {
		panic("require equal")
	}
	if nodeNewId != chainConfig.Consensus.Nodes[3].NodeId[1] {
		panic("require equal")
	}

	// 11) [ConsensusNodeAddrDelete]
	nodeOrgId = examples.OrgId4
	testChainConfigConsensusNodeIdDelete(client, admin1, admin2, admin3, admin4, nodeOrgId, nodeNewId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if nodeOrgId != chainConfig.Consensus.Nodes[3].OrgId {
		panic("require equal")
	}
	if 1 != len(chainConfig.Consensus.Nodes[3].NodeId) {
		panic("require equal")
	}

	// 12) [ConsensusNodeOrgAdd]
	raw, err = ioutil.ReadFile("../../testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt")
	if err != nil {
		panic(err)
	}
	trustRootOrgId = examples.OrgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootAdd(client, admin1, admin2, admin3, admin4, trustRootOrgId, trustRootCrt)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 5 != len(chainConfig.TrustRoots) {
		panic("require equal")
	}
	if trustRootOrgId != chainConfig.TrustRoots[4].OrgId {
		panic("require equal")
	}
	if trustRootCrt != chainConfig.TrustRoots[4].Root {
		panic("require equal")
	}
	nodeOrgId = examples.OrgId5
	nodeIds = []string{nodePeerId1}
	testChainConfigConsensusNodeOrgAdd(client, admin1, admin2, admin3, admin4, nodeOrgId, nodeIds)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 5 != len(chainConfig.Consensus.Nodes) {
		panic("require equal")
	}
	if nodeOrgId != chainConfig.Consensus.Nodes[4].OrgId {
		panic("require equal")
	}
	if 1 != len(chainConfig.Consensus.Nodes[4].NodeId) {
		panic("require equal")
	}
	if nodeIds[0] != chainConfig.Consensus.Nodes[4].NodeId[0] {
		panic("require equal")
	}

	// 13) [ConsensusNodeOrgUpdate]
	nodeOrgId = examples.OrgId5
	nodeIds = []string{nodePeerId2}
	testChainConfigConsensusNodeOrgUpdate(client, admin1, admin2, admin3, admin4, nodeOrgId, nodeIds)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 5 != len(chainConfig.Consensus.Nodes) {
		panic("require equal")
	}
	if nodeOrgId != chainConfig.Consensus.Nodes[4].OrgId {
		panic("require equal")
	}
	if 1 != len(chainConfig.Consensus.Nodes[4].NodeId) {
		panic("require equal")
	}
	if nodeIds[0] != chainConfig.Consensus.Nodes[4].NodeId[0] {
		panic("require equal")
	}

	// 14) [ConsensusNodeOrgDelete]
	nodeOrgId = examples.OrgId5
	testChainConfigConsensusNodeOrgDelete(client, admin1, admin2, admin3, admin4, nodeOrgId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 4 != len(chainConfig.Consensus.Nodes) {
		panic("require equal")
	}

	// 15) [ConsensusExtAdd]
	kvs := []*common.KeyValuePair{
		{
			Key:   testKey,
			Value: "test_value",
		},
	}
	testChainConfigConsensusExtAdd(client, admin1, admin2, admin3, admin4, kvs)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 2 != len(chainConfig.Consensus.ExtConfig) {
		panic("require equal")
	}
	if !proto.Equal(kvs[0], chainConfig.Consensus.ExtConfig[1]) {
		panic("require equal")
	}

	// 16) [ConsensusExtUpdate]
	kvs = []*common.KeyValuePair{
		{
			Key:   testKey,
			Value: "updated_value",
		},
	}
	testChainConfigConsensusExtUpdate(client, admin1, admin2, admin3, admin4, kvs)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 2 != len(chainConfig.Consensus.ExtConfig) {
		panic("require equal")
	}
	if !proto.Equal(kvs[0], chainConfig.Consensus.ExtConfig[1]) {
		panic("require equal")
	}

	// 16) [ConsensusExtDelete]
	keys := []string{testKey}
	testChainConfigConsensusExtDelete(client, admin1, admin2, admin3, admin4, keys)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 1 != len(chainConfig.Consensus.ExtConfig) {
		panic("require equal")
	}
}

func testGetChainConfig(client *sdk.ChainClient) *config.ChainConfig {
	resp, err := client.GetChainConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetChainConfig resp: %+v\n", resp)
	return resp
}

func testGetChainConfigByBlockHeight(client *sdk.ChainClient) {
	resp, err := client.GetChainConfigByBlockHeight(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("GetChainConfig resp: %+v\n", resp)
}

func testGetChainConfigSeq(client *sdk.ChainClient) {
	seq, err := client.GetChainConfigSequence()
	if err != nil {
		panic(err)
	}
	fmt.Printf("chainconfig seq: %d\n", seq)
}

func testChainConfigCoreUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, txSchedulerTimeout,
	txSchedulerValidateTimeout int) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigCoreUpdatePayload(
		txSchedulerTimeout, txSchedulerValidateTimeout)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigBlockUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, txTimestampVerify bool,
	txTimeout, blockTxCapacity, blockSize, blockInterval int) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigBlockUpdatePayload(
		txTimestampVerify, txTimeout, blockTxCapacity, blockSize, blockInterval)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigTrustRootAdd(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	trustRootOrgId, trustRootCrt string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigTrustRootAddPayload(trustRootOrgId, trustRootCrt)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigTrustRootUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	trustRootOrgId, trustRootCrt string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigTrustRootUpdatePayload(trustRootOrgId, trustRootCrt)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigTrustRootDelete(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, trustRootOrgId string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigTrustRootDeletePayload(trustRootOrgId)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigPermissionAdd(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	permissionResourceName string, policy *accesscontrol.Policy) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigPermissionAddPayload(permissionResourceName, policy)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigPermissionUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	permissionResourceName string, policy *accesscontrol.Policy) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigPermissionUpdatePayload(permissionResourceName, policy)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigPermissionDelete(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	permissionResourceName string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigPermissionDeletePayload(permissionResourceName)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeIdAdd(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId string, nodeIds []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeIdAddPayload(nodeAddrOrgId, nodeIds)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeIdUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId, nodeOldIds, nodeNewIds string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeIdUpdatePayload(nodeAddrOrgId, nodeOldIds, nodeNewIds)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeIdDelete(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId, nodeId string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeIdDeletePayload(nodeAddrOrgId, nodeId)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeOrgAdd(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId string, nodeIds []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeOrgAddPayload(nodeAddrOrgId, nodeIds)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeOrgUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId string, nodeIds []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeOrgUpdatePayload(nodeAddrOrgId, nodeIds)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeOrgDelete(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeOrgDeletePayload(nodeAddrOrgId)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusExtAdd(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	kvs []*common.KeyValuePair) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusExtAddPayload(kvs)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusExtUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	kvs []*common.KeyValuePair) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusExtUpdatePayload(kvs)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusExtDelete(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	keys []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusExtDeletePayload(keys)
	if err != nil {
		panic(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func signAndSendRequest(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, payloadBytes []byte) {
	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignChainConfigPayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes2, err := admin2.SignChainConfigPayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes3, err := admin3.SignChainConfigPayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes4, err := admin4.SignChainConfigPayload(payloadBytes)
	if err != nil {
		panic(err)
	}

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeChainConfigSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	if err != nil {
		panic(err)
	}

	// 发送配置更新请求
	resp, err := client.SendChainConfigUpdateRequest(mergeSignedPayloadBytes)
	if err != nil {
		panic(err)
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("chain config [CoreUpdate] resp: %+v", resp)
}
