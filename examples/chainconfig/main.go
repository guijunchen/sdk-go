/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
		log.Fatalln(err)
	}

	admin1, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Admin1Path)
	if err != nil {
		log.Fatalln(err)
	}
	admin2, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg2Admin1Path)
	if err != nil {
		log.Fatalln(err)
	}
	admin3, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg3Admin1Path)
	if err != nil {
		log.Fatalln(err)
	}
	admin4, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg4Admin1Path)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("====================== 根据区块高度获取链配置 ======================")
	testGetChainConfigByBlockHeight(client, 1)

	fmt.Println("====================== 获取链Sequence ======================")
	testGetChainConfigSeq(client)

	fmt.Println("====================== 更新CoreConfig ======================")
	rand.Seed(time.Now().UnixNano())
	txSchedulerTimeout := uint64(rand.Intn(61))
	txSchedulerValidateTimeout := uint64(rand.Intn(61))
	testChainConfigCoreUpdate(client, admin1, admin2, admin3, admin4, txSchedulerTimeout, txSchedulerValidateTimeout)
	time.Sleep(5 * time.Second)
	chainConfig = testGetChainConfig(client)
	fmt.Printf("txSchedulerTimeout: %d, txSchedulerValidateTimeout: %d\n", txSchedulerTimeout, txSchedulerValidateTimeout)
	fmt.Printf("chainConfig txSchedulerTimeout: %d, txSchedulerValidateTimeout: %d\n",
		chainConfig.Core.TxSchedulerTimeout, chainConfig.Core.TxSchedulerValidateTimeout)
	if txSchedulerTimeout != chainConfig.Core.TxSchedulerTimeout {
		log.Fatalln("require txSchedulerTimeout == int(chainConfig.Core.TxSchedulerTimeout)")
	}
	if txSchedulerValidateTimeout != chainConfig.Core.TxSchedulerValidateTimeout {
		log.Fatalln("require txSchedulerValidateTimeout == int(chainConfig.Core.TxSchedulerValidateTimeout)")
	}

	fmt.Println("====================== 更新BlockConfig ======================")
	txTimestampVerify := rand.Intn(2) == 0
	txTimeout := uint32(rand.Intn(1000)) + 600
	blockTxCapacity := uint32(rand.Intn(1000)) + 1
	blockSize := uint32(rand.Intn(10)) + 1
	blockInterval := uint32(rand.Intn(10000)) + 10
	testChainConfigBlockUpdate(client, admin1, admin2, admin3, admin4, txTimestampVerify, txTimeout, blockTxCapacity, blockSize, blockInterval)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	fmt.Printf("tx_timestamp_verify: %s, txTimeout: %d, blockTxCapacity: %d, blockSize: %d, blockInterval: %d\n", strconv.FormatBool(txTimestampVerify), txTimeout, blockTxCapacity, blockSize, blockInterval)
	fmt.Printf("chainConfig txSchedulerTimeout: tx_timestamp_verify: %s, txTimeout: %d, blockTxCapacity: %d, blockSize: %d, blockInterval: %d\n",
		strconv.FormatBool(chainConfig.Block.TxTimestampVerify), chainConfig.Block.TxTimeout, chainConfig.Block.BlockTxCapacity, chainConfig.Block.BlockSize, chainConfig.Block.BlockInterval)
	if chainConfig.Block.TxTimestampVerify != txTimestampVerify {
		log.Fatalln("require chainConfig.Block.TxTimestampVerify == txTimestampVerify")
	}
	if txTimeout != chainConfig.Block.TxTimeout {
		log.Fatalln("require txTimeout == int(chainConfig.Block.TxTimeout)")
	}
	if blockTxCapacity != chainConfig.Block.BlockTxCapacity {
		log.Fatalln("require equal")
	}
	if blockSize != chainConfig.Block.BlockSize {
		log.Fatalln("require equal")
	}
	if blockInterval != chainConfig.Block.BlockInterval {
		log.Fatalln("require equal")
	}

	fmt.Println("====================== 新增trust root ca ======================")
	trustCount := len(testGetChainConfig(client).TrustRoots)
	raw, err := ioutil.ReadFile("../../testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt")
	if err != nil {
		log.Fatalln(err)
	}
	trustRootOrgId := examples.OrgId5
	trustRootCrt := string(raw)
	testChainConfigTrustRootAdd(client, admin1, admin2, admin3, admin4, trustRootOrgId, trustRootCrt)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if trustCount+1 != len(chainConfig.TrustRoots) {
		log.Fatalln("require equal")
	}
	if trustRootOrgId != chainConfig.TrustRoots[trustCount].OrgId {
		log.Fatalln("require equal")
	}
	if trustRootCrt != chainConfig.TrustRoots[trustCount].Root {
		log.Fatalln("require equal")
	}

	fmt.Println("====================== 更新trust root ca ======================")
	admin5, err := examples.CreateChainClientWithSDKConfDisableCertHash(sdkConfigOrg5Admin1Path)
	if err != nil {
		log.Fatalln(err)
	}
	raw, err = ioutil.ReadFile("../../testdata/crypto-config/wx-org6.chainmaker.org/ca/ca.crt")
	if err != nil {
		log.Fatalln(err)
	}
	trustRootOrgId = examples.OrgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootUpdate(client, admin1, admin2, admin3, admin5, trustRootOrgId, trustRootCrt)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if trustCount+1 != len(chainConfig.TrustRoots) {
		log.Fatalln("require equal")
	}
	if trustRootOrgId != chainConfig.TrustRoots[trustCount].OrgId {
		log.Fatalln("require equal")
	}
	if trustRootCrt != chainConfig.TrustRoots[trustCount].Root {
		log.Fatalln("require equal")
	}

	fmt.Println("====================== 删除trust root ca ======================")
	trustRootOrgId = examples.OrgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootDelete(client, admin1, admin2, admin3, admin5, trustRootOrgId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if trustCount != len(chainConfig.TrustRoots) {
		log.Fatalln("require equal")
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
		log.Fatalln("require equal")
	}
	if !proto.Equal(policy, chainConfig.ResourcePolicies[permissionCount].Policy) {
		log.Fatalln("require true")
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
		log.Fatalln("require equal")
	}
	if !proto.Equal(policy, chainConfig.ResourcePolicies[permissionCount].Policy) {
		log.Fatalln("require true")
	}

	// 8) [PermissionDelete]
	testChainConfigPermissionDelete(client, admin1, admin2, admin3, admin4, permissionResourceName)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if permissionCount != len(chainConfig.ResourcePolicies) {
		log.Fatalln("require equal")
	}

	// 9) [ConsensusNodeAddrAdd]
	nodeOrgId := examples.OrgId4
	nodeIds := []string{nodePeerId1}
	testChainConfigConsensusNodeIdAdd(client, admin1, admin2, admin3, admin4, nodeOrgId, nodeIds)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if nodeOrgId != chainConfig.Consensus.Nodes[3].OrgId {
		log.Fatalln("require equal")
	}
	if 2 != len(chainConfig.Consensus.Nodes[3].NodeId) {
		log.Fatalln("require equal")
	}
	if nodeIds[0] != chainConfig.Consensus.Nodes[3].NodeId[1] {
		log.Fatalln("require equal")
	}

	// 10) [ConsensusNodeAddrUpdate]
	nodeOrgId = examples.OrgId4
	nodeOldId := nodePeerId1
	nodeNewId := nodePeerId2
	testChainConfigConsensusNodeIdUpdate(client, admin1, admin2, admin3, admin4, nodeOrgId, nodeOldId, nodeNewId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if nodeOrgId != chainConfig.Consensus.Nodes[3].OrgId {
		log.Fatalln("require equal")
	}
	if 2 != len(chainConfig.Consensus.Nodes[3].NodeId) {
		log.Fatalln("require equal")
	}
	if nodeNewId != chainConfig.Consensus.Nodes[3].NodeId[1] {
		log.Fatalln("require equal")
	}

	// 11) [ConsensusNodeAddrDelete]
	nodeOrgId = examples.OrgId4
	testChainConfigConsensusNodeIdDelete(client, admin1, admin2, admin3, admin4, nodeOrgId, nodeNewId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if nodeOrgId != chainConfig.Consensus.Nodes[3].OrgId {
		log.Fatalln("require equal")
	}
	if 1 != len(chainConfig.Consensus.Nodes[3].NodeId) {
		log.Fatalln("require equal")
	}

	// 12) [ConsensusNodeOrgAdd]
	raw, err = ioutil.ReadFile("../../testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt")
	if err != nil {
		log.Fatalln(err)
	}
	trustRootOrgId = examples.OrgId5
	trustRootCrt = string(raw)
	testChainConfigTrustRootAdd(client, admin1, admin2, admin3, admin4, trustRootOrgId, trustRootCrt)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 5 != len(chainConfig.TrustRoots) {
		log.Fatalln("require equal")
	}
	if trustRootOrgId != chainConfig.TrustRoots[4].OrgId {
		log.Fatalln("require equal")
	}
	if trustRootCrt != chainConfig.TrustRoots[4].Root {
		log.Fatalln("require equal")
	}
	nodeOrgId = examples.OrgId5
	nodeIds = []string{nodePeerId1}
	testChainConfigConsensusNodeOrgAdd(client, admin1, admin2, admin3, admin4, nodeOrgId, nodeIds)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 5 != len(chainConfig.Consensus.Nodes) {
		log.Fatalln("require equal")
	}
	if nodeOrgId != chainConfig.Consensus.Nodes[4].OrgId {
		log.Fatalln("require equal")
	}
	if 1 != len(chainConfig.Consensus.Nodes[4].NodeId) {
		log.Fatalln("require equal")
	}
	if nodeIds[0] != chainConfig.Consensus.Nodes[4].NodeId[0] {
		log.Fatalln("require equal")
	}

	// 13) [ConsensusNodeOrgUpdate]
	nodeOrgId = examples.OrgId5
	nodeIds = []string{nodePeerId2}
	testChainConfigConsensusNodeOrgUpdate(client, admin1, admin2, admin3, admin4, nodeOrgId, nodeIds)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 5 != len(chainConfig.Consensus.Nodes) {
		log.Fatalln("require equal")
	}
	if nodeOrgId != chainConfig.Consensus.Nodes[4].OrgId {
		log.Fatalln("require equal")
	}
	if 1 != len(chainConfig.Consensus.Nodes[4].NodeId) {
		log.Fatalln("require equal")
	}
	if nodeIds[0] != chainConfig.Consensus.Nodes[4].NodeId[0] {
		log.Fatalln("require equal")
	}

	// 14) [ConsensusNodeOrgDelete]
	nodeOrgId = examples.OrgId5
	testChainConfigConsensusNodeOrgDelete(client, admin1, admin2, admin3, admin4, nodeOrgId)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 4 != len(chainConfig.Consensus.Nodes) {
		log.Fatalln("require equal")
	}

	// 15) [ConsensusExtAdd]
	kvs := []*common.KeyValuePair{
		{
			Key:   testKey,
			Value: []byte("test_value"),
		},
	}
	testChainConfigConsensusExtAdd(client, admin1, admin2, admin3, admin4, kvs)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 2 != len(chainConfig.Consensus.ExtConfig) {
		log.Fatalln("require equal")
	}
	if !proto.Equal(kvs[0], chainConfig.Consensus.ExtConfig[1]) {
		log.Fatalln("require equal")
	}

	// 16) [ConsensusExtUpdate]
	kvs = []*common.KeyValuePair{
		{
			Key:   testKey,
			Value: []byte("updated_value"),
		},
	}
	testChainConfigConsensusExtUpdate(client, admin1, admin2, admin3, admin4, kvs)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 2 != len(chainConfig.Consensus.ExtConfig) {
		log.Fatalln("require equal")
	}
	if !proto.Equal(kvs[0], chainConfig.Consensus.ExtConfig[1]) {
		log.Fatalln("require equal")
	}

	// 16) [ConsensusExtDelete]
	keys := []string{testKey}
	testChainConfigConsensusExtDelete(client, admin1, admin2, admin3, admin4, keys)
	time.Sleep(2 * time.Second)
	chainConfig = testGetChainConfig(client)
	if 1 != len(chainConfig.Consensus.ExtConfig) {
		log.Fatalln("require equal")
	}
}

func testGetChainConfig(client *sdk.ChainClient) *config.ChainConfig {
	resp, err := client.GetChainConfig()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("GetChainConfig resp: %+v\n", resp)
	return resp
}

func testGetChainConfigByBlockHeight(client *sdk.ChainClient, blockHeight uint64) {
	resp, err := client.GetChainConfigByBlockHeight(blockHeight)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("GetChainConfig resp: %+v\n", resp)
}

func testGetChainConfigSeq(client *sdk.ChainClient) {
	seq, err := client.GetChainConfigSequence()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("chainconfig seq: %d\n", seq)
}

func testChainConfigCoreUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, txSchedulerTimeout,
	txSchedulerValidateTimeout uint64) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigCoreUpdatePayload(
		txSchedulerTimeout, txSchedulerValidateTimeout)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payload)
}

func testChainConfigBlockUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, txTimestampVerify bool,
	txTimeout, blockTxCapacity, blockSize, blockInterval uint32) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigBlockUpdatePayload(
		txTimestampVerify, txTimeout, blockTxCapacity, blockSize, blockInterval)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigTrustRootAdd(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	trustRootOrgId, trustRootCrt string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigTrustRootAddPayload(trustRootOrgId, trustRootCrt)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigTrustRootUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	trustRootOrgId, trustRootCrt string) {

	// 配置块更新payload生成
	payload, err := client.CreateChainConfigTrustRootUpdatePayload(trustRootOrgId, trustRootCrt)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payload)
}

func testChainConfigTrustRootDelete(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, trustRootOrgId string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigTrustRootDeletePayload(trustRootOrgId)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigPermissionAdd(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	permissionResourceName string, policy *accesscontrol.Policy) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigPermissionAddPayload(permissionResourceName, policy)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigPermissionUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	permissionResourceName string, policy *accesscontrol.Policy) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigPermissionUpdatePayload(permissionResourceName, policy)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigPermissionDelete(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	permissionResourceName string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigPermissionDeletePayload(permissionResourceName)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeIdAdd(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId string, nodeIds []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeIdAddPayload(nodeAddrOrgId, nodeIds)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeIdUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId, nodeOldIds, nodeNewIds string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeIdUpdatePayload(nodeAddrOrgId, nodeOldIds, nodeNewIds)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeIdDelete(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId, nodeId string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeIdDeletePayload(nodeAddrOrgId, nodeId)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeOrgAdd(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId string, nodeIds []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeOrgAddPayload(nodeAddrOrgId, nodeIds)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeOrgUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId string, nodeIds []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeOrgUpdatePayload(nodeAddrOrgId, nodeIds)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusNodeOrgDelete(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	nodeAddrOrgId string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusNodeOrgDeletePayload(nodeAddrOrgId)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusExtAdd(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	kvs []*common.KeyValuePair) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusExtAddPayload(kvs)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusExtUpdate(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	kvs []*common.KeyValuePair) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusExtUpdatePayload(kvs)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func testChainConfigConsensusExtDelete(client, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	keys []string) {

	// 配置块更新payload生成
	payloadBytes, err := client.CreateChainConfigConsensusExtDeletePayload(keys)
	if err != nil {
		log.Fatalln(err)
	}

	signAndSendRequest(client, admin1, admin2, admin3, admin4, payloadBytes)
}

func signAndSendRequest(client, admin1, admin2, admin3, admin4 *sdk.ChainClient, payload *common.Payload) {
	// 各组织Admin权限用户签名
	endorsementEntry1, err := admin1.SignChainConfigPayload(payload)
	if err != nil {
		log.Fatalln(err)
	}

	endorsementEntry2, err := admin2.SignChainConfigPayload(payload)
	if err != nil {
		log.Fatalln(err)
	}

	endorsementEntry3, err := admin3.SignChainConfigPayload(payload)
	if err != nil {
		log.Fatalln(err)
	}

	endorsementEntry4, err := admin4.SignChainConfigPayload(payload)
	if err != nil {
		log.Fatalln(err)
	}

	// 发送配置更新请求
	resp, err := client.SendChainConfigUpdateRequest(payload, []*common.EndorsementEntry{endorsementEntry1, endorsementEntry2, endorsementEntry3, endorsementEntry4}, -1, true)
	if err != nil {
		log.Fatalln(err)
	}

	err = examples.CheckProposalRequestResp(resp, false)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("ChainConfigUpdate resp: %+v\n", resp)
}
