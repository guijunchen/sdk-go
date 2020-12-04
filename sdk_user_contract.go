/**
 * @Author: jasonruan
 * @Date:   2020-12-02 10:08:52
 **/
package chainmaker_sdk_go

import (
	"fmt"
	"io/ioutil"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/golang/protobuf/proto"
)

type ContractManageType int

const (
	TYPE_CREATE  ContractManageType = 0
	TYPE_UPGRADE ContractManageType = 1
)

var (
	mamageType = map[ContractManageType]string{
		TYPE_CREATE:  "init",
		TYPE_UPGRADE: "upgrade",
	}
)

func (cc ChainClient) CreateContractManagePayload(manageType ContractManageType, contractName, version, byteCodePath string, runtime pb.RuntimeType, kvs []*pb.KeyValuePair) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [ContractManage] to be signed payload")

	codeBytes, err := ioutil.ReadFile(byteCodePath)
	if err != nil {
		return nil, fmt.Errorf("Read from file %s error: %s", byteCodePath, err)
	}

	payload := &pb.ContractMgmtPayload{
		ChainId: cc.chainId,
		ContractId: &pb.ContractId{
			ContractName:    contractName,
			ContractVersion: version,
			RuntimeType:     runtime,
		},
		Method:     mamageType[manageType],
		Parameters: kvs,
		ByteCode:   codeBytes,
	}

	bytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("construct contract manage payload failed, %s", err)
	}

	return bytes, nil
}

func (cc ChainClient) SignContractManagePayload(payloadBytes []byte) ([]byte, error) {
	payload := &pb.ContractMgmtPayload{}
	if err := proto.Unmarshal(payloadBytes, payload); err != nil {
		return nil, fmt.Errorf("unmarshal contract manage payload failed, %s", err)
	}

	signBytes, err := signPayload(cc.privateKey, cc.userCrt, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("SignPayload failed, %s", err)
	}

	// TODO: 后续支持证书索引，减小交易大小
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

func (cc ChainClient) MergeContractManageSignedPayload(signedPayloadBytes [][]byte) ([]byte, error) {
	return mergeContractManageSignedPayload(signedPayloadBytes)
}

func (cc ChainClient) SendContractManageRequest(manageType ContractManageType, mergeSignedPayloadBytes []byte) (*pb.TxResponse, error) {
	txId := GetRandTxId()

	txType := pb.TxType_CREATE_USER_CONTRACT
	if manageType == TYPE_UPGRADE {
		txType = pb.TxType_UPGRADE_USER_CONTRACT
	}
	resp, err := cc.proposalRequest(txType, txId, mergeSignedPayloadBytes)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", txType.String(), err.Error())
	}

	resp.ContractResult = &pb.ContractResult{
		Code:    pb.ContractResultCode_OK,
		Message: pb.ContractResultCode_OK.String(),
		Result:  []byte(txId),
	}

	return resp, nil
}

func (cc ChainClient) InvokeContract(contractName, method, txId string, params map[string]string) (*pb.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	cc.logger.Debugf("[SDK] begin to INVOKE contract, [contractName:%s]/[method:%s]/[txId:%s]/[params:%+v]",
		contractName, method, txId, params)

	pairs := paramsMap2KVPairs(params)

	payloadBytes, err := constructTransactPayload(contractName, method, pairs)
	if err != nil {
		return nil, fmt.Errorf("construct transact payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_INVOKE_USER_CONTRACT, txId, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("%s failed, %s", pb.TxType_INVOKE_USER_CONTRACT.String(), err.Error())
	}

	resp.ContractResult = &pb.ContractResult{
		Code:    pb.ContractResultCode_OK,
		Message: pb.ContractResultCode_OK.String(),
		Result:  []byte(txId),
	}

	return resp, nil
}

func (cc ChainClient) QueryContract(contractName, method string, params map[string]string) (*pb.TxResponse, error) {
	txId := GetRandTxId()

	cc.logger.Debugf("[SDK] begin to QUERY contract, [contractName:%s]/[method:%s]/[txId:%s]/[params:%+v]",
		contractName, method, txId, params)

	pairs := paramsMap2KVPairs(params)

	payloadBytes, err := constructQueryPayload(contractName, method, pairs)
	if err != nil {
		return nil, fmt.Errorf("construct query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequest(pb.TxType_QUERY_USER_CONTRACT, txId, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", pb.TxType_QUERY_USER_CONTRACT.String(), err.Error())
	}

	return resp, nil
}
