/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/pb-go/discovery"
	"chainmaker.org/chainmaker/pb-go/store"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
)

const (
	SYSTEM_CHAIN   = "system_chain"
	keyWithRWSet   = "withRWSet"
	keyBlockHash   = "blockHash"
	keyBlockHeight = "blockHeight"
	keyTxId        = "txId"
)

func (cc *ChainClient) GetTxByTxId(txId string) (*common.TransactionInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]",
		common.QueryFunction_GET_TX_BY_TX_ID.String(), txId)

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		common.QueryFunction_GET_TX_BY_TX_ID.String(),
		[]*common.KeyValuePair{
			{
				Key:   keyTxId,
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
		if IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
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
				Key:   keyBlockHeight,
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
		if IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
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
				Key:   keyBlockHash,
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
		if IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
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
				Key:   keyTxId,
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
		if IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
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

func (cc *ChainClient) GetFullBlockByHeight(blockHeight int64) (*store.BlockWithRWSet, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHeight:%d]",
		common.QueryFunction_GET_FULL_BLOCK_BY_HEIGHT.String(), blockHeight)

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		common.QueryFunction_GET_FULL_BLOCK_BY_HEIGHT.String(),
		[]*common.KeyValuePair{
			{
				Key:   keyBlockHeight,
				Value: strconv.FormatInt(blockHeight, 10),
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("GetFullBlockByHeight marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		if IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}

		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	fullBlockInfo := &store.BlockWithRWSet{}
	if err = proto.Unmarshal(resp.ContractResult.Result, fullBlockInfo); err != nil {
		return nil, fmt.Errorf("GetFullBlockByHeight unmarshal block info payload failed, %s", err.Error())
	}

	return fullBlockInfo, nil
}

func (cc *ChainClient) GetArchivedBlockHeight() (int64, error) {
	return cc.getBlockHeight("", "")
}

func (cc *ChainClient) GetBlockHeightByTxId(txId string) (int64, error) {
	return cc.getBlockHeight(txId, "")
}

func (cc *ChainClient) GetBlockHeightByHash(blockHash string) (int64, error) {
	return cc.getBlockHeight("", blockHash)
}

func (cc *ChainClient) getBlockHeight(txId, blockHash string) (int64, error) {
	var (
		contractName string
		method       string
		pairs        []*common.KeyValuePair
	)

	contractName = common.ContractName_SYSTEM_CONTRACT_QUERY.String()
	if txId != "" {
		method = common.QueryFunction_GET_BLOCK_HEIGHT_BY_TX_ID.String()
		pairs = []*common.KeyValuePair{
			{
				Key:   keyTxId,
				Value: txId,
			},
		}

		cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]", method, txId)
	} else if blockHash != "" {
		method = common.QueryFunction_GET_BLOCK_HEIGHT_BY_HASH.String()
		pairs = []*common.KeyValuePair{
			{
				Key:   keyBlockHash,
				Value: blockHash,
			},
		}

		cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHash:%s]", method, blockHash)
	} else {
		method = common.QueryFunction_GET_ARCHIVED_BLOCK_HEIGHT.String()
		pairs = []*common.KeyValuePair{}

		cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]", method)
	}

	payloadBytes, err := constructQueryPayload(contractName, method, pairs)
	if err != nil {
		return -1, fmt.Errorf("%s marshal query payload failed, %s", method, err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, txId, payloadBytes)
	if err != nil {
		return -1, fmt.Errorf("%s, proposal request failed, %s", common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return -1, fmt.Errorf("%s, check resp faield, %s", common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	blockHeight, err := strconv.ParseInt(string(resp.ContractResult.Result), 10, 64)
	if err != nil {
		return -1, fmt.Errorf("%s, parse block height failed, %s", common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	return blockHeight, nil
}

func (cc *ChainClient) GetLastBlock(withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[withRWSet:%s]",
		common.QueryFunction_GET_LAST_BLOCK.String(), strconv.FormatBool(withRWSet))

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		common.QueryFunction_GET_LAST_BLOCK.String(),
		[]*common.KeyValuePair{
			{
				Key:   keyWithRWSet,
				Value: strconv.FormatBool(withRWSet),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("GetLastBlock marshal query payload failed, %s", err.Error())
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
		return nil, fmt.Errorf("GetLastBlock unmarshal block info payload failed, %s", err.Error())
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetCurrentBlockHeight() (int64, error) {
	block, err := cc.GetLastBlock(false)
	if err != nil {
		return -1, err
	}

	return block.Block.Header.BlockHeight, nil
}

func (cc *ChainClient) GetBlockHeaderByHeight(blockHeight int64) (*common.BlockHeader, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHeight:%d]",
		common.QueryFunction_GET_BLOCK_HEADER_BY_HEIGHT.String(), blockHeight)

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		common.QueryFunction_GET_BLOCK_HEADER_BY_HEIGHT.String(),
		[]*common.KeyValuePair{
			{
				Key:   keyBlockHeight,
				Value: strconv.FormatInt(blockHeight, 10),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("GetBlockHeaderByHeight marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	blockHeader := &common.BlockHeader{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockHeader); err != nil {
		return nil, fmt.Errorf("GetBlockHeaderByHeight unmarshal block header payload failed, %s", err.Error())
	}

	return blockHeader, nil
}
