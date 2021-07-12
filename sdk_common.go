/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker/pb-go/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/common"
	"chainmaker.org/chainmaker/sdk-go/utils"
	"fmt"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"time"
)

const (
	// 轮训交易结果最大次数
	retryCnt = 10
)

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
		strategy.Backoff(backoff.Fibonacci(retryInterval * time.Millisecond)),
	)

	if err != nil {
		return nil, fmt.Errorf("get tx by txId [%s] failed, %s", txId, err.Error())
	}

	if txInfo == nil || txInfo.Transaction == nil || txInfo.Transaction.Result == nil {
		return nil, fmt.Errorf("get result by txId [%s] failed, %+v", txId, txInfo)
	}

	return txInfo.Transaction.Result.ContractResult, nil
}

func (cc *ChainClient) sendContractRequest(payload *common.Payload, endosers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {

	resp, err := cc.proposalRequestWithTimeout(payload, endosers, timeout)
	if err != nil {
		return resp, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if resp.Code == common.TxStatusCode_SUCCESS {
		if !withSyncResult {
			resp.TxId = payload.TxId
		}
	} else {
		contractResult, err := cc.getSyncResult(payload.TxId)
		if err != nil {
			return nil, fmt.Errorf("get sync result failed, %s", err.Error())
		}
		resp.ContractResult = contractResult
	}

	return resp, nil
}

func (cc *ChainClient) createPayload(txId string, txType common.TxType, contractName, method string, kvs []*common.KeyValuePair) *common.Payload {
	if txId == "" {
		txId = utils.GetRandTxId()
	}

	payload := utils.NewPayload(
		utils.WithChainId(cc.chainId),
		utils.WithTxType(txType),
		utils.WithTxId(txId),
		utils.WithTimestamp(time.Now().Unix()),
		utils.WithContractName(contractName),
		utils.WithMethod(method),
		utils.WithParameters(kvs),
	)

	return payload
}

func (cc *ChainClient) SignPayload(payload *common.Payload) (*common.EndorsementEntry, error) {

	signBytes, err := utils.SignPayload(cc.privateKey, cc.userCrt, payload)
	if err != nil {
		return nil, fmt.Errorf("SignPayload failed, %s", err)
	}

	sender := &accesscontrol.SerializedMember{
		OrgId:      cc.orgId,
		MemberInfo: cc.userCrtBytes,
		MemberType: accesscontrol.MemberType_CERT,
	}

	entry := &common.EndorsementEntry{
		Signer:    sender,
		Signature: signBytes,
	}

	return entry, nil
}