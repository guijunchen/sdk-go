package chainmaker_sdk_go

import (
	"errors"
	"fmt"
	"strconv"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
)

func (cc *ChainClient) CreateSetGasAdminPayload(adminPubKey crypto.PublicKey) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [CreateSetGasAdminPayload] payload")

	bz, err := adminPubKey.Bytes()
	if err != nil {
		return nil, err
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyGasPublicKey,
			Value: bz,
		},
	}

	return cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		syscontract.GasAccountFunction_SET_ADMIN.String(), pairs, defaultSeq, nil), nil
}

func (cc *ChainClient) GetGasAdmin() (string, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		syscontract.GasAccountFunction_GET_ADMIN)

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		syscontract.GasAccountFunction_GET_ADMIN.String(), nil, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return "", fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err = utils.CheckProposalRequestResp(resp, true); err != nil {
		return "", fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	return string(resp.ContractResult.Result), nil
}

func (cc *ChainClient) CreateRechargeGasPayload(rechargeGasList []*syscontract.RechargeGas) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [CreateRechargeGasPayload] payload")

	rechargeGasReq := syscontract.RechargeGasReq{BatchRechargeGas: rechargeGasList}
	bz, err := rechargeGasReq.Marshal()
	if err != nil {
		return nil, err
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyGasBatchRecharge,
			Value: bz,
		},
	}

	return cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		syscontract.GasAccountFunction_RECHARGE_GAS.String(), pairs, defaultSeq, nil), nil
}

func (cc *ChainClient) GetGasBalance(pubKey crypto.PublicKey) (int64, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		syscontract.GasAccountFunction_GET_BALANCE)

	pubKeyBytes, err := pubKey.Bytes()
	if err != nil {
		return 0, err
	}

	payload := cc.createPayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		syscontract.GasAccountFunction_GET_BALANCE.String(), []*common.KeyValuePair{
			{
				Key:   utils.KeyGasBalancePublicKey,
				Value: pubKeyBytes,
			},
		}, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return 0, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err := utils.CheckProposalRequestResp(resp, true); err != nil {
		return 0, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	balance, err := strconv.ParseInt(string(resp.ContractResult.Result), 10, 64)
	if err != nil {
		return 0, fmt.Errorf(errStringFormat, "strconv.ParseInt", err)
	}

	return balance, nil
}

func (cc *ChainClient) CreateRefundGasPayload(pubKey crypto.PublicKey, amount int64) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [CreateRefundGasPayload] payload")

	pubKeyBytes, err := pubKey.Bytes()
	if err != nil {
		return nil, err
	}

	if amount <= 0 {
		return nil, errors.New("amount must > 0")
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyGasChargePublicKey,
			Value: pubKeyBytes,
		},
		{
			Key:   utils.KeyGasChargeGasAmount,
			Value: []byte(strconv.FormatInt(amount, 10)),
		},
	}

	return cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		syscontract.GasAccountFunction_REFUND_GAS.String(), pairs, defaultSeq, nil), nil
}

func (cc *ChainClient) CreateFrozenGasAccountPayload(pubKey crypto.PublicKey) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [CreateFrozenGasAccountPayload] payload")

	return cc.createFrozenUnfrozenGasAccountPayload(syscontract.GasAccountFunction_FROZEN_ACCOUNT.String(), pubKey)
}

func (cc *ChainClient) CreateUnfrozenGasAccountPayload(pubKey crypto.PublicKey) (*common.Payload, error) {
	cc.logger.Debugf("[SDK] create [CreateFrozenGasAccountPayload] payload")

	return cc.createFrozenUnfrozenGasAccountPayload(syscontract.GasAccountFunction_UNFROZEN_ACCOUNT.String(), pubKey)
}

func (cc *ChainClient) GetGasAccountStatus(pubKey crypto.PublicKey) (bool, error) {
	cc.logger.Debugf("[SDK] begin to QUERY system contract, [method:%s]",
		syscontract.GasAccountFunction_ACCOUNT_STATUS)

	payload, err := cc.createFrozenUnfrozenGasAccountPayload(syscontract.GasAccountFunction_ACCOUNT_STATUS.String(),
		pubKey)
	if err != nil {
		return false, err
	}

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return false, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	if err := utils.CheckProposalRequestResp(resp, true); err != nil {
		return false, fmt.Errorf(errStringFormat, payload.TxType, err)
	}

	return string(resp.ContractResult.Result) == "0", nil
}

func (cc *ChainClient) SendGasManageRequest(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	cc.logger.Debug("[SDK] begin to get chain config")
	return cc.sendContractRequest(payload, endorsers, timeout, withSyncResult)
}

func (cc *ChainClient) createFrozenUnfrozenGasAccountPayload(method string,
	pubKey crypto.PublicKey) (*common.Payload, error) {

	bz, err := pubKey.Bytes()
	if err != nil {
		return nil, err
	}

	pairs := []*common.KeyValuePair{
		{
			Key:   utils.KeyGasFrozenPublicKey,
			Value: bz,
		},
	}

	return cc.createPayload("", common.TxType_INVOKE_CONTRACT, syscontract.SystemContract_ACCOUNT_MANAGER.String(),
		method, pairs, defaultSeq, nil), nil
}
