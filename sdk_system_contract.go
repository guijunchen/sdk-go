/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/pb-go/discovery"
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
)

const (
	SYSTEM_CHAIN = "system_chain"
	keyWithRWSet = "withRWSet"
)

func (cc *ChainClient) GetTxByTxId(txId string) (*common.TransactionInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]",
		common.QueryFunction_GET_TX_BY_TX_ID.String(), txId)

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		common.QueryFunction_GET_TX_BY_TX_ID.String(),
		[]*common.KeyValuePair{
			{
				Key:   "txId",
				Value: txId,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("GetTxByTxId marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, txId, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	transactionInfo := &common.TransactionInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, transactionInfo); err != nil {
		return nil, fmt.Errorf("unmarshal transaction info payload failed, %s", err.Error())
	}

	return transactionInfo, nil
}

func (cc *ChainClient) GetBlockByHeight(blockHeight int64, withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHeight:%d]/[withRWSet:%s]",
		common.QueryFunction_GET_BLOCK_BY_HEIGHT.String(), blockHeight, strconv.FormatBool(withRWSet))

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		common.QueryFunction_GET_BLOCK_BY_HEIGHT.String(),
		[]*common.KeyValuePair{
			{
				Key:   "blockHeight",
				Value: strconv.FormatInt(blockHeight, 10),
			},
			{
				Key:   keyWithRWSet,
				Value: strconv.FormatBool(withRWSet),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("GetBlockByHeight marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetBlockByHeight unmarshal block info payload failed, %s", err.Error())
	}

	return blockInfo, nil

}

func (cc *ChainClient) GetBlockByHash(blockHash string, withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHash:%s]/[withRWSet:%s]",
		common.QueryFunction_GET_BLOCK_BY_HASH.String(), blockHash, strconv.FormatBool(withRWSet))

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		common.QueryFunction_GET_BLOCK_BY_HASH.String(),
		[]*common.KeyValuePair{
			{
				Key:   "blockHash",
				Value: blockHash,
			},
			{
				Key:   keyWithRWSet,
				Value: strconv.FormatBool(withRWSet),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("GetBlockByHash marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetBlockByHash unmarshal block info payload failed, %s", err.Error())
	}

	return blockInfo, nil

}

func (cc *ChainClient) GetBlockByTxId(txId string, withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]/[withRWSet:%s]",
		common.QueryFunction_GET_BLOCK_BY_TX_ID.String(), txId, strconv.FormatBool(withRWSet))

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		common.QueryFunction_GET_BLOCK_BY_TX_ID.String(),
		[]*common.KeyValuePair{
			{
				Key:   "txId",
				Value: txId,
			},
			{
				Key:   keyWithRWSet,
				Value: strconv.FormatBool(withRWSet),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("GetBlockByTxId marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetBlockByTxId unmarshal block info payload failed, %s", err.Error())
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetLastConfigBlock(withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[withRWSet:%s]",
		common.QueryFunction_GET_LAST_CONFIG_BLOCK.String(), strconv.FormatBool(withRWSet))

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		common.QueryFunction_GET_LAST_CONFIG_BLOCK.String(),
		[]*common.KeyValuePair{
			{
				Key:   keyWithRWSet,
				Value: strconv.FormatBool(withRWSet),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("GetLastConfigBlock marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetLastConfigBlock unmarshal block info payload failed, %s", err.Error())
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetChainInfo() (*discovery.ChainInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		common.QueryFunction_GET_CHAIN_INFO.String())

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		common.QueryFunction_GET_CHAIN_INFO.String(),
		[]*common.KeyValuePair{},
	)
	if err != nil {
		return nil, fmt.Errorf("GetChainInfo marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	chainInfo := &discovery.ChainInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, chainInfo); err != nil {
		return nil, fmt.Errorf("unmarshal chain info payload failed, %s", err.Error())
	}

	return chainInfo, nil
}

func (cc *ChainClient) GetNodeChainList() (*discovery.ChainList, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		common.QueryFunction_GET_NODE_CHAIN_LIST.String())

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		common.QueryFunction_GET_NODE_CHAIN_LIST.String(),
		[]*common.KeyValuePair{},
	)
	if err != nil {
		return nil, fmt.Errorf("GetNodeChainList marshar query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	chainList := &discovery.ChainList{}
	if err = proto.Unmarshal(resp.ContractResult.Result, chainList); err != nil {
		return nil, fmt.Errorf("unmarshal chain list payload failed, %s", err.Error())
	}

	return chainList, nil
}
