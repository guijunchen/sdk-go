/**
 * @Author: jasonruan
 * @Date:   2020-12-02 10:08:52
 **/
package chainmaker_sdk_go

import (
	"fmt"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
)

func (cc ChainClient) ContractCreate(txId string, multiSignedPayload []byte) (*pb.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Debugf("[SDK] begin to CREATE contract, [txId:%s]/[payload size:%d]",
		txId, len(multiSignedPayload))

	resp, err := cc.proposalRequest(pb.TxType_CREATE_USER_CONTRACT, txId, multiSignedPayload)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_CREATE_USER_CONTRACT.String(), err.Error())
	}

	resp.ContractResult = &pb.ContractResult{
		Code:    pb.ContractResultCode_OK,
		Message: pb.ContractResultCode_OK.String(),
		Result:  []byte(txId),
	}

	return resp, nil
}

func (cc ChainClient) ContractUpgrade(txId string, multiSignedPayload []byte) (*pb.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Debugf("[SDK] begin to UPGRADE contract, [txId:%s]/[payload size:%d]",
		txId, len(multiSignedPayload))

	resp, err := cc.proposalRequest(pb.TxType_UPGRADE_USER_CONTRACT, txId, multiSignedPayload)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_UPGRADE_USER_CONTRACT.String(), err.Error())
	}

	resp.ContractResult = &pb.ContractResult{
		Code:    pb.ContractResultCode_OK,
		Message: pb.ContractResultCode_OK.String(),
		Result:  []byte(txId),
	}

	return resp, nil
}

func (cc ChainClient) ContractInvoke(contractName, method, txId string, params map[string]string) (*pb.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Debugf("[SDK] begin to INVOKE contract, [contractName:%s]/[method:%s]/[txId:%s]/[params:%+v]",
		contractName, method, txId, params)

	pairs := paramsMap2KVPairs(params)

	payloadBytes, err := constructTransactPayload(contractName, method, pairs)
	if err != nil {
		return nil, fmt.Errorf("construct transact payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_INVOKE_USER_CONTRACT, txId, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_INVOKE_USER_CONTRACT.String(), err.Error())
	}

	resp.ContractResult = &pb.ContractResult{
		Code:    pb.ContractResultCode_OK,
		Message: pb.ContractResultCode_OK.String(),
		Result:  []byte(txId),
	}

	return resp, nil
}

func (cc ChainClient) ContractQuery(contractName, method string, params map[string]string) (*pb.TxResponse, error) {
	txId := GetRandTxId()

	cc.logger.Debugf("[SDK] begin to QUERY contract, [contractName:%s]/[method:%s]/[txId:%s]/[params:%+v]",
		contractName, method, txId, params)

	pairs := paramsMap2KVPairs(params)

	payloadBytes, err := constructQueryPayload(contractName, method, pairs)
	if err != nil {
		return nil, fmt.Errorf("construct query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_USER_CONTRACT, txId, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", pb.TxType_QUERY_USER_CONTRACT.String(), err.Error())
	}

	return resp, nil
}
