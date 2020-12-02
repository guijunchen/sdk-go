/**
 * @Author: jasonruan
 * @Date:   2020-12-02 10:08:52
 **/
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"fmt"
	"github.com/golang/protobuf/proto"
)

func (cc ChainClient) ContractCreate(txId string, multiSignPayload []byte) (*pb.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Infof("[SDK] begin to CREATE contract, [txId:%s]/[payload size:%d]",
		txId, len(multiSignPayload))

	resp, err := cc.proposalRequest(pb.TxType_CREATE_USER_CONTRACT, txId, multiSignPayload)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_CREATE_USER_CONTRACT.String(), err.Error())
	}

	resp.ContractResult = &pb.ContractResult{
		Code: pb.ContractResultCode_OK,
		Message: pb.ContractResultCode_OK.String(),
		Result: []byte(txId),
	}

	return resp, nil
}

func (cc ChainClient) ContractInvoke(contractName, method, txId string, params map[string]string) (*pb.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Infof("[SDK] begin to INVOKE contract, [contractName:%s]/[method:%s]/[txId:%s]/[params:%+v]",
		contractName, method, txId, params)

	pairs := paramsMap2KVPairs(params)

	payload := &pb.TransactPayload{
		ContractName: contractName,
		Method:       method,
		Parameters:   pairs,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal transact payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_INVOKE_USER_CONTRACT, txId, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_INVOKE_USER_CONTRACT.String(), err.Error())
	}


	resp.ContractResult = &pb.ContractResult{
		Code: pb.ContractResultCode_OK,
		Message: pb.ContractResultCode_OK.String(),
		Result: []byte(txId),
	}

	return resp, nil
}

func (cc ChainClient) ContractQuery(contractName, method, txId string, params map[string]string) (*pb.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Infof("[SDK] begin to QUERY contract, [contractName:%s]/[method:%s]/[txId:%s]/[params:%+v]",
		contractName, method, txId, params)

	pairs := paramsMap2KVPairs(params)

	payload := &pb.QueryPayload{
		ContractName: contractName,
		Method:       method,
		Parameters:   pairs,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_USER_CONTRACT, txId, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_USER_CONTRACT.String(), err.Error())
	}

	return resp, nil
}
