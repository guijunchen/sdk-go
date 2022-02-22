/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package chainmaker_sdk_go

import (
	"fmt"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/gogo/protobuf/proto"
)

const (
	KEY_ALIAS        = "alias"
	KEY_CERT         = "cert"
	KEY_BLOCK_HEIGHT = "block_height"
)

func (cc *ChainClient) AddAlias() (*common.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to add alias, [contract:%s]/[method:%s]",
		syscontract.SystemContract_CERT_MANAGE.String(), syscontract.CertManageFunction_CERT_ALIAS_ADD.String())

	kvs := []*common.KeyValuePair{
		{
			Key:   KEY_ALIAS,
			Value: []byte(cc.alias),
		},
	}

	payload := cc.CreateCertManagePayload(syscontract.CertManageFunction_CERT_ALIAS_ADD.String(), kvs)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	resp.ContractResult = &common.ContractResult{
		Code:   utils.SUCCESS,
		Result: []byte(cc.alias),
	}

	return resp, nil
}

func (cc *ChainClient) QueryCurrentAlias(alias string) (*common.AliasInfo, error) {
	cc.logger.Infof("[SDK] begin to query cert by alias, [contract:%s]/[method:%s]",
		syscontract.SystemContract_CERT_MANAGE.String(), syscontract.CertManageFunction_CERTS_ALIAS_QUERY.String())

	kvs := []*common.KeyValuePair{
		{
			Key:   KEY_ALIAS,
			Value: []byte(cc.alias),
		},
	}

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CERT_MANAGE.String(),
		syscontract.CertManageFunction_CERTS_ALIAS_QUERY.String(), kvs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	aliasInfo := &common.AliasInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, aliasInfo); err != nil {
		return nil, fmt.Errorf("unmarshal cert infos payload failed, %s", err.Error())
	}

	return aliasInfo, nil
}

func (cc *ChainClient) QueryCertByAliasAndBlockHeight(alias, blockHeight string) (*common.AliasCertInfo, error) {
	cc.logger.Infof("[SDK] begin to query cert by alias, [contract:%s]/[method:%s]",
		syscontract.SystemContract_CERT_MANAGE.String(),
		syscontract.CertManageFunction_CERTS_ALIAS_QUERY.String())

	kvs := []*common.KeyValuePair{
		{
			Key:   KEY_ALIAS,
			Value: []byte(cc.alias),
		},
		{
			Key:   KEY_BLOCK_HEIGHT,
			Value: []byte(blockHeight),
		},
	}

	payload := cc.CreatePayload("", common.TxType_QUERY_CONTRACT, syscontract.SystemContract_CERT_MANAGE.String(),
		syscontract.CertManageFunction_CERTS_ALIAS_QUERY.String(), kvs, defaultSeq, nil)

	resp, err := cc.proposalRequest(payload, nil)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, payload.TxType.String(), err.Error())
	}

	aliasCertInfo := &common.AliasCertInfo{}
	if err = proto.Unmarshal(resp.ContractResult.Result, aliasCertInfo); err != nil {
		return nil, fmt.Errorf("unmarshal cert infos payload failed, %s", err.Error())
	}

	return aliasCertInfo, nil
}

func (cc *ChainClient) CreateUpdateAliasPayload(alias, certPEM string) *common.Payload {
	cc.logger.Debugf("[SDK] create [UpdateAlias] to be signed payload")

	pairs := []*common.KeyValuePair{
		{
			Key:   KEY_ALIAS,
			Value: []byte(alias),
		},
		{
			Key:   KEY_CERT,
			Value: []byte(certPEM),
		},
	}

	return cc.CreateCertManagePayload(syscontract.CertManageFunction_CERT_ALIAS_UPDATE.String(), pairs)
}

func (cc *ChainClient) SignUpdateAliasPayload(payload *common.Payload) (*common.EndorsementEntry, error) {
	return cc.SignCertManagePayload(payload)
}

func (cc *ChainClient) UpdateAlias(payload *common.Payload, endorsers []*common.EndorsementEntry,
	timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	return cc.SendCertManageRequest(payload, endorsers, timeout, withSyncResult)
}
