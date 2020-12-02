/**
 * @Author: jasonruan
 * @Date:   2020-12-02 10:30:23
 **/
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
)

func (cc ChainClient) GetChainConfig() (*pb.ChainConfig, error) {
	cc.logger.Debug("[SDK] begin to get chain config")

	pairs := make([]*pb.KeyValuePair, 0)
	payloadBytes, err := constructQueryPayload(pb.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		pb.ConfigFunction_GET_CHAIN_CONFIG.String(), pairs)
	if err != nil {
		return nil, fmt.Errorf("construct query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_SYSTEM_CONTRACT, "", payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err := checkProposalRequestResp(resp, true); err != nil {
		return nil, err
	}

	chainConfig := &pb.ChainConfig{}
	err = proto.Unmarshal(resp.ContractResult.Result, chainConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshal contract result failed, %s", err.Error())
	}

	return chainConfig, nil
}

func (cc ChainClient) GetChainConfigByBlockHeight(blockHeight int) (*pb.ChainConfig, error) {
	cc.logger.Debugf("[SDK] begin to get chain config by block height [%d]", blockHeight)

	pairs := make([]*pb.KeyValuePair, 0)
	pairs = append(pairs, &pb.KeyValuePair{
		Key:   "block_height",
		Value: strconv.Itoa(blockHeight),
	})

	payloadBytes, err := constructQueryPayload(pb.ContractName_SYSTEM_CONTRACT_CHAIN_CONFIG.String(),
		pb.ConfigFunction_GET_CHAIN_CONFIG_AT.String(), pairs)
	if err != nil {
		return nil, fmt.Errorf("construct query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_SYSTEM_CONTRACT, "", payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err := checkProposalRequestResp(resp, true); err != nil {
		return nil, err
	}

	chainConfig := &pb.ChainConfig{}
	err = proto.Unmarshal(resp.ContractResult.Result, chainConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshal contract result failed, %s", err.Error())
	}

	return chainConfig, nil
}
