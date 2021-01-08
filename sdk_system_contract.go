/**
 * @Author: zghh
 * @Date:   2020-12-02 10:09:05
 **/
package chainmaker_sdk_go

import (
	"fmt"
	"strconv"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/golang/protobuf/proto"
)

const (
	SYSTEM_CHAIN = "system_chain"
)

func (cc ChainClient) GetTxByTxId(txId string) (*pb.TransactionInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]",
		pb.QueryFunction_GET_TX_BY_TX_ID.String(), txId)

	payloadBytes, err := constructQueryPayload(
		pb.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		pb.QueryFunction_GET_TX_BY_TX_ID.String(),
		[]*pb.KeyValuePair{
			{
				Key:   "txId",
				Value: txId,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_SYSTEM_CONTRACT, txId, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	transactionInfo := &pb.TransactionInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, transactionInfo); err != nil {
		return nil, fmt.Errorf("unmarshal transaction info payload failed, %s", err.Error())
	}

	return transactionInfo, nil
}

func (cc ChainClient) GetBlockByHeight(blockHeight int64, withRWSet bool) (*pb.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHeight:%d]/[withRWSet:%s]",
		pb.QueryFunction_GET_BLOCK_BY_HEIGHT.String(), blockHeight, strconv.FormatBool(withRWSet))

	payloadBytes, err := constructQueryPayload(
		pb.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		pb.QueryFunction_GET_BLOCK_BY_HEIGHT.String(),
		[]*pb.KeyValuePair{
			{
				Key:   "blockHeight",
				Value: strconv.FormatInt(blockHeight, 10),
			},
			{
				Key:   "withRWSet",
				Value: strconv.FormatBool(withRWSet),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	blockInfo := &pb.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("unmarshal block info payload failed, %s", err.Error())
	}

	return blockInfo, nil

}

func (cc ChainClient) GetBlockByHash(blockHash string, withRWSet bool) (*pb.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[blockHash:%s]/[withRWSet:%s]",
		pb.QueryFunction_GET_BLOCK_BY_HASH.String(), blockHash, strconv.FormatBool(withRWSet))

	payloadBytes, err := constructQueryPayload(
		pb.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		pb.QueryFunction_GET_BLOCK_BY_HASH.String(),
		[]*pb.KeyValuePair{
			{
				Key:   "blockHash",
				Value: blockHash,
			},
			{
				Key:   "withRWSet",
				Value: strconv.FormatBool(withRWSet),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	blockInfo := &pb.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("unmarshal block info payload failed, %s", err.Error())
	}

	return blockInfo, nil

}

func (cc ChainClient) GetBlockByTxId(txId string, withRWSet bool) (*pb.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[txId:%s]/[withRWSet:%s]",
		pb.QueryFunction_GET_BLOCK_BY_TX_ID.String(), txId, strconv.FormatBool(withRWSet))

	payloadBytes, err := constructQueryPayload(
		pb.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		pb.QueryFunction_GET_BLOCK_BY_TX_ID.String(),
		[]*pb.KeyValuePair{
			{
				Key:   "txId",
				Value: txId,
			},
			{
				Key:   "withRWSet",
				Value: strconv.FormatBool(withRWSet),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	blockInfo := &pb.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("unmarshal block info payload failed, %s", err.Error())
	}

	return blockInfo, nil
}

func (cc ChainClient) GetLastConfigBlock(withRWSet bool) (*pb.BlockInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]/[withRWSet:%s]",
		pb.QueryFunction_GET_LAST_CONFIG_BLOCK.String(), strconv.FormatBool(withRWSet))

	payloadBytes, err := constructQueryPayload(
		pb.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		pb.QueryFunction_GET_LAST_CONFIG_BLOCK.String(),
		[]*pb.KeyValuePair{
			{
				Key:   "withRWSet",
				Value: strconv.FormatBool(withRWSet),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	blockInfo := &pb.BlockInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, blockInfo); err != nil {
		return nil, fmt.Errorf("unmarshal block info payload failed, %s", err.Error())
	}

	return blockInfo, nil
}

func (cc ChainClient) GetChainInfo() (*pb.ChainInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		pb.QueryFunction_GET_CHAIN_INFO.String())

	payloadBytes, err := constructQueryPayload(
		pb.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		pb.QueryFunction_GET_CHAIN_INFO.String(),
		[]*pb.KeyValuePair{},
	)
	if err != nil {
		return nil, fmt.Errorf("marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	chainInfo := &pb.ChainInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, chainInfo); err != nil {
		return nil, fmt.Errorf("unmarshal chain info payload failed, %s", err.Error())
	}

	return chainInfo, nil
}

func (cc ChainClient) GetContractInfo() (*pb.ContractInfo, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		pb.QueryFunction_GET_CONTRACT_INFO.String())

	payloadBytes, err := constructQueryPayload(
		pb.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		pb.QueryFunction_GET_CONTRACT_INFO.String(),
		[]*pb.KeyValuePair{},
	)
	if err != nil {
		return nil, fmt.Errorf("marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	contractInfo := &pb.ContractInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, contractInfo); err != nil {
		return nil, fmt.Errorf("unmarshal contract info payload failed, %s", err.Error())
	}

	return contractInfo, nil
}

func (cc ChainClient) GetNodeChainList() (*pb.ChainList, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		pb.QueryFunction_GET_NODE_CHAIN_LIST.String())

	payloadBytes, err := constructQueryPayload(
		pb.ContractName_SYSTEM_CONTRACT_QUERY.String(),
		pb.QueryFunction_GET_NODE_CHAIN_LIST.String(),
		[]*pb.KeyValuePair{},
	)
	if err != nil {
		return nil, fmt.Errorf("marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = CheckProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	chainList := &pb.ChainList{}
	if err = proto.Unmarshal(resp.ContractResult.Result, chainList); err != nil {
		return nil, fmt.Errorf("unmarshal chain list payload failed, %s", err.Error())
	}

	return chainList, nil
}
