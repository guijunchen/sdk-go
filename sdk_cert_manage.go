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

func (cc ChainClient) AddCert() (*pb.TxResponse, error) {
	cc.logger.Infof("[SDK] begin to INVOKE system contract, [contract:%s]/[method:%s]",
		pb.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(), pb.CertManageFunction_CERT_ADD.String())

	certHash, err := cc.GetCertHash()
	if err != nil {
		return nil, fmt.Errorf("get cert hash in hex failed, %s", err.Error())
	}

	payloadBytes, err := constructSystemContractPayload("", pb.ContractName_SYSTEM_CONTRACT_CERT_MANAGE.String(),
		pb.CertManageFunction_CERT_ADD.String(), []*pb.KeyValuePair{}, 0)
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
		0,
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
