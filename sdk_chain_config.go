/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"

	"chainmaker.org/chainmaker/pb-go/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/pb-go/config"
	"chainmaker.org/chainmaker/pb-go/syscontract"
	"chainmaker.org/chainmaker/sdk-go/utils"
)

const (
	getCCSeqErrStringFormat         = "get chain config sequence failed, %s"
	genConfigPayloadErrStringFormat = "construct config update payload failed, %s"
)

func (cc *ChainClient) GetChainConfig() (*config.ChainConfig, error) {
	cc.logger.Debug("[SDK] begin to get chain config")

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_GET_CHAIN_CONFIG.String(), nil, 0)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if err := utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, err
	}

	chainConfig := &config.ChainConfig{}
	err = proto.Unmarshal(resp.ContractResult.Result, chainConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshal contract result failed, %s", err.Error())
	}

	return chainConfig, nil
}

func (cc *ChainClient) GetChainConfigByBlockHeight(blockHeight uint64) (*config.ChainConfig, error) {
	cc.logger.Debugf("[SDK] begin to get chain config by block height [%d]", blockHeight)

	var pairs = []*common.KeyValuePair{{
		Key:   keyBlockHeight,
		Value: utils.U64ToBytes(blockHeight),
	}}

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_GET_CHAIN_CONFIG_AT.String(), pairs, 0)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err := utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	chainConfig := &config.ChainConfig{}
	err = proto.Unmarshal(resp.ContractResult.Result, chainConfig)
	if err != nil {
		return nil, fmt.Errorf("GetChainConfigByBlockHeight unmarshal contract result failed, %s", err)
	}

	return chainConfig, nil
}

func (cc *ChainClient) GetChainConfigSequence() (uint64, error) {
	cc.logger.Debug("[SDK] begin to get chain config sequence")

	chainConfig, err := cc.GetChainConfig()
	if err != nil {
		return 0, err
	}

	return chainConfig.Sequence, nil
}

func (cc *ChainClient) SignChainConfigPayload(payloadBytes []byte) ([]byte, error) {
	signature, err := utils.SignPayloadBytes(cc.privateKey, cc.userCrt, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, "SignChainConfigPayload", err)
	}

	return signature, nil
}

//func (cc *ChainClient) MergeChainConfigSignedPayload(signedPayloadBytes [][]byte) ([]byte, error) {
//	return mergeSystemContractSignedPayload(signedPayloadBytes)
//}

func (cc *ChainClient) CreateChainConfigCoreUpdatePayload(txSchedulerTimeout, txSchedulerValidateTimeout int) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [CoreUpdate] to be signed payload")

	if txSchedulerTimeout > 60 {
		return nil, fmt.Errorf("[tx_scheduler_timeout] should be [0,60]")
	}

	if txSchedulerValidateTimeout > 60 {
		return nil, fmt.Errorf("[tx_scheduler_validate_timeout] should be [0,60]")
	}

	pairs := make([]*common.KeyValuePair, 0)
	if txSchedulerTimeout > 0 {
		pairs = append(pairs, &common.KeyValuePair{
			Key:   "tx_scheduler_timeout",
			Value: utils.IntToBytes(txSchedulerTimeout),
		})
	}

	if txSchedulerValidateTimeout > 0 {
		pairs = append(pairs, &common.KeyValuePair{
			Key:   "tx_scheduler_validate_timeout",
			Value: utils.IntToBytes(txSchedulerValidateTimeout),
		})
	}

	if len(pairs) == 0 {
		return nil, fmt.Errorf("update nothing")
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_CORE_UPDATE.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigBlockUpdatePayload(txTimestampVerify bool, txTimeout, blockTxCapacity,
	blockSize, blockInterval int) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [BlockUpdate] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   "tx_timestamp_verify",
			Value: []byte(strconv.FormatBool(txTimestampVerify)),
		},
	}

	if txTimeout < 600 {
		return nil, fmt.Errorf("[tx_timeout] should be [600, +∞)")
	}

	if blockTxCapacity < 1 {
		return nil, fmt.Errorf("[block_tx_capacity] should be (0, +∞]")
	}

	if blockSize < 1 {
		return nil, fmt.Errorf("[block_size] should be (0, +∞]")
	}

	if blockInterval < 10 {
		return nil, fmt.Errorf("[block_interval] should be [10, +∞]")
	}

	if txTimeout > 0 {
		pairs = append(pairs, &common.KeyValuePair{
			Key:   "tx_timeout",
			Value: utils.IntToBytes(txTimeout),
		})
	}
	if blockTxCapacity > 0 {
		pairs = append(pairs, &common.KeyValuePair{
			Key:   "block_tx_capacity",
			Value: utils.IntToBytes(blockTxCapacity),
		})
	}
	if blockSize > 0 {
		pairs = append(pairs, &common.KeyValuePair{
			Key:   "block_size",
			Value: utils.IntToBytes(blockSize),
		})
	}
	if blockInterval > 0 {
		pairs = append(pairs, &common.KeyValuePair{
			Key:   "block_interval",
			Value: utils.IntToBytes(blockInterval),
		})
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_BLOCK_UPDATE.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigTrustRootAddPayload(trustRootOrgId, trustRootCrt string) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [TrustRootAdd] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   keyOrgId,
			Value: []byte(trustRootOrgId),
		},
		{
			Key:   "root",
			Value: []byte(trustRootCrt),
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_TRUST_ROOT_ADD.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigTrustRootUpdatePayload(trustRootOrgId, trustRootCrt string) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [TrustRootUpdate] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   keyOrgId,
			Value: []byte(trustRootOrgId),
		},
		{
			Key:   "root",
			Value: []byte(trustRootCrt),
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_TRUST_ROOT_UPDATE.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigTrustRootDeletePayload(trustRootOrgId string) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [TrustRootDelete] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   keyOrgId,
			Value: []byte(trustRootOrgId),
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_TRUST_ROOT_DELETE.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigPermissionAddPayload(permissionResourceName string, policy *accesscontrol.Policy) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [PermissionAdd] to be signed payload")

	policyBytes, err := proto.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("marshal policy failed, %s", err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   permissionResourceName,
			Value: policyBytes,
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_PERMISSION_ADD.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigPermissionUpdatePayload(permissionResourceName string, policy *accesscontrol.Policy) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [PermissionUpdate] to be signed payload")

	policyBytes, err := proto.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("marshal policy failed, %s", err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   permissionResourceName,
			Value: policyBytes,
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_PERMISSION_UPDATE.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigPermissionDeletePayload(permissionResourceName string) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [PermissionDelete] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key: permissionResourceName,
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_PERMISSION_DELETE.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigConsensusNodeIdAddPayload(nodeOrgId string, nodeIds []string) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeAddrAdd] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   keyOrgId,
			Value: []byte(nodeOrgId),
		},
		{
			Key:   keyNodeIds,
			Value: []byte(strings.Join(nodeIds, ",")),
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_NODE_ID_ADD.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigConsensusNodeIdUpdatePayload(nodeOrgId, nodeOldIds, nodeNewIds string) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeAddrUpdate] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   keyOrgId,
			Value: []byte(nodeOrgId),
		},
		{
			Key:   keyNodeId,
			Value: []byte(nodeOldIds),
		},
		{
			Key:   keyNewNodeId,
			Value: []byte(nodeNewIds),
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_NODE_ID_UPDATE.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigConsensusNodeIdDeletePayload(nodeOrgId, nodeId string) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeAddrDelete] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   keyOrgId,
			Value: []byte(nodeOrgId),
		},
		{
			Key:   keyNodeId,
			Value: []byte(nodeId),
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_NODE_ID_DELETE.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigConsensusNodeOrgAddPayload(nodeOrgId string, nodeIds []string) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeOrgAdd] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   keyOrgId,
			Value: []byte(nodeOrgId),
		},
		{
			Key:   keyNodeIds,
			Value: []byte(strings.Join(nodeIds, ",")),
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_NODE_ORG_ADD.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigConsensusNodeOrgUpdatePayload(nodeOrgId string, nodeIds []string) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeOrgUpdate] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   keyOrgId,
			Value: []byte(nodeOrgId),
		},
		{
			Key:   keyNodeIds,
			Value: []byte(strings.Join(nodeIds, ",")),
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_NODE_ORG_UPDATE.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigConsensusNodeOrgDeletePayload(nodeOrgId string) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeOrgAdd] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   keyOrgId,
			Value: []byte(nodeOrgId),
		},
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_NODE_ORG_DELETE.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigConsensusExtAddPayload(kvs []*common.KeyValuePair) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusExtAdd] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_CONSENSUS_EXT_ADD.String(), kvs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigConsensusExtUpdatePayload(kvs []*common.KeyValuePair) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusExtUpdate] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_CONSENSUS_EXT_UPDATE.String(), kvs, seq+1)

	return payload, nil
}

func (cc *ChainClient) CreateChainConfigConsensusExtDeletePayload(keys []string) (*common.Payload, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusExtDelete] to be signed payload")

	var pairs = make([]*common.KeyValuePair, len(keys))
	for i, key := range keys {
		pairs[i] = &common.KeyValuePair{
			Key: key,
		}
	}

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_CHAIN_CONFIG.String(),
		syscontract.ChainConfigFunction_CONSENSUS_EXT_DELETE.String(), pairs, seq+1)

	return payload, nil
}

func (cc *ChainClient) SendChainConfigUpdateRequest(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	return cc.sendContractRequest(payload, endorsers, timeout, withSyncResult)
}
