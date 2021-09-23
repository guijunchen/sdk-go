package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"fmt"
)

const (
	createContractTimeout    = 5
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

var (
	timestamp int64
	txId      string
)

func (cc *ChainClient) MultiSignContractReq(payload *common.Payload) (*common.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to multisign req, [contract:%s]/[method:%s]",
		syscontract.SystemContract_MULTI_SIGN.String(), syscontract.MultiSignFunction_REQ.String())

	//timestamp = payload.Timestamp
	//txId = payload.TxId
	fmt.Println("testMultiSignReq timestamp", payload.Timestamp)
	fmt.Printf("testMultiSignReq txid %s", payload.TxId)
	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = utils.CheckProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	return resp, nil
}

//func (cc *ChainClient) MultiSignContractVote(payload1 *common.Payload,SignKeyPath string,SignCrtPath string) (*common.TxResponse, error) {
func (cc *ChainClient) MultiSignContractVote(payload1 *common.Payload, endorser *common.EndorsementEntry) (*common.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to multisign vote, [contract:%s]/[method:%s]",
		syscontract.SystemContract_MULTI_SIGN.String(), syscontract.MultiSignFunction_VOTE.String())

	//payload1 := cc.createContractMultiSignReqPayload()
	//payload1.Timestamp = timestamp
	//payload1.TxId = txId
	fmt.Println("testMultiSignVote timestamp", payload1.Timestamp)
	fmt.Printf("testMultiSignVote txid %s", payload1.TxId)

	//ee, err := SignPayloadWithPath(SignKeyPath, SignCrtPath, payload1)
	//if err != nil {
	//	return nil, err
	//}

	msvi := &syscontract.MultiSignVoteInfo{
		Vote:        syscontract.VoteStatus_AGREE,
		Endorsement: endorser,
		//Endorsement: endorsers[0],
	}
	msviByte, _ := msvi.Marshal()
	pairs := []*common.KeyValuePair{
		{
			Key:   syscontract.MultiVote_VOTE_INFO.String(),
			Value: msviByte,
		},
		{
			Key:   syscontract.MultiVote_TX_ID.String(),
			Value: []byte(payload1.TxId),
		},
	}
	payload := cc.CreateContractMultiSignVotePayload(syscontract.ContractManageFunction_INIT_CONTRACT.String(), pairs)

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
	cc.logger.Infof("[SDK] begin to multisign vote, [contract:%s]/[method:%s]",
		syscontract.SystemContract_MULTI_SIGN.String(), syscontract.MultiSignFunction_VOTE.String())

	pairs := []*common.KeyValuePair{
		{
			Key:   syscontract.MultiVote_TX_ID.String(),
			Value: []byte(txId),
		},
	}
	payload := cc.CreateContractMultiSignQueryPayload(syscontract.ContractManageFunction_INIT_CONTRACT.String(), pairs)

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

func (cc *ChainClient) CreateContractMultiSignVotePayload(method string, pairs []*common.KeyValuePair) *common.Payload {
	cc.logger.Debugf("[SDK] create ContractMultiSignVotePayload, method: %s", method)
	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_MULTI_SIGN.String(),
		syscontract.MultiSignFunction_VOTE.String(), pairs, defaultSeq)

	return payload
}

func (cc *ChainClient) CreateContractMultiSignQueryPayload(method string, pairs []*common.KeyValuePair) *common.Payload {
	cc.logger.Debugf("[SDK] create ContractMultiSignVotePayload, method: %s", method)
	payload := cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_MULTI_SIGN.String(),
		syscontract.MultiSignFunction_QUERY.String(), pairs, defaultSeq)

	return payload
}
