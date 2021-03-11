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
	"strconv"
	"strings"
)

const (
	orgId = "org_id"
	addrs = "addresses"
)

const (
	getCCSeqErrStringFormat         = "get chain config sequence failed, %s"
	genConfigPayloadErrStringFormat = "construct config update payload failed, %s"
)

func (cc ChainClient) GetChainConfig() (*config.ChainConfig, error) {
	cc.logger.Debug("[SDK] begin to get chain config")

	pairs := make([]*common.KeyValuePair, 0)
	payloadBytes, err := constructQueryPayload(common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_GET_CHAIN_CONFIG.String(), pairs)
	if err != nil {
		return nil, fmt.Errorf("construct query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, "", payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err := checkProposalRequestResp(resp, true); err != nil {
		return nil, err
	}

	chainConfig := &config.ChainConfig{}
	err = proto.Unmarshal(resp.ContractResult.Result, chainConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshal contract result failed, %s", err.Error())
	}

	return chainConfig, nil
}

func (cc ChainClient) GetChainConfigByBlockHeight(blockHeight int) (*config.ChainConfig, error) {
	cc.logger.Debugf("[SDK] begin to get chain config by block height [%d]", blockHeight)

	pairs := make([]*common.KeyValuePair, 0)
	pairs = append(pairs, &common.KeyValuePair{
		Key:   "block_height",
		Value: strconv.Itoa(blockHeight),
	})

	payloadBytes, err := constructQueryPayload(common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_GET_CHAIN_CONFIG_AT.String(), pairs)
	if err != nil {
		return nil, fmt.Errorf("construct query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, "", payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("get chain config by block height %s failed, %s", common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err := checkProposalRequestResp(resp, true); err != nil {
		return nil, err
	}

	chainConfig := &config.ChainConfig{}
	err = proto.Unmarshal(resp.ContractResult.Result, chainConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshal contract result failed, %s", err.Error())
	}

	return chainConfig, nil
}

func (cc ChainClient) GetChainConfigSequence() (int, error) {
	cc.logger.Debug("[SDK] begin to get chain config sequence")

	chainConfig, err := cc.GetChainConfig()
	if err != nil {
		return -1, err
	}

	return int(chainConfig.Sequence), nil
}

func (cc ChainClient) SignChainConfigPayload(payloadBytes []byte) ([]byte, error) {
	payload := &common.SystemContractPayload{}
	if err := proto.Unmarshal(payloadBytes, payload); err != nil {
		return nil, fmt.Errorf("unmarshal config update payload failed, %s", err)
	}

	signBytes, err := signPayload(cc.privateKey, cc.userCrt, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("SignPayload failed, %s", err)
	}

	// TODO: 后续支持证书索引，减小交易大小
	sender := &accesscontrol.SerializedMember{
		OrgId:      cc.orgId,
		MemberInfo: cc.userCrtPEM,
		IsFullCert: true,
	}

	entry := &common.EndorsementEntry{
		Signer:    sender,
		Signature: signBytes,
	}

	payload.Endorsement = []*common.EndorsementEntry{
		entry,
	}

	signedPayloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal config update sigend payload failed, %s", err)
	}

	return signedPayloadBytes, nil
}

func (cc ChainClient) MergeChainConfigSignedPayload(signedPayloadBytes [][]byte) ([]byte, error) {
	return mergeSystemContractSignedPayload(signedPayloadBytes)
}

func (cc ChainClient) CreateChainConfigCoreUpdatePayload(txSchedulerTimeout, txSchedulerValidateTimeout int) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [CoreUpdate] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

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
			Value: strconv.Itoa(txSchedulerTimeout),
		})
	}

	if txSchedulerValidateTimeout > 0 {
		pairs = append(pairs, &common.KeyValuePair{
			Key:   "tx_scheduler_validate_timeout",
			Value: strconv.Itoa(txSchedulerValidateTimeout),
		})
	}

	if len(pairs) == 0 {
		return nil, fmt.Errorf("update nothing")
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_CORE_UPDATE.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigBlockUpdatePayload(txTimestampVerify bool, txTimeout, blockTxCapacity, blockSize, blockInterval int) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [BlockUpdate] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   "tx_timestamp_verify",
			Value: strconv.FormatBool(txTimestampVerify),
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
			Value: strconv.Itoa(txTimeout),
		})
	}
	if blockTxCapacity > 0 {
		pairs = append(pairs, &common.KeyValuePair{
			Key:   "block_tx_capacity",
			Value: strconv.Itoa(blockTxCapacity),
		})
	}
	if blockSize > 0 {
		pairs = append(pairs, &common.KeyValuePair{
			Key:   "block_size",
			Value: strconv.Itoa(blockSize),
		})
	}
	if blockInterval > 0 {
		pairs = append(pairs, &common.KeyValuePair{
			Key:   "block_interval",
			Value: strconv.Itoa(blockInterval),
		})
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_BLOCK_UPDATE.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigTrustRootAddPayload(trustRootOrgId, trustRootCrt string) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [TrustRootAdd] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   orgId,
			Value: trustRootOrgId,
		},
		{
			Key:   "root",
			Value: trustRootCrt,
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_TRUST_ROOT_ADD.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigTrustRootUpdatePayload(trustRootOrgId, trustRootCrt string) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [TrustRootUpdate] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   orgId,
			Value: trustRootOrgId,
		},
		{
			Key:   "root",
			Value: trustRootCrt,
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_TRUST_ROOT_UPDATE.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigTrustRootDeletePayload(trustRootOrgId string) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [TrustRootDelete] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   orgId,
			Value: trustRootOrgId,
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_TRUST_ROOT_DELETE.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigPermissionAddPayload(permissionResourceName string, policy *accesscontrol.Policy) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [PermissionAdd] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	bytes, err := proto.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("marshal policy failed, %s", err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   permissionResourceName,
			Value: string(bytes),
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_PERMISSION_ADD.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigPermissionUpdatePayload(permissionResourceName string, policy *accesscontrol.Policy) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [PermissionUpdate] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	bytes, err := proto.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("marshal policy failed, %s", err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   permissionResourceName,
			Value: string(bytes),
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_PERMISSION_UPDATE.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigPermissionDeletePayload(permissionResourceName string) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [PermissionDelete] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key: permissionResourceName,
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_PERMISSION_DELETE.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigConsensusNodeAddrAddPayload(nodeOrgId string, nodeAddresses []string) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeAddrAdd] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   orgId,
			Value: nodeOrgId,
		},
		{
			Key:   addrs,
			Value: strings.Join(nodeAddresses, ","),
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_NODE_ADDR_ADD.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigConsensusNodeAddrUpdatePayload(nodeOrgId, nodeOldAddress, nodeNewAddress string) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeAddrUpdate] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   orgId,
			Value: nodeOrgId,
		},
		{
			Key:   "address",
			Value: nodeOldAddress,
		},
		{
			Key:   "new_address",
			Value: nodeNewAddress,
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_NODE_ADDR_UPDATE.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigConsensusNodeAddrDeletePayload(nodeOrgId, nodeAddress string) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeAddrDelete] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   orgId,
			Value: nodeOrgId,
		},
		{
			Key:   "address",
			Value: nodeAddress,
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_NODE_ADDR_DELETE.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigConsensusNodeOrgAddPayload(nodeOrgId string, nodeAddresses []string) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeOrgAdd] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   orgId,
			Value: nodeOrgId,
		},
		{
			Key:   addrs,
			Value: strings.Join(nodeAddresses, ","),
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_NODE_ORG_ADD.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigConsensusNodeOrgUpdatePayload(nodeOrgId string, nodeAddresses []string) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeOrgUpdate] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   orgId,
			Value: nodeOrgId,
		},
		{
			Key:   addrs,
			Value: strings.Join(nodeAddresses, ","),
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_NODE_ORG_UPDATE.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigConsensusNodeOrgDeletePayload(nodeOrgId string) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusNodeOrgAdd] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   orgId,
			Value: nodeOrgId,
		},
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_NODE_ORG_DELETE.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigConsensusExtAddPayload(kvs []*common.KeyValuePair) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusExtAdd] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_CONSENSUS_EXT_ADD.String(), kvs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigConsensusExtUpdatePayload(kvs []*common.KeyValuePair) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusExtUpdate] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_CONSENSUS_EXT_UPDATE.String(), kvs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) CreateChainConfigConsensusExtDeletePayload(keys []string) ([]byte, error) {
	cc.logger.Debug("[SDK] begin to create [ConsensusExtDelete] to be signed payload")

	seq, err := cc.GetChainConfigSequence()
	if err != nil {
		return nil, fmt.Errorf(getCCSeqErrStringFormat, err)
	}

	pairs := []*common.KeyValuePair{}
	for _, key := range keys {
		pairs = append(pairs, &common.KeyValuePair{
			Key: key,
		})
	}

	payload, err := constructConfigUpdatePayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		common.ConfigFunction_CONSENSUS_EXT_DELETE.String(), pairs, seq+1)
	if err != nil {
		return nil, fmt.Errorf(genConfigPayloadErrStringFormat, err)
	}

	return payload, nil
}

func (cc ChainClient) SendChainConfigUpdateRequest(mergeSignedPayloadBytes []byte) (*common.TxResponse, error) {
	txId := GetRandTxId()

	resp, err := cc.proposalRequest(common.TxType_UPDATE_CHAIN_CONFIG, txId, mergeSignedPayloadBytes)
	if err != nil {
		return resp, fmt.Errorf("send %s failed, %s", common.TxType_UPDATE_CHAIN_CONFIG.String(), err.Error())
	}

	resp.ContractResult = &common.ContractResult{
		Code:    common.ContractResultCode_OK,
		Message: common.ContractResultCode_OK.String(),
		Result:  []byte(txId),
	}

	return resp, nil
}
