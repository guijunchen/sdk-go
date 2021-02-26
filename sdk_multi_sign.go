/**
 * @Author: jasonruan
 * @Date:   2020-12-30 16:46:01
 **/
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
)

func (cc ChainClient) SignMultiSignPayload(payloadBytes []byte) (*pb.EndorsementEntry, error) {
	signBytes, err := signPayload(cc.privateKey, cc.userCrt, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("MultiSignPayloadSign failed, %s", err)
	}

	sender := &pb.SerializedMember{
		OrgId:      cc.orgId,
		MemberInfo: cc.userCrtPEM,
		IsFullCert: true,
	}

	entry := &pb.EndorsementEntry{
		Signer:    sender,
		Signature: signBytes,
	}

	return entry, nil
}

func (cc ChainClient) SendMultiSignReq(txType pb.TxType, payloadBytes []byte, endorsementEntry *pb.EndorsementEntry,
	deadlineBlockHeight int, timeout int64) (*pb.TxResponse, error) {
	var (
		resp pb.TxResponse
		errMsg string
	)

	multiSignReqPayload, err := cc.createMultiSignReqPayload(txType, payloadBytes, endorsementEntry, deadlineBlockHeight)
	if err != nil {
		errMsg = fmt.Sprintf("create multi sign req payload failed, %s", err.Error())
		cc.logger.Error(errMsg)
		resp.Code = pb.TxStatusCode_INVALID_PARAMETER
		resp.Message = errMsg
		return &resp, nil
	}

	return cc.sendContractRequest(pb.TxType_SYSTEM_CONTRACT, multiSignReqPayload, timeout, false)
}

func (cc ChainClient) SendMultiSignVote(voteStatus pb.VoteStatus, multiSignReqTxId, payloadHash string,
	endorsementEntry *pb.EndorsementEntry, timeout int64) (*pb.TxResponse, error) {

	var (
		resp pb.TxResponse
		errMsg string
	)

	multiSignVotePayload, err := cc.createMultiSignVotePayload(voteStatus, multiSignReqTxId, payloadHash, endorsementEntry)
	if err != nil {
		errMsg = fmt.Sprintf("create multi sign vote payload failed, %s", err.Error())
		cc.logger.Error(errMsg)
		resp.Code = pb.TxStatusCode_INVALID_PARAMETER
		resp.Message = errMsg
		return &resp, nil
	}

	return cc.sendContractRequest(pb.TxType_SYSTEM_CONTRACT, multiSignVotePayload, timeout, false)
}

func (cc ChainClient) QueryMultiSignResult(multiSignReqTxId, payloadHash string) (*pb.TxResponse, error) {
	var (
		resp pb.TxResponse
		errMsg string
	)

	multiSignVotePayload, err := cc.createQueryMultiSignResultPayload(multiSignReqTxId, payloadHash)
	if err != nil {
		errMsg = fmt.Sprintf("create query multi sign result payload failed, %s", err.Error())
		cc.logger.Error(errMsg)
		resp.Code = pb.TxStatusCode_INVALID_PARAMETER
		resp.Message = errMsg
		return &resp, nil
	}

	return cc.sendContractRequest(pb.TxType_SYSTEM_CONTRACT, multiSignVotePayload, -1, false)
}

func (cc ChainClient) createMultiSignReqPayload(txType pb.TxType, payloadBytes []byte,
	endorsementEntry *pb.EndorsementEntry, deadlineBlockHeight int) ([]byte, error) {

	voteInfo := &pb.MultSignVoteInfo{
		Vote:        pb.VoteStatus_AGREE,
		Endorsement: endorsementEntry,
	}
	voteInfoBytes, err := proto.Marshal(voteInfo)
	if err != nil {
		return nil, fmt.Errorf("marshal vote info failed, %s", err.Error())
	}

	// 构造Payload
	pairs := []*pb.KeyValuePair{
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

	payload, err := constructSystemContractPayload(cc.chainId, pb.ContractName_SYSTEM_CONTRACT_MULT_SIGN.String(),
		pb.MultSignFunction_REQ.String(), pairs, defaultSequence)
	if err != nil {
		return nil, fmt.Errorf("constructSystemContractPayload failed, %s", err.Error())
	}

	return payload, nil
}

func (cc ChainClient) createMultiSignVotePayload(voteStatus pb.VoteStatus, multiSignReqTxId, payloadHash string,
	endorsementEntry *pb.EndorsementEntry) ([]byte, error) {

	var voteInfo *pb.MultSignVoteInfo
	if voteStatus == pb.VoteStatus_AGREE {
		voteInfo = &pb.MultSignVoteInfo{
			Vote:        pb.VoteStatus_AGREE,
			Endorsement: endorsementEntry,
		}
	} else {
		// 不同意时，不需要用户签名
		voteInfo = &pb.MultSignVoteInfo{
			Vote: pb.VoteStatus_DISAGREE,
		}
	}
	voteInfoBytes, err := proto.Marshal(voteInfo)
	if err != nil {
		return nil, fmt.Errorf("marshal vote info failed, %s", err.Error())
	}


	// 构造Payload
	pairs := []*pb.KeyValuePair{
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

	payload, err := constructSystemContractPayload(cc.chainId, pb.ContractName_SYSTEM_CONTRACT_MULT_SIGN.String(),
		pb.MultSignFunction_VOTE.String(), pairs, defaultSequence)
	if err != nil {
		return nil, fmt.Errorf("createMultiSignVotePayload failed, %s", err.Error())
	}

	return payload, nil
}

func (cc ChainClient) createQueryMultiSignResultPayload(multiSignReqTxId, payloadHash string) ([]byte, error) {
	// 构造Payload
	pairs := []*pb.KeyValuePair{
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

	payload, err := constructSystemContractPayload(cc.chainId, pb.ContractName_SYSTEM_CONTRACT_MULT_SIGN.String(),
		pb.MultSignFunction_VOTE.String(), pairs, defaultSequence)
	if err != nil {
		return nil, fmt.Errorf("createQueryMultiSignResultPayload failed, %s", err.Error())
	}

	return payload, nil
}
