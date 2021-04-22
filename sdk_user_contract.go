/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/accesscontrol"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"strings"
	"time"
)

const (
	// 轮训交易结果最大次数
	retryCnt = 10
)

func (cc *ChainClient) CreateContractCreatePayload(contractName, version, byteCode string, runtime common.RuntimeType, kvs []*common.KeyValuePair) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [ContractCreate] to be signed payload")
	return cc.createContractManagePayload(contractName, common.ManageUserContractFunction_INIT_CONTRACT.String(), version, byteCode, runtime, kvs)
}

func (cc *ChainClient) CreateContractUpgradePayload(contractName, version, byteCode string, runtime common.RuntimeType, kvs []*common.KeyValuePair) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [ContractUpgrade] to be signed payload")
	return cc.createContractManagePayload(contractName, common.ManageUserContractFunction_UPGRADE_CONTRACT.String(), version, byteCode, runtime, kvs)
}

func (cc *ChainClient) CreateContractFreezePayload(contractName string) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [ContractFreeze] to be signed payload")
	return cc.createContractOpPayload(contractName, common.ManageUserContractFunction_FREEZE_CONTRACT.String())
}

func (cc *ChainClient) CreateContractUnfreezePayload(contractName string) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [ContractUnfreeze] to be signed payload")
	return cc.createContractOpPayload(contractName, common.ManageUserContractFunction_UNFREEZE_CONTRACT.String())
}

func (cc *ChainClient) CreateContractRevokePayload(contractName string) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [ContractRevoke] to be signed payload")
	return cc.createContractOpPayload(contractName, common.ManageUserContractFunction_REVOKE_CONTRACT.String())
}

func (cc *ChainClient) createContractManagePayload(contractName, method, version, byteCode string, runtime common.RuntimeType, kvs []*common.KeyValuePair) ([]byte, error) {
	var (
		err       error
		codeBytes []byte
	)
	cc.logger.Debugf("[SDK] create [ContractManage] to be signed payload")

	exists := Exists(byteCode)
	if exists {
		codeBytes, err = ioutil.ReadFile(byteCode)
		if err != nil {
			return nil, fmt.Errorf("read from byteCode file %s failed, %s", byteCode, err)
		}
	} else {
		for {
			byteCode = strings.TrimSpace(byteCode)

			codeBytes, err = hex.DecodeString(byteCode)
			if err == nil {
				break
			}

			codeBytes, err = base64.StdEncoding.DecodeString(byteCode)
			if err == nil {
				break
			}

			return nil, fmt.Errorf("decode byteCode failed, %s", err)
		}
	}

	payload := &common.ContractMgmtPayload{
		ChainId: cc.chainId,
		ContractId: &common.ContractId{
			ContractName:    contractName,
			ContractVersion: version,
			RuntimeType:     runtime,
		},
		Method:     method,
		Parameters: kvs,
		ByteCode:   codeBytes,
	}

	bytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("construct contract manage payload failed, %s", err)
	}

	return bytes, nil
}

func (cc *ChainClient) createContractOpPayload(contractName, method string) ([]byte, error) {
	payload := &common.ContractMgmtPayload{
		ChainId: cc.chainId,
		ContractId: &common.ContractId{
			ContractName: contractName,
		},
		Method: method,
	}

	bytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("construct contract manage payload failed, %s", err)
	}

	return bytes, nil
}

func (cc *ChainClient) SignContractManagePayload(payloadBytes []byte) ([]byte, error) {
	payload := &common.ContractMgmtPayload{}
	if err := proto.Unmarshal(payloadBytes, payload); err != nil {
		return nil, fmt.Errorf("unmarshal contract manage payload failed, %s", err)
	}

	signBytes, err := signPayload(cc.privateKey, cc.userCrt, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("SignPayload failed, %s", err)
	}

	sender := &accesscontrol.SerializedMember{
		OrgId:      cc.orgId,
		MemberInfo: cc.userCrtBytes,
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
		return nil, fmt.Errorf("marshal contract manage sigend payload failed, %s", err)
	}

	return signedPayloadBytes, nil
}

func (cc *ChainClient) MergeContractManageSignedPayload(signedPayloadBytes [][]byte) ([]byte, error) {
	return mergeContractManageSignedPayload(signedPayloadBytes)
}

func (cc *ChainClient) SendContractManageRequest(mergeSignedPayloadBytes []byte, timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	return cc.sendContractRequest(common.TxType_MANAGE_USER_CONTRACT, mergeSignedPayloadBytes, timeout, withSyncResult)
}

func (cc *ChainClient) sendContractRequest(txType common.TxType, mergeSignedPayloadBytes []byte, timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	txId := GetRandTxId()

	resp, err := cc.proposalRequestWithTimeout(txType, txId, mergeSignedPayloadBytes, timeout)
	if err != nil {
		return resp, fmt.Errorf("send %s failed, %s", txType.String(), err.Error())
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
			resp.ContractResult = contractResult
		}
	}

	return resp, nil
}

func (cc *ChainClient) InvokeContract(contractName, method, txId string, params map[string]string, timeout int64, withSyncResult bool) (*common.TxResponse, error) {
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

	resp, err := cc.proposalRequestWithTimeout(common.TxType_INVOKE_USER_CONTRACT, txId, payloadBytes, timeout)
	if err != nil {
		return resp, fmt.Errorf("%s failed, %s", common.TxType_INVOKE_USER_CONTRACT.String(), err.Error())
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

func (cc *ChainClient) QueryContract(contractName, method string, params map[string]string, timeout int64) (*common.TxResponse, error) {
	txId := GetRandTxId()

	cc.logger.Debugf("[SDK] begin to QUERY contract, [contractName:%s]/[method:%s]/[txId:%s]/[params:%+v]",
		contractName, method, txId, params)

	pairs := paramsMap2KVPairs(params)

	payloadBytes, err := constructQueryPayload(contractName, method, pairs)
	if err != nil {
		return nil, fmt.Errorf("construct query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequestWithTimeout(common.TxType_QUERY_USER_CONTRACT, txId, payloadBytes, timeout)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", common.TxType_QUERY_USER_CONTRACT.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) getSyncResult(txId string) (*common.ContractResult, error) {
	var (
		txInfo *common.TransactionInfo
		err    error
	)

	err = retry.Retry(func(uint) error {
		txInfo, err = cc.GetTxByTxId(txId)
		if err != nil {
			return err
		}

		return nil
	},
		strategy.Limit(retryCnt),
		strategy.Backoff(backoff.Fibonacci(retryInterval*time.Millisecond)),
	)

	if err != nil {
		return nil, fmt.Errorf("get tx by txId [%s] failed, %s", txId, err.Error())
	}
	if txInfo == nil || txInfo.Transaction == nil || txInfo.Transaction.Result == nil {
		return nil, fmt.Errorf("get result by txId [%s] failed, %+v", txId, txInfo)
	}
	return txInfo.Transaction.Result.ContractResult, nil
}

func (cc *ChainClient) GetTxRequest(contractName, method, txId string, params map[string]string) (*common.TxRequest, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Debugf("[SDK] begin to create TxRequest, [contractName:%s]/[method:%s]/[txId:%s]/[params:%+v]",
		contractName, method, txId, params)

	pairs := paramsMap2KVPairs(params)

	payloadBytes, err := constructTransactPayload(contractName, method, pairs)
	if err != nil {
		return nil, fmt.Errorf("construct transact payload failed, %s", err.Error())
	}

	req, err := cc.generateTxRequest(txId, common.TxType_INVOKE_USER_CONTRACT, payloadBytes)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (cc *ChainClient) SendTxRequest(txRequest *common.TxRequest, timeout int64, withSyncResult bool) (*common.TxResponse, error) {

	resp, err := cc.sendTxRequest(txRequest, timeout)
	if err != nil {
		return resp, fmt.Errorf("%s failed, %s", common.TxType_INVOKE_USER_CONTRACT.String(), err.Error())
	}

	if resp.Code == common.TxStatusCode_SUCCESS {
		if !withSyncResult {
			resp.ContractResult = &common.ContractResult{
				Code:    common.ContractResultCode_OK,
				Message: common.ContractResultCode_OK.String(),
				Result:  []byte(txRequest.Header.TxId),
			}
		} else {
			contractResult, err := cc.getSyncResult(txRequest.Header.TxId)
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

