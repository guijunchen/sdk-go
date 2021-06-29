/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/accesscontrol"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"encoding/hex"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"strconv"
)

func (cc *ChainClient) SignMultiSignPayload(payloadBytes []byte) (*common.EndorsementEntry, error) {
	signBytes, err := signPayload(cc.privateKey, cc.userCrt, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("MultiSignPayloadSign failed, %s", err)
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

	return entry, nil
}

func (cc *ChainClient) SendMultiSignReq(txType common.TxType, payloadBytes []byte, endorsementEntry *common.EndorsementEntry,
	deadlineBlockHeight int, timeout int64) (*common.TxResponse, error) {
	var (
		resp   common.TxResponse
		errMsg string
	)

	multiSignReqPayload, err := cc.createMultiSignReqPayload(txType, payloadBytes, endorsementEntry, deadlineBlockHeight)
	if err != nil {
		errMsg = fmt.Sprintf("create multi sign req payload failed, %s", err.Error())
		cc.logger.Error(errMsg)
		resp.Code = common.TxStatusCode_INVALID_PARAMETER
		resp.Message = errMsg
		return &resp, nil
	}

	return cc.sendContractRequest(common.TxType_INVOKE_SYSTEM_CONTRACT, multiSignReqPayload, timeout, false)
}

func (cc *ChainClient) SendMultiSignVote(voteStatus common.VoteStatus, multiSignReqTxId, payloadHash string,
	endorsementEntry *common.EndorsementEntry, timeout int64) (*common.TxResponse, error) {

	var (
		resp   common.TxResponse
		errMsg string
	)

	multiSignVotePayload, err := cc.createMultiSignVotePayload(voteStatus, multiSignReqTxId, payloadHash, endorsementEntry)
	if err != nil {
		errMsg = fmt.Sprintf("create multi sign vote payload failed, %s", err.Error())
		cc.logger.Error(errMsg)
		resp.Code = common.TxStatusCode_INVALID_PARAMETER
		resp.Message = errMsg
		return &resp, nil
	}

	return cc.sendContractRequest(common.TxType_INVOKE_SYSTEM_CONTRACT, multiSignVotePayload, timeout, false)
}

func (cc *ChainClient) QueryMultiSignResult(multiSignReqTxId, payloadHash string) (*common.TxResponse, error) {
	var (
		resp   common.TxResponse
		errMsg string
	)

	multiSignVotePayload, err := cc.createQueryMultiSignResultPayload(multiSignReqTxId, payloadHash)
	if err != nil {
		errMsg = fmt.Sprintf("create query multi sign result payload failed, %s", err.Error())
		cc.logger.Error(errMsg)
		resp.Code = common.TxStatusCode_INVALID_PARAMETER
		resp.Message = errMsg
		return &resp, nil
	}

	return cc.sendContractRequest(common.TxType_INVOKE_SYSTEM_CONTRACT, multiSignVotePayload, -1, false)
}

func (cc *ChainClient) createMultiSignReqPayload(txType common.TxType, payloadBytes []byte,
	endorsementEntry *common.EndorsementEntry, deadlineBlockHeight int) ([]byte, error) {

	voteInfo := &common.MultiSignVoteInfo{
		Vote:        common.VoteStatus_AGREE,
		Endorsement: endorsementEntry,
	}
	voteInfoBytes, err := proto.Marshal(voteInfo)
	if err != nil {
		return nil, fmt.Errorf("marshal vote info failed, %s", err.Error())
	}

	// 构造Payload
	pairs := []*common.KeyValuePair{
		{
			Key:   "tx_type", // 多签内的交易类型
			Value: txType.String(),
		},
		{
			Key:   "deadline_block", // 过期的区块高度
			Value: strconv.Itoa(deadlineBlockHeight),
		},
		{
			Key:   "payload",
			Value: hex.EncodeToString(payloadBytes),
		},
		{
			Key:   "vote_info",
			Value: hex.EncodeToString(voteInfoBytes),
		},
	}

	payload, err := constructSystemContractPayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_MULTI_SIGN.String(),
		common.MultiSignFunction_REQ.String(), pairs, defaultSequence)
	if err != nil {
		return nil, fmt.Errorf("constructSystemContractPayload failed, %s", err.Error())
	}

	return payload, nil
}

func (cc *ChainClient) createMultiSignVotePayload(voteStatus common.VoteStatus, multiSignReqTxId, payloadHash string,
	endorsementEntry *common.EndorsementEntry) ([]byte, error) {

	var voteInfo *common.MultiSignVoteInfo
	if voteStatus == common.VoteStatus_AGREE {
		voteInfo = &common.MultiSignVoteInfo{
			Vote:        common.VoteStatus_AGREE,
			Endorsement: endorsementEntry,
		}
	} else {
		// 不同意时，不需要用户签名
		voteInfo = &common.MultiSignVoteInfo{
			Vote: common.VoteStatus_DISAGREE,
		}
	}
	voteInfoBytes, err := proto.Marshal(voteInfo)
	if err != nil {
		return nil, fmt.Errorf("marshal vote info failed, %s", err.Error())
	}

	// 构造Payload
	pairs := []*common.KeyValuePair{
		// tx_id或payload_hash，如果有tx_id，会优先选择tx_id
		{
			Key:   "tx_id",
			Value: multiSignReqTxId,
		},
		{
			Key:   "payload_hash",
			Value: payloadHash,
		},
		{
			Key:   "vote_info",
			Value: hex.EncodeToString(voteInfoBytes),
		},
	}

	payload, err := constructSystemContractPayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_MULTI_SIGN.String(),
		common.MultiSignFunction_VOTE.String(), pairs, defaultSequence)
	if err != nil {
		return nil, fmt.Errorf("createMultiSignVotePayload failed, %s", err.Error())
	}

	return payload, nil
}

func (cc *ChainClient) createQueryMultiSignResultPayload(multiSignReqTxId, payloadHash string) ([]byte, error) {
	// 构造Payload
	pairs := []*common.KeyValuePair{
		// tx_id或payload_hash，如果有tx_id，会优先选择tx_id
		{
			Key:   "tx_id",
			Value: multiSignReqTxId,
		},
		{
			Key:   "payload_hash",
			Value: payloadHash,
		},
	}

	payload, err := constructSystemContractPayload(cc.chainId, common.ContractName_SYSTEM_CONTRACT_MULTI_SIGN.String(),
		common.MultiSignFunction_VOTE.String(), pairs, defaultSequence)
	if err != nil {
		return nil, fmt.Errorf("createQueryMultiSignResultPayload failed, %s", err.Error())
	}

	return payload, nil
}
