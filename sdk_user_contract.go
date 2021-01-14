/**
 * @Author: jasonruan
 * @Date:   2020-12-02 10:08:52
 **/
package chainmaker_sdk_go

import (
	"fmt"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"io/ioutil"
	"time"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/golang/protobuf/proto"
)

const (
	// 轮训交易结果最大次数
	retryCnt = 5
)

func (cc ChainClient) CreateContractCreatePayload(contractName, version, byteCodePath string, runtime pb.RuntimeType, kvs []*pb.KeyValuePair) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [ContractCreate] to be signed payload")
	return cc.createContractManagePayload(contractName, pb.ManageUserContractFunction_INIT_CONTRACT.String(), version, byteCodePath, runtime, kvs)
}

func (cc ChainClient) CreateContractUpgradePayload(contractName, version, byteCodePath string, runtime pb.RuntimeType, kvs []*pb.KeyValuePair) ([]byte, error) {
	cc.logger.Debugf("[SDK] create [ContractUpgrade] to be signed payload")
	return cc.createContractManagePayload(contractName, pb.ManageUserContractFunction_UPGRADE_CONTRACT.String(), version, byteCodePath, runtime, kvs)
}

func (cc ChainClient) CreateContractFreezePayload(contractName string) ([]byte, error) {
	return cc.createContractOpPayload(contractName, pb.ManageUserContractFunction_FREEZE_CONTRACT.String())
}

func (cc ChainClient) CreateContractUnfreezePayload(contractName string) ([]byte, error) {
	return cc.createContractOpPayload(contractName, pb.ManageUserContractFunction_UNFREEZE_CONTRACT.String())
}

func (cc ChainClient) CreateContractRevokePayload(contractName string) ([]byte, error) {
	return cc.createContractOpPayload(contractName, pb.ManageUserContractFunction_REVOKE_CONTRACT.String())
}

func (cc ChainClient) createContractManagePayload(contractName, method, version, byteCodePath string, runtime pb.RuntimeType, kvs []*pb.KeyValuePair) ([]byte, error) {
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
		Method:     method,
		Parameters: kvs,
		ByteCode:   codeBytes,
	}

	bytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("construct contract manage payload failed, %s", err)
	}

	return bytes, nil
}

func (cc ChainClient) createContractOpPayload(contractName, method string) ([]byte, error) {
	payload := &pb.ContractMgmtPayload{
		ChainId: cc.chainId,
		ContractId: &pb.ContractId{
			ContractName:    contractName,
		},
		Method: method,
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

func (cc ChainClient) SendContractManageRequest(mergeSignedPayloadBytes []byte, timeout int64, withSyncResult bool) (*pb.TxResponse, error) {
	return cc.sendContractRequest(pb.TxType_MANAGE_USER_CONTRACT, mergeSignedPayloadBytes, timeout, withSyncResult)
}

func (cc ChainClient) sendContractRequest(txType pb.TxType, mergeSignedPayloadBytes []byte, timeout int64, withSyncResult bool) (*pb.TxResponse, error) {
	txId := GetRandTxId()

	resp, err := cc.proposalRequestWithTimeout(txType, txId, mergeSignedPayloadBytes, timeout)
	if err != nil {
		return resp, fmt.Errorf("send %s failed, %s", txType.String(), err.Error())
	}

	if resp.Code == pb.TxStatusCode_SUCCESS {
		if !withSyncResult {
			resp.ContractResult = &pb.ContractResult{
				Code:    pb.ContractResultCode_OK,
				Message: pb.ContractResultCode_OK.String(),
				Result:  []byte(txId),
			}
		} else {
			contractResult, err := cc.getSyncResult(txId)
			if err != nil {
				return nil, fmt.Errorf("get sync result failed, %s", err.Error())
			}
			resp.ContractResult = contractResult;
		}
	}

	return resp, nil
}

func (cc ChainClient) InvokeContract(contractName, method, txId string, params map[string]string, timeout int64, withSyncResult bool) (*pb.TxResponse, error) {
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

	resp, err := cc.proposalRequestWithTimeout(pb.TxType_INVOKE_USER_CONTRACT, txId, payloadBytes, timeout)
	if err != nil {
		return resp, fmt.Errorf("%s failed, %s", pb.TxType_INVOKE_USER_CONTRACT.String(), err.Error())
	}

	if resp.Code == pb.TxStatusCode_SUCCESS {
		if !withSyncResult {
			resp.ContractResult = &pb.ContractResult{
				Code:    pb.ContractResultCode_OK,
				Message: pb.ContractResultCode_OK.String(),
				Result:  []byte(txId),
			}
		} else {
			contractResult, err := cc.getSyncResult(txId)
			if err != nil {
				return nil, fmt.Errorf("get sync result failed, %s", err.Error())
			}
			resp.ContractResult = contractResult;
		}
	}

	return resp, nil
}

func (cc ChainClient) QueryContract(contractName, method string, params map[string]string, timeout int64) (*pb.TxResponse, error) {
	txId := GetRandTxId()

	cc.logger.Debugf("[SDK] begin to QUERY contract, [contractName:%s]/[method:%s]/[txId:%s]/[params:%+v]",
		contractName, method, txId, params)

	pairs := paramsMap2KVPairs(params)

	payloadBytes, err := constructQueryPayload(contractName, method, pairs)
	if err != nil {
		return nil, fmt.Errorf("construct query payload failed, %s", err.Error())
	}

	resp, err := cc.proposalRequestWithTimeout(pb.TxType_QUERY_USER_CONTRACT, txId, payloadBytes, timeout)
	if err != nil {
		return nil, fmt.Errorf("send %s failed, %s", pb.TxType_QUERY_USER_CONTRACT.String(), err.Error())
	}

	return resp, nil
}

func (cc ChainClient) getSyncResult(txId string) (*pb.ContractResult, error) {
	var (
		txInfo *pb.TransactionInfo
		err error
	)

	err = retry.Retry(func(attempt uint) error {
		txInfo, err = cc.GetTxByTxId(txId)
		if err != nil {
			return err
		}

		return nil
	},
		strategy.Limit(retryCnt),
		strategy.Backoff(backoff.Fibonacci(retryInterval * time.Millisecond)),
	)

	if err != nil {
		return nil, fmt.Errorf("get tx by txId [%s] failed, %s", txId, err.Error())
	}
	if (txInfo == nil || txInfo.Transaction == nil || txInfo.Transaction.Result == nil) {
		return nil, fmt.Errorf("get result by txId [%s] failed, %+v", txId, txInfo)
	}
	return txInfo.Transaction.Result.ContractResult, nil;
}
