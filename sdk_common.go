/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package chainmaker_sdk_go

import (
	"fmt"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
)

const (
	// DefaultRetryLimit 默认轮训交易结果最大次数
	DefaultRetryLimit = 10
	// DefaultRetryInterval 默认每次轮训交易结果时的等待时间，单位ms
	DefaultRetryInterval = 500
	// defaultSeq default sequence
	defaultSeq = 0
)

func (cc *ChainClient) GetSyncResult(txId string) (*common.Result, error) {
	if cc.enableTxResultDispatcher {
		r, err := cc.asyncTxResult(txId)
		if err != nil {
			return nil, err
		}
		return r.Result, nil
	}
	return cc.pollingTxResult(txId)
}

func (cc *ChainClient) GetSyncResultV2(txId string) (*txResult, error) {
	if cc.enableTxResultDispatcher {
		return cc.asyncTxResult(txId)
	}
	r, err := cc.pollingTxResult(txId)
	if err != nil {
		return nil, err
	}
	return &txResult{Result: r}, nil
}

func (cc *ChainClient) pollingTxResult(txId string) (*common.Result, error) {
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
		strategy.Wait(time.Duration(cc.retryInterval)*time.Millisecond),
		strategy.Limit(uint(cc.retryLimit)),
	)

	if err != nil {
		return nil, fmt.Errorf("get tx by txId [%s] failed, %s", txId, err.Error())
	}

	if txInfo == nil || txInfo.Transaction == nil || txInfo.Transaction.Result == nil {
		return nil, fmt.Errorf("get result by txId [%s] failed, %+v", txId, txInfo)
	}

	return txInfo.Transaction.Result, nil
}

func (cc *ChainClient) asyncTxResult(txId string) (*txResult, error) {
	txResultC := cc.txResultDispatcher.register(txId)
	defer cc.txResultDispatcher.unregister(txId)

	timeout := time.Duration(cc.retryInterval*cc.retryLimit) * time.Millisecond
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()
	select {
	case r := <-txResultC:
		return r, nil
	case <-ticker.C:
		return nil, fmt.Errorf("get transaction result timed out, timeout=%s", timeout)
	}
}

func (cc *ChainClient) sendContractRequest(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {

	resp, err := cc.proposalRequestWithTimeout(payload, endorsers, timeout)
	if err != nil {
		return resp, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	if resp.Code == common.TxStatusCode_SUCCESS {
		if withSyncResult {
			result, err := cc.GetSyncResult(payload.TxId)
			if err != nil {
				return nil, fmt.Errorf("get sync result failed, %s", err.Error())
			}
			resp.Code = result.Code
			resp.Message = result.Message
			resp.ContractResult = result.ContractResult
			resp.TxId = payload.TxId
		}
	}

	return resp, nil
}

type TxResponse struct {
	Response    *common.TxResponse
	TxTimestamp int64
	BlockHeight uint64
}

func (cc *ChainClient) sendContractRequestV2(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*TxResponse, error) {

	resp, err := cc.proposalRequestWithTimeout(payload, endorsers, timeout)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", payload.TxType.String(), err.Error())
	}

	txResp := &TxResponse{Response: resp}
	if resp.Code == common.TxStatusCode_SUCCESS {
		if withSyncResult {
			result, err := cc.GetSyncResultV2(payload.TxId)
			if err != nil {
				return nil, fmt.Errorf("get sync result failed, %s", err.Error())
			}
			resp.Code = result.Result.Code
			resp.Message = result.Result.Message
			resp.ContractResult = result.Result.ContractResult
			resp.TxId = payload.TxId
			txResp.TxTimestamp = result.TxTimestamp
			txResp.BlockHeight = result.BlockHeight
		}
	}

	return txResp, nil
}

func (cc *ChainClient) CreatePayload(txId string, txType common.TxType, contractName, method string,
	kvs []*common.KeyValuePair, seq uint64, limit *common.Limit) *common.Payload {
	if txId == "" {
		if cc.enableNormalKey {
			txId = utils.GetRandTxId()
		} else {
			txId = utils.GetTimestampTxId()
		}
	}

	payload := utils.NewPayload(
		utils.WithChainId(cc.chainId),
		utils.WithTxType(txType),
		utils.WithTxId(txId),
		utils.WithTimestamp(time.Now().Unix()),
		utils.WithContractName(contractName),
		utils.WithMethod(method),
		utils.WithParameters(kvs),
		utils.WithSequence(seq),
		utils.WithLimit(limit),
	)

	return payload
}

func (cc *ChainClient) SignPayload(payload *common.Payload) (*common.EndorsementEntry, error) {
	var (
		sender    *accesscontrol.Member
		signBytes []byte
		err       error
	)
	if cc.authType == PermissionedWithCert {

		hashalgo, err := bcx509.GetHashFromSignatureAlgorithm(cc.userCrt.SignatureAlgorithm)
		if err != nil {
			return nil, fmt.Errorf("invalid algorithm: %s", err.Error())
		}

		signBytes, err = utils.SignPayloadWithHashType(cc.privateKey, hashalgo, payload)
		if err != nil {
			return nil, fmt.Errorf("SignPayload failed, %s", err)
		}

		sender = &accesscontrol.Member{
			OrgId:      cc.orgId,
			MemberInfo: cc.userCrtBytes,
			MemberType: accesscontrol.MemberType_CERT,
		}

	} else {
		signBytes, err = utils.SignPayloadWithHashType(cc.privateKey, crypto.HashAlgoMap[cc.hashType], payload)
		if err != nil {
			return nil, fmt.Errorf("SignPayload failed, %s", err.Error())
		}
		sender = &accesscontrol.Member{
			OrgId:      cc.orgId,
			MemberInfo: cc.pkBytes,
			MemberType: accesscontrol.MemberType_PUBLIC_KEY,
		}
	}

	entry := &common.EndorsementEntry{
		Signer:    sender,
		Signature: signBytes,
	}
	return entry, nil
}
