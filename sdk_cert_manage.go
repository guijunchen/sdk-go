/**
 * @Author: zghh
 * @Date:   2020-12-03 10:16:38
 **/
package chainmaker_sdk_go

import (
	"fmt"
	"strings"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/golang/protobuf/proto"
)

const (
	defaultSequence = 0
)

func (cc ChainClient) AddCert() (*pb.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to INVOKE system contract, [contract:%s]/[method:%s]",
		pb.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(), pb.CertManageFunction_CERT_ADD.String())

	certHash, err := cc.GetCertHash()
	if err != nil {
		return nil, fmt.Errorf("get cert hash in hex failed, %s", err.Error())
	}

	payloadBytes, err := constructSystemContractPayload("", pb.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(),
		pb.CertManageFunction_CERT_ADD.String(), []*pb.KeyValuePair{}, defaultSequence)
	if err != nil {
		return nil, fmt.Errorf("construct transact payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return resp, fmt.Errorf("%s failed, %s", pb.TxType_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_SYSTEM_CONTRACT.String(), err.Error())
	}

	resp.ContractResult = &pb.ContractResult{
		Code:    pb.ContractResultCode_OK,
		Message: pb.ContractResultCode_OK.String(),
		Result:  certHash,
	}

	return resp, nil
}

func (cc ChainClient) DeleteCert(certHashes []string) (*pb.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to INVOKE system contract, [contract:%s]/[method:%s]",
		pb.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(), pb.CertManageFunction_CERTS_DELETE.String())

	payloadBytes, err := constructSystemContractPayload(
		"",
		pb.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(),
		pb.CertManageFunction_CERTS_DELETE.String(),
		[]*pb.KeyValuePair{
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

	resp, err := cc.proposalRequest(pb.TxType_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return resp, fmt.Errorf("%s failed, %s", pb.TxType_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, false); err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_SYSTEM_CONTRACT.String(), err.Error())
	}

	resp.ContractResult = &pb.ContractResult{
		Code:    pb.ContractResultCode_OK,
		Message: pb.ContractResultCode_OK.String(),
	}

	return resp, nil
}

func (cc ChainClient) QueryCert(certHashes []string) (*pb.CertInfos, error) {
	cc.logger.Infof("[SDK] begin to INVOKE system contract, [contract:%s]/[method:%s]",
		pb.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(), pb.CertManageFunction_CERTS_QUERY.String())

	payloadBytes, err := constructQueryPayload(
		pb.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(),
		pb.CertManageFunction_CERTS_QUERY.String(),
		[]*pb.KeyValuePair{
			{
				Key:   "cert_hashes",
				Value: strings.Join(certHashes, ","),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("marshal query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_SYSTEM_CONTRACT, GetRandTxId(), payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_QUERY_SYSTEM_CONTRACT.String(), err.Error())
	}

	certInfos := &pb.CertInfos{}
	if err = proto.Unmarshal(resp.ContractResult.Result, certInfos); err != nil {
		return nil, fmt.Errorf("unmarshal cert infos payload failed, %s", err.Error())
	}

	return certInfos, nil
}

func (cc ChainClient) GetCertHash() ([]byte, error) {
	chainConfig, err := cc.GetChainConfig()

	if err != nil {
		return nil, fmt.Errorf("get cert hash failed, %s", err.Error())
	}

	member := &pb.SerializedMember{
		OrgId:      cc.orgId,
		MemberInfo: cc.userCrtPEM,
		IsFullCert: true,
	}

	certHash, err := getCertificateId(member.GetMemberInfo(), chainConfig.Crypto.Hash)
	if err != nil {
		return nil, fmt.Errorf("calc cert hash failed, %s", err.Error())
	}

	return certHash, nil
}

func (cc ChainClient) CreateCertManagePayload(method string, kvs []*pb.KeyValuePair) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [CertManage] to be signed payload")

	payload := &pb.SystemContractPayload{
		ChainId:      cc.chainId,
		ContractName: pb.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(),
		Method:       method,
		Parameters:   kvs,
	}

	bytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("construct cert manage payload failed, %s", err)
	}

	return bytes, nil
}

func (cc ChainClient) CreateCertManageFrozenPayload(certs []string) ([]byte, error) {
	pairs := []*pb.KeyValuePair{
		{
			Key:   "certs",
			Value: strings.Join(certs, ","),
		},
	}

	return cc.CreateCertManagePayload(pb.CertManageFunction_CERTS_FROZEN.String(), pairs)
}

func (cc ChainClient) CreateCertManageUnfrozenPayload(certs []string) ([]byte, error) {
	pairs := []*pb.KeyValuePair{
		{
			Key:   "certs",
			Value: strings.Join(certs, ","),
		},
	}

	return cc.CreateCertManagePayload(pb.CertManageFunction_CERTS_UNFROZEN.String(), pairs)
}

func (cc ChainClient) CreateCertManageRevocationPayload(certCrl string) ([]byte, error) {
	pairs := []*pb.KeyValuePair{
		{
			Key:   "cert_crl",
			Value: certCrl,
		},
	}

	return cc.CreateCertManagePayload(pb.CertManageFunction_CERTS_REVOCATION.String(), pairs)
}

func (cc ChainClient) SignCertManagePayload(payloadBytes []byte) ([]byte, error) {
	payload := &pb.SystemContractPayload{}
	if err := proto.Unmarshal(payloadBytes, payload); err != nil {
		return nil, fmt.Errorf("unmarshal contract manage payload failed, %s", err)
	}

	signBytes, err := signPayload(cc.privateKey, cc.userCrt, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("SignPayload failed, %s", err)
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

	payload.Endorsement = []*pb.EndorsementEntry{
		entry,
	}

	signedPayloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal contract manage sigend payload failed, %s", err)
	}

	return signedPayloadBytes, nil
}

func (cc ChainClient) MergeCertManageSignedPayload(signedPayloadBytes [][]byte) ([]byte, error) {
	return mergeSystemContractSignedPayload(signedPayloadBytes)
}

func (cc ChainClient) SendCertManageRequest(mergeSignedPayloadBytes []byte, timeout int64, withSyncResult bool) (*pb.TxResponse, error) {
	return cc.sendContractRequest(pb.TxType_SYSTEM_CONTRACT, mergeSignedPayloadBytes, timeout, withSyncResult)
}
