/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker/pb-go/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/common"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"strings"
)

func (cc *ChainClient) AddCert() (*common.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to add cert, [contract:%s]/[method:%s]",
		common.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(), common.CertManageFunction_CERT_ADD.String())

	certHash, err := cc.GetCertHash()
	if err != nil {
		return nil, fmt.Errorf("get cert hash in hex failed, %s", err.Error())
	}

	payloadBytes, err := constructSystemContractPayload("", common.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(),
		common.CertManageFunction_CERT_ADD.String(), []*common.KeyValuePair{}, defaultSequence)
	if err != nil {
		return nil, fmt.Errorf("construct transact payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_INVOKE_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	resp.ContractResult = &common.ContractResult{
		Code:    common.ContractResultCode_OK,
		Message: common.ContractResultCode_OK.String(),
		Result:  certHash,
	}

	return resp, nil
}

func (cc *ChainClient) DeleteCert(certHashes []string) (*common.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to delete cert, [contract:%s]/[method:%s]",
		common.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(), common.CertManageFunction_CERTS_DELETE.String())

	payloadBytes, err := constructSystemContractPayload(
		"",
		common.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(),
		common.CertManageFunction_CERTS_DELETE.String(),
		[]*common.KeyValuePair{
			{
				Key:   "cert_hashes",
				Value: strings.Join(certHashes, ","),
			},
		},
		defaultSequence,
	)
	if err != nil {
		return nil, fmt.Errorf("marshal transact payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_INVOKE_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return resp, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_INVOKE_SYSTEM_CONTRACT.String(), err.Error())
	}

	resp.ContractResult = &common.ContractResult{
		Code:    common.ContractResultCode_OK,
		Message: common.ContractResultCode_OK.String(),
	}

	return resp, nil
}

func (cc *ChainClient) QueryCert(certHashes []string) (*common.CertInfos, error) {
	cc.logger.Infof("[SDK] begin to query cert, [contract:%s]/[method:%s]",
		common.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(), common.CertManageFunction_CERTS_QUERY.String())

	payloadBytes, err := constructQueryPayload(
		common.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(),
		common.CertManageFunction_CERTS_QUERY.String(),
		[]*common.KeyValuePair{
			{
				Key:   "cert_hashes",
				Value: strings.Join(certHashes, ","),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(common.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf(errStringFormat, common.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	certInfos := &common.CertInfos{}
	if err = proto.Unmarshal(resp.ContractResult.Result, certInfos); err != nil {
		return nil, fmt.Errorf("unmarshal cert infos payload failed, %s", err.Error())
	}

	return certInfos, nil
}

func (cc *ChainClient) GetCertHash() ([]byte, error) {
	chainConfig, err := cc.GetChainConfig()

	if err != nil {
		return nil, fmt.Errorf("get cert hash failed, %s", err.Error())
	}

	member := &accesscontrol.SerializedMember{
		OrgId:      cc.orgId,
		MemberInfo: cc.userCrtBytes,
		IsFullCert: true,
	}

	certHash, err := getCertificateId(member.GetMemberInfo(), chainConfig.Crypto.Hash)
	if err != nil {
		return nil, fmt.Errorf("calc cert hash failed, %s", err.Error())
	}

	return certHash, nil
}

func (cc *ChainClient) CreateCertManagePayload(method string, kvs []*common.KeyValuePair) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [CertManage] to be signed payload")

	payload := &common.SystemContractPayload{
		ChainId:      cc.chainId,
		ContractName: common.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(),
		Method:       method,
		Parameters:   kvs,
	}

	bytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("construct cert manage payload failed, %s", err)
	}

	return bytes, nil
}

func (cc *ChainClient) CreateCertManageFrozenPayload(certs []string) ([]byte, error) {
	pairs := []*common.KeyValuePair{
		{
			Key:   "certs",
			Value: strings.Join(certs, ","),
		},
	}

	return cc.CreateCertManagePayload(common.CertManageFunction_CERTS_FREEZE.String(), pairs)
}

func (cc *ChainClient) CreateCertManageUnfrozenPayload(certs []string) ([]byte, error) {
	pairs := []*common.KeyValuePair{
		{
			Key:   "certs",
			Value: strings.Join(certs, ","),
		},
	}

	return cc.CreateCertManagePayload(common.CertManageFunction_CERTS_UNFREEZE.String(), pairs)
}

func (cc *ChainClient) CreateCertManageRevocationPayload(certCrl string) ([]byte, error) {
	pairs := []*common.KeyValuePair{
		{
			Key:   "cert_crl",
			Value: certCrl,
		},
	}

	return cc.CreateCertManagePayload(common.CertManageFunction_CERTS_REVOKE.String(), pairs)
}

func (cc *ChainClient) SignCertManagePayload(payloadBytes []byte) ([]byte, error) {
	payload := &common.SystemContractPayload{}
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

func (cc *ChainClient) MergeCertManageSignedPayload(signedPayloadBytes [][]byte) ([]byte, error) {
	return mergeSystemContractSignedPayload(signedPayloadBytes)
}

func (cc *ChainClient) SendCertManageRequest(mergeSignedPayloadBytes []byte, timeout int64, withSyncResult bool) (*common.TxResponse, error) {
	return cc.sendContractRequest(common.TxType_INVOKE_SYSTEM_CONTRACT, mergeSignedPayloadBytes, timeout, withSyncResult)
}
