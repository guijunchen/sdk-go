/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gogo/protobuf/proto"

	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/pb-go/discovery"
	"chainmaker.org/chainmaker/pb-go/store"
	"chainmaker.org/chainmaker/pb-go/syscontract"
	"chainmaker.org/chainmaker/sdk-go/utils"
)

func (cc *ChainClient) GetTxByTxId(txId string) (*common.TransactionInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]",
		syscontract.ChainQueryFunction_GET_TX_BY_TX_ID, txId)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_TX_BY_TX_ID.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractTxId,
				Value: []byte(txId),
			},
		}, 0,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	transactionInfo := &common.TransactionInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, transactionInfo); err != nil {
		return nil, fmt.Errorf("GetTxByTxId unmarshal transaction info payload failed, %s", err)
	}

	return transactionInfo, nil
}

func (cc *ChainClient) GetBlockByHeight(blockHeight uint64, withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHeight:%d]/[withRWSet:%t]",
		syscontract.ChainQueryFunction_GET_BLOCK_BY_HEIGHT, blockHeight, withRWSet)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_BLOCK_BY_HEIGHT.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractBlockHeight,
				Value: []byte(strconv.FormatUint(blockHeight, 10)),
			},
		}, 0,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetBlockByHeight unmarshal block info payload failed, %s", err)
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetBlockByHash(blockHash string, withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHash:%s]/[withRWSet:%t]",
		syscontract.ChainQueryFunction_GET_BLOCK_BY_HASH, blockHash, withRWSet)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_BLOCK_BY_HASH.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractBlockHash,
				Value: []byte(blockHash),
			},
			{
				Key:   utils.KeyBlockContractWithRWSet,
				Value: []byte(strconv.FormatBool(withRWSet)),
			},
		}, 0,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetBlockByHash unmarshal block info payload failed, %s", err)
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetBlockByTxId(txId string, withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]/[withRWSet:%t]",
		syscontract.ChainQueryFunction_GET_BLOCK_BY_TX_ID, txId, withRWSet)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_BLOCK_BY_TX_ID.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractTxId,
				Value: []byte(txId),
			},
			{
				Key:   utils.KeyBlockContractWithRWSet,
				Value: []byte(strconv.FormatBool(withRWSet)),
			},
		}, 0,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetBlockByTxId unmarshal block info payload failed, %s", err)
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetLastConfigBlock(withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[withRWSet:%t]",
		syscontract.ChainQueryFunction_GET_LAST_CONFIG_BLOCK, withRWSet)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_LAST_CONFIG_BLOCK.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractWithRWSet,
				Value: []byte(strconv.FormatBool(withRWSet)),
			},
		}, 0,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetBlockByTxId unmarshal block info payload failed, %s", err)
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetChainInfo() (*discovery.ChainInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]", syscontract.ChainQueryFunction_GET_CHAIN_INFO)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_CHAIN_INFO.String(), []*common.KeyValuePair{}, 0,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	chainInfo := &discovery.ChainInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, chainInfo); err != nil {
		return nil, fmt.Errorf("GetChainInfo unmarshal chain info payload failed, %s", err)
	}

	return chainInfo, nil
}

func (cc *ChainClient) GetNodeChainList() (*discovery.ChainList, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		syscontract.ChainQueryFunction_GET_NODE_CHAIN_LIST)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_NODE_CHAIN_LIST.String(), []*common.KeyValuePair{}, 0,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	chainList := &discovery.ChainList{}
	if err = proto.Unmarshal(resp.ContractResult.Result, chainList); err != nil {
		return nil, fmt.Errorf("GetNodeChainList unmarshal chain list payload failed, %s", err)
	}

	return chainList, nil
}

func (cc *ChainClient) GetFullBlockByHeight(blockHeight uint64) (*store.BlockWithRWSet, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHeight:%d]",
		syscontract.ChainQueryFunction_GET_FULL_BLOCK_BY_HEIGHT, blockHeight)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_FULL_BLOCK_BY_HEIGHT.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractBlockHeight,
				Value: []byte(strconv.FormatUint(blockHeight, 10)),
			},
		}, 0,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		if utils.IsArchived(resp.Code) {
			return nil, errors.New(resp.Code.String())
		}
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	fullBlockInfo := &store.BlockWithRWSet{}
	if err = proto.Unmarshal(resp.ContractResult.Result, fullBlockInfo); err != nil {
		return nil, fmt.Errorf("GetFullBlockByHeight unmarshal block info payload failed, %s", err)
	}

	return fullBlockInfo, nil
}

func (cc *ChainClient) GetArchivedBlockHeight() (uint64, error) {
	return cc.getBlockHeight("", "")
}

func (cc *ChainClient) GetBlockHeightByTxId(txId string) (uint64, error) {
	return cc.getBlockHeight(txId, "")
}

func (cc *ChainClient) GetBlockHeightByHash(blockHash string) (uint64, error) {
	return cc.getBlockHeight("", blockHash)
}

func (cc *ChainClient) getBlockHeight(txId, blockHash string) (uint64, error) {
	var (
		txType       = common.TxType_QUERY_CONTRACT
		contractName = syscontract.SystemContract_CHAIN_QUERY.String()
		method       string
		pairs        []*common.KeyValuePair
	)

	if txId != "" {
		method = syscontract.ChainQueryFunction_GET_BLOCK_HEIGHT_BY_TX_ID.String()
		pairs = []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractTxId,
				Value: []byte(txId),
			},
		}

		cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]", method, txId)
	} else if blockHash != "" {
		method = syscontract.ChainQueryFunction_GET_BLOCK_HEIGHT_BY_HASH.String()
		pairs = []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractBlockHash,
				Value: []byte(blockHash),
			},
		}

		cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHash:%s]", method, blockHash)
	} else {
		method = syscontract.ChainQueryFunction_GET_ARCHIVED_BLOCK_HEIGHT.String()
		pairs = []*common.KeyValuePair{}

		cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]", method)
	}

	payload := cc.createPayload("", txType, contractName, method, pairs, 0)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return 0, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return 0, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockHeight, err := strconv.ParseUint(string(resp.ContractResult.Result), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s, parse block height failed, %s", payload.TxType, err)
	}

	return blockHeight, nil
}

func (cc *ChainClient) GetLastBlock(withRWSet bool) (*common.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[withRWSet:%t]",
		syscontract.ChainQueryFunction_GET_LAST_BLOCK, withRWSet)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_LAST_BLOCK.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractWithRWSet,
				Value: []byte(strconv.FormatBool(withRWSet)),
			},
		}, 0,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockInfo := &common.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("GetLastBlock unmarshal block info payload failed, %s", err)
	}

	return blockInfo, nil
}

func (cc *ChainClient) GetCurrentBlockHeight() (uint64, error) {
	block, err := cc.GetLastBlock(false)
	if err != nil {
		return 0, err
	}

	return block.Block.Header.BlockHeight, nil
}

func (cc *ChainClient) GetBlockHeaderByHeight(blockHeight uint64) (*common.BlockHeader, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHeight:%d]",
		syscontract.ChainQueryFunction_GET_BLOCK_HEADER_BY_HEIGHT, blockHeight)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CHAIN_QUERY.String(),
		syscontract.ChainQueryFunction_GET_BLOCK_HEADER_BY_HEIGHT.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyBlockContractBlockHeight,
				Value: []byte(strconv.FormatUint(blockHeight, 10)),
			},
		}, 0,
	)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	blockHeader := &common.BlockHeader{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockHeader); err != nil {
		return nil, fmt.Errorf("GetBlockHeaderByHeight unmarshal block header payload failed, %s", err)
	}

	return blockHeader, nil
}

func (cc *ChainClient) InvokeSystemContract(contractName, method, txId string, params []*common.KeyValuePair,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	cc.logger.Debugf("[SDK] begin to INVOKE system contract, [contractName:%s]/[method:%s]/[txId:%s]/[params:%+v]",
		contractName, method, txId, params)

	payload := cc.createPayload(txId, common.TxType_INVOKE_CONTRACT, contractName, method, params, 0)

	return cc.sendContractRequest(payload, nil, timeout, withSyncResult)
}

func (cc *ChainClient) QuerySystemContract(contractName, method string, params []*common.KeyValuePair,
	timeout int64) (*common.TxResponse, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [contractName:%s]/[method:%s]/[params:%+v]",
		contractName, method, params)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, contractName, method, params, 0)

	resp, err := cc.proposalRequestWithTimeout(payload, nil, timeout)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	return resp, nil
}
