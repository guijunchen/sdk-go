package chainmaker_sdk_go

import (
	"fmt"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

func (cc *ChainClient) MultiSignContractReq(payload *common.Payload) (*common.TxResponse, error) {

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) MultiSignContractVote(multisignReqPayload *common.Payload,
	endorser *common.EndorsementEntry) (*common.TxResponse, error) {

	msvi := &syscontract.MultiSignVoteInfo{
		Vote:        syscontract.VoteStatus_AGREE,
		Endorsement: endorser,
	}
	msviByte, _ := msvi.Marshal()
	pairs := []*common.KeyValuePair{
		{
			Key:   syscontract.MultiVote_VOTE_INFO.String(),
			Value: msviByte,
		},
		{
			Key:   syscontract.MultiVote_TX_ID.String(),
			Value: []byte(multisignReqPayload.TxId),
		},
	}
	payload := cc.CreateContractMultiSignVotePayload(pairs)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) MultiSignContractQuery(txId string) (*common.TxResponse, error) {

	pairs := []*common.KeyValuePair{
		{
			Key:   syscontract.MultiVote_TX_ID.String(),
			Value: []byte(txId),
		},
	}
	payload := cc.CreateContractMultiSignQueryPayload(pairs)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	return resp, nil
}

func (cc *ChainClient) CreateContractMultiSignReqPayload(pairs []*common.KeyValuePair) *common.Payload {

	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_MULTI_SIGN.String(),
		syscontract.MultiSignFunction_REQ.String(), pairs, defaultSeq)
	return payload
}

func (cc *ChainClient) CreateContractMultiSignVotePayload(pairs []*common.KeyValuePair) *common.Payload {
	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_MULTI_SIGN.String(),
		syscontract.MultiSignFunction_VOTE.String(), pairs, defaultSeq)

	return payload
}

func (cc *ChainClient) CreateContractMultiSignQueryPayload(pairs []*common.KeyValuePair) *common.Payload {
	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_MULTI_SIGN.String(),
		syscontract.MultiSignFunction_QUERY.String(), pairs, defaultSeq)

	return payload
}
