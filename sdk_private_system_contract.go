/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"fmt"
	"github.com/golang/protobuf/proto"
)

func (cc *ChainClient) SaveCert(enclaveCert, enclaveId, txId string, withSyncResult bool, timeout int64) (*common.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Infof("[SDK] begin to save cert , [contract:%s]/[method:%s]/[txId:%s]",
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_SAVE_CERT.String(),
		txId,
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string]string{
		"enclave_cert": enclaveCert,
		"enclave_id":   enclaveId,
	})

	payloadBytes, err := constructSystemContractPayload(
		cc.chainId,
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_SAVE_CERT.String(),
		pairs,
		defaultSequence,
	)
	if err != nil {
		return nil, fmt.Errorf("construct save cert payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequestWithTimeout(common.TxType_INVOKE_SYSTEM_CONTRACT, txId, payloadBytes, timeout)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	if resp.Code == common.TxStatusCode_SUCCESS {
		if !withSyncResult {
			resp.ContractResult = &common.ContractResult{
				Code:    common.ContractResultCode_OK,
				Message: common.ContractResultCode_OK.String(),
				Result:  []byte(txId),
			}
		} else {
			contractResult, err := cc.getSyncResult(txId)
			if err != nil {
				return nil, fmt.Errorf("get sync result failed, %s", err.Error())
			}

			if contractResult.Code != common.ContractResultCode_OK {
				resp.Code = common.TxStatusCode_CONTRACT_FAIL
				resp.Message = contractResult.Message
			}

			resp.ContractResult = contractResult
		}
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) SaveDir(orderId, txId string,
	privateDir *common.StrSlice, withSyncResult bool, timeout int64) (*common.TxResponse, error) {

	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Infof("[SDK] begin to save dir , [contract:%s]/[method:%s]/[txId:%s]",
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_SAVE_DIR.String(),
		txId,
	)

	// 构造Payload
	priDirBytes, err := privateDir.Marshal()
	if err != nil {
		return nil, fmt.Errorf("serielized private dir failed, %s", err.Error())
	}

	pairs := paramsMap2KVPairs(map[string]string{
		"order_id":    orderId,
		"private_dir": string(priDirBytes),
	})

	payloadBytes, err := constructSystemContractPayload(
		cc.chainId,
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_SAVE_DIR.String(),
		pairs,
		defaultSequence,
	)
	if err != nil {
		return nil, fmt.Errorf("construct save dir payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequestWithTimeout(common.TxType_INVOKE_SYSTEM_CONTRACT, txId, payloadBytes, timeout)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	if resp.Code == common.TxStatusCode_SUCCESS {
		if !withSyncResult {
			resp.ContractResult = &common.ContractResult{
				Code:    common.ContractResultCode_OK,
				Message: common.ContractResultCode_OK.String(),
				Result:  []byte(txId),
			}
		} else {
			contractResult, err := cc.getSyncResult(txId)
			if err != nil {
				return nil, fmt.Errorf("get sync result failed, %s", err.Error())
			}

			if contractResult.Code != common.ContractResultCode_OK {
				resp.Code = common.TxStatusCode_CONTRACT_FAIL
				resp.Message = contractResult.Message
			}

			resp.ContractResult = contractResult
		}
	}

	if resp.Code != common.TxStatusCode_SUCCESS || resp.Message != "OK" {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) GetContract(contractName, codeHash string) (*common.PrivateGetContract, error) {

	cc.logger.Infof("[SDK] begin to get contract , [contract:%s]/[method:%s]",
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_GET_CONTRACT.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string]string{
		"contract_name": contractName,
		"code_hash":     codeHash,
	})

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_GET_CONTRACT.String(),
		pairs,
	)
	if err != nil {
		return nil, fmt.Errorf("marshal get contract payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	contractInfo := &common.PrivateGetContract{}
	if err = proto.Unmarshal(resp.ContractResult.Result, contractInfo); err != nil {
		return nil, fmt.Errorf("GetContract unmarshal contract info payload failed, %s", err.Error())
	}

	return contractInfo, nil
}

func (cc *ChainClient) SaveData(contractName string, contractVersion string, codeHash []byte, reportHash []byte, result *common.ContractResult, txId string, rwSet *common.TxRWSet, sign []byte,
	events *common.StrSlice, withSyncResult bool, timeout int64) (*common.TxResponse, error) { //todo   change params   return TxResponse
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Infof("[SDK] begin to save data , [contract:%s]/[method:%s]/[txId:%s]",
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_SAVE_DATA.String(),
		txId,
	)

	// 构造Payload
	var rwSetStr string
	if rwSet != nil {
		rwb, err := rwSet.Marshal()
		if err != nil {
			return nil, fmt.Errorf("construct save data payload failed, %s", err.Error())
		}
		rwSetStr = string(rwb)
	}

	var eventsStr string
	if events != nil {
		eb, err := events.Marshal()
		if err != nil {
			return nil, fmt.Errorf("construct save data payload failed, %s", err.Error())
		}
		eventsStr = string(eb)
	}

	var resultStr string
	if result != nil {
		result, err := result.Marshal()
		if err != nil {
			return nil, fmt.Errorf("construct save data payload failed, %s", err.Error())
		}
		resultStr = string(result)
	}

	pairs := paramsMap2KVPairs(map[string]string{
		"result":        resultStr,
		"contract_name": contractName,
		"version":       contractVersion,
		"code_hash":     string(codeHash),
		"rw_set":        rwSetStr,
		"events":        eventsStr,
		"report_hash":   string(reportHash),
		"sign":          string(sign),
	})

	payloadBytes, err := constructSystemContractPayload(
		cc.chainId,
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_SAVE_DATA.String(),
		pairs,
		defaultSequence,
	)
	if err != nil {
		return nil, fmt.Errorf("construct save data payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequestWithTimeout(common.TxType_INVOKE_SYSTEM_CONTRACT, txId, payloadBytes, timeout)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	if resp.Code == common.TxStatusCode_SUCCESS {
		if !withSyncResult {
			resp.ContractResult = &common.ContractResult{
				Code:    common.ContractResultCode_OK,
				Message: common.ContractResultCode_OK.String(),
				Result:  []byte(txId),
			}
		} else {
			contractResult, err := cc.getSyncResult(txId)
			if err != nil {
				return nil, fmt.Errorf("get sync result failed, %s", err.Error())
			}

			if contractResult.Code != common.ContractResultCode_OK {
				resp.Code = common.TxStatusCode_CONTRACT_FAIL
				resp.Message = contractResult.Message
			}

			resp.ContractResult = contractResult
		}
	}

	if resp.Code != common.TxStatusCode_SUCCESS || resp.Message != "OK" {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) GetData(contractName, key string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get data , [contract:%s]/[method:%s]",
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_GET_DATA.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string]string{
		"contract_name": contractName,
		"key":           key,
	})

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_GET_DATA.String(),
		pairs,
	)
	if err != nil {
		return nil, fmt.Errorf("marshal get data payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) SaveContract(codeBytes []byte, codeHash, contractName, version, txId string,
	withSyncResult bool, timeout int64) (*common.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Infof("[SDK] begin to save contract code , [contract:%s]/[method:%s]/[txId:%s]",
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_SAVE_CONTRACT.String(),
		txId,
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string]string{
		"contract_code": string(codeBytes),
		"code_hash":     codeHash,
		"contract_name": contractName,
		"version":       version,
	})

	payloadBytes, err := constructSystemContractPayload(
		cc.chainId,
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_SAVE_CONTRACT.String(),
		pairs,
		defaultSequence,
	)
	if err != nil {
		return nil, fmt.Errorf("construct save contract code payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequestWithTimeout(common.TxType_INVOKE_SYSTEM_CONTRACT, txId, payloadBytes, timeout)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	if resp.Code == common.TxStatusCode_SUCCESS {
		if !withSyncResult {
			resp.ContractResult = &common.ContractResult{
				Code:    common.ContractResultCode_OK,
				Message: common.ContractResultCode_OK.String(),
				Result:  []byte(txId),
			}
		} else {
			contractResult, err := cc.getSyncResult(txId)
			if err != nil {
				return nil, fmt.Errorf("get sync result failed, %s", err.Error())
			}

			if contractResult.Code != common.ContractResultCode_OK {
				resp.Code = common.TxStatusCode_CONTRACT_FAIL
				resp.Message = contractResult.Message
			}

			resp.ContractResult = contractResult
		}
	}

	return resp, nil
}

func (cc *ChainClient) UpdateContract(codeBytes []byte, codeHash, contractName, version, txId string,
	withSyncResult bool, timeout int64) (*common.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Infof("[SDK] begin to update contract code , [contract:%s]/[method:%s]/[txId:%s]",
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_UPDATE_CONTRACT.String(),
		txId,
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string]string{
		"contract_code": string(codeBytes),
		"code_hash":     codeHash,
		"contract_name": contractName,
		"version":       version,
	})

	payloadBytes, err := constructSystemContractPayload(
		cc.chainId,
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_UPDATE_CONTRACT.String(),
		pairs,
		defaultSequence,
	)
	if err != nil {
		return nil, fmt.Errorf("construct update contract code payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequestWithTimeout(common.TxType_INVOKE_SYSTEM_CONTRACT, txId, payloadBytes, timeout)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	if resp.Code == common.TxStatusCode_SUCCESS {
		if !withSyncResult {
			resp.ContractResult = &common.ContractResult{
				Code:    common.ContractResultCode_OK,
				Message: common.ContractResultCode_OK.String(),
				Result:  []byte(txId),
			}
		} else {
			contractResult, err := cc.getSyncResult(txId)
			if err != nil {
				return nil, fmt.Errorf("get sync result failed, %s", err.Error())
			}

			if contractResult.Code != common.ContractResultCode_OK {
				resp.Code = common.TxStatusCode_CONTRACT_FAIL
				resp.Message = contractResult.Message
			}

			resp.ContractResult = contractResult
		}
	}

	return resp, nil
}

func (cc *ChainClient) SaveQuote(enclaveId, quoteId, quote, sign, txId string, withSyncResult bool, timeout int64) (*common.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Infof("[SDK] begin to save contract code , [contract:%s]/[method:%s]/[txId:%s]",
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_SAVE_QUOTE.String(),
		txId,
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string]string{
		"enclave_id": enclaveId,
		"quote_id":   quoteId,
		"quote":      quote,
		"sign":       sign,
	})

	payloadBytes, err := constructSystemContractPayload(
		cc.chainId,
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_SAVE_QUOTE.String(),
		pairs,
		defaultSequence,
	)
	if err != nil {
		return nil, fmt.Errorf("construct save quote  payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequestWithTimeout(common.TxType_INVOKE_SYSTEM_CONTRACT, txId, payloadBytes, timeout)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	if resp.Code == common.TxStatusCode_SUCCESS {
		if !withSyncResult {
			resp.ContractResult = &common.ContractResult{
				Code:    common.ContractResultCode_OK,
				Message: common.ContractResultCode_OK.String(),
				Result:  []byte(txId),
			}
		} else {
			contractResult, err := cc.getSyncResult(txId)
			if err != nil {
				return nil, fmt.Errorf("get sync result failed, %s", err.Error())
			}

			if contractResult.Code != common.ContractResultCode_OK {
				resp.Code = common.TxStatusCode_CONTRACT_FAIL
				resp.Message = contractResult.Message
			}

			resp.ContractResult = contractResult
		}
	}

	return resp, nil
}

func (cc *ChainClient) GetCert(enclaveId string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get data , [contract:%s]/[method:%s]",
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_GET_CERT.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string]string{
		"enclave_id": enclaveId,
	})

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_GET_CERT.String(),
		pairs,
	)
	if err != nil {
		return nil, fmt.Errorf("marshal get data payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) GetDir(orderId string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get data , [contract:%s]/[method:%s]",
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_GET_DIR.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string]string{
		"order_id": orderId,
	})

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_GET_DIR.String(),
		pairs,
	)
	if err != nil {
		return nil, fmt.Errorf("marshal get data payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}

func (cc *ChainClient) GetQuote(quoteId string) ([]byte, error) {
	cc.logger.Infof("[SDK] begin to get data , [contract:%s]/[method:%s]",
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_GET_QUOTE.String(),
	)

	// 构造Payload
	pairs := paramsMap2KVPairs(map[string]string{
		"quote_id": quoteId,
	})

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_PRIVATE_COMPUTE.String(),
		common.PrivateComputeContractFunction_GET_QUOTE.String(),
		pairs,
	)
	if err != nil {
		return nil, fmt.Errorf("marshal get data payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	return resp.ContractResult.Result, nil
}
