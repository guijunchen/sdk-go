/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker/common/crypto"
	bcx509 "chainmaker.org/chainmaker/common/crypto/x509"
	"chainmaker.org/chainmaker/common/evmutils"
	"chainmaker.org/chainmaker/common/serialize"
	"chainmaker.org/chainmaker/pb-go/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/common"
	"context"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"strings"
	"time"
)

const (
	errStringFormat    = "%s failed, %s"
	sdkErrStringFormat = "[SDK] %s"
)

var _ SDKInterface = (*ChainClient)(nil)

type ChainClient struct {
	logger       Logger
	pool         *ConnectionPool
	chainId      string
	orgId        string
	userCrtBytes []byte
	userCrt      *bcx509.Certificate
	privateKey   crypto.PrivateKey
	// 用户压缩证书
	enabledCrtHash bool
	userCrtHash    []byte

	// archive config
	archiveConfig *ArchiveConfig
}

func (cc *ChainClient) CreateArchivePayload(method string, kvs []*common.KeyValuePair) ([]byte, error) {
	panic("implement me")
}

func NewNodeConfig(opts ...NodeOption) *NodeConfig {
	config := &NodeConfig{}
	for _, opt := range opts {
		opt(config)
	}

	return config
}

func NewArchiveConfig(opts ...ArchiveOption) *ArchiveConfig {
	config := &ArchiveConfig{}
	for _, opt := range opts {
		opt(config)
	}

	return config
}

func NewChainClient(opts ...ChainClientOption) (*ChainClient, error) {
	config, err := generateConfig(opts...)
	if err != nil {
		return nil, err
	}

	pool, err := NewConnPool(config)
	if err != nil {
		return nil, err
	}

	return &ChainClient{
		pool:          pool,
		logger:        config.logger,
		chainId:       config.chainId,
		orgId:         config.orgId,
		userCrtBytes:  config.userCrtBytes,
		userCrt:       config.userCrt,
		privateKey:    config.privateKey,
		archiveConfig: config.archiveConfig,
	}, nil
}

func (cc *ChainClient) Stop() error {
	return cc.pool.Close()
}

func (cc *ChainClient) EnableCertHash() error {
	var (
		err error
	)

	// 0.已经启用压缩证书
	if cc.enabledCrtHash {
		return nil
	}

	// 1.如尚未获取证书Hash，便进行获取
	if len(cc.userCrtHash) == 0 {
		// 获取证书Hash
		cc.userCrtHash, err = cc.GetCertHash()
		if err != nil {
			errMsg := fmt.Sprintf("get cert hash failed, %s", err.Error())
			cc.logger.Errorf(sdkErrStringFormat, errMsg)
			return errors.New(errMsg)
		}
	}

	// 2.链上查询证书是否存在
	ok, err := cc.getCheckCertHash()
	if err != nil {
		errMsg := fmt.Sprintf("enable cert hash, get and check cert hash failed, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	// 3.1 若证书已经上链，直接返回
	if ok {
		cc.enabledCrtHash = true
		return nil
	}

	// 3.2 若证书未上链，添加证书
	resp, err := cc.AddCert()
	if err != nil {
		errMsg := fmt.Sprintf("enable cert hash AddCert failed, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		errMsg := fmt.Sprintf("enable cert hash AddCert got invalid resp, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	// 循环检查证书是否成功上链
	err = cc.checkUserCertOnChain()
	if err != nil {
		errMsg := fmt.Sprintf("check user cert on chain failed, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	cc.enabledCrtHash = true

	return nil
}

func (cc *ChainClient) DisableCertHash() error {
	cc.enabledCrtHash = false
	return nil
}

func (cc *ChainClient) EasyCodecItemToParamsMap(items []*serialize.EasyCodecItem) map[string]string {
	return serialize.EasyCodecItemToParamsMap(items)
}

// 检查证书是否成功上链
func (cc *ChainClient) checkUserCertOnChain() error {
	err := retry.Retry(func(uint) error {
		ok, err := cc.getCheckCertHash()
		if err != nil {
			errMsg := fmt.Sprintf("check user cert on chain, get and check cert hash failed, %s", err.Error())
			cc.logger.Errorf(sdkErrStringFormat, errMsg)
			return errors.New(errMsg)
		}

		if !ok {
			errMsg := fmt.Sprintf("user cert havenot on chain yet, and try again")
			cc.logger.Debugf(sdkErrStringFormat, errMsg)
			return errors.New(errMsg)
		}

		return nil
	}, strategy.Limit(10), strategy.Wait(time.Second))

	if err != nil {
		errMsg := fmt.Sprintf("check user upload cert on chain failed, try again later, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return errors.New(errMsg)
	}

	return nil
}

func (cc *ChainClient) getCheckCertHash() (bool, error) {
	// 根据已缓存证书Hash，查链上是否存在
	certInfo, err := cc.QueryCert([]string{hex.EncodeToString(cc.userCrtHash)})
	if err != nil {
		errMsg := fmt.Sprintf("QueryCert failed, %s", err.Error())
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return false, errors.New(errMsg)
	}

	if len(certInfo.CertInfos) == 0 {
		return false, nil
	}

	// 返回链上证书列表长度不为1，即报错
	if len(certInfo.CertInfos) > 1 {
		errMsg := fmt.Sprintf("CertInfos != 1")
		cc.logger.Errorf(sdkErrStringFormat, errMsg)
		return false, errors.New(errMsg)
	}

	// 如果链上证书Hash不为空
	if len(certInfo.CertInfos[0].Cert) > 0 {
		// 如果和缓存的证书Hash不一致则报错
		if hex.EncodeToString(cc.userCrtHash) != certInfo.CertInfos[0].Hash {
			errMsg := fmt.Sprintf("not equal certHash, [expected:%s]/[actual:%s]",
				cc.userCrtHash, certInfo.CertInfos[0].Hash)
			cc.logger.Errorf(sdkErrStringFormat, errMsg)
			return false, errors.New(errMsg)
		}

		// 如果和缓存的证书Hash一致，则说明已经上传好了证书，具备提交压缩证书交易的能力
		return true, nil
	}

	return false, nil
}

func (cc *ChainClient) generateTxRequest(txId string, txType common.TxType, payloadBytes []byte) (*common.TxRequest, error) {
	var (
		sender *accesscontrol.SerializedMember
	)

	// 构造Sender
	if cc.enabledCrtHash && len(cc.userCrtHash) > 0 {
		sender = &accesscontrol.SerializedMember{
			OrgId:      cc.orgId,
			MemberInfo: cc.userCrtHash,
			IsFullCert: false,
		}
	} else {
		sender = &accesscontrol.SerializedMember{
			OrgId:      cc.orgId,
			MemberInfo: cc.userCrtBytes,
			IsFullCert: true,
		}
	}

	// 构造Header
	header := &common.TxHeader{
		ChainId:        cc.chainId,
		Sender:         sender,
		TxType:         txType,
		TxId:           txId,
		Timestamp:      time.Now().Unix(),
		ExpirationTime: 0,
	}

	req := &common.TxRequest{
		Header:    header,
		Payload:   payloadBytes,
		Signature: nil,
	}

	// 拼接后，计算Hash，对hash计算签名
	rawTxBytes, err := CalcUnsignedTxRequestBytes(req)
	if err != nil {
		return nil, err
	}

	signBytes, err := SignTx(cc.privateKey, cc.userCrt, rawTxBytes)
	if err != nil {
		return nil, fmt.Errorf("SignTx failed, %s", err)
	}

	req.Signature = signBytes

	return req, nil
}

func (cc *ChainClient) proposalRequest(txType common.TxType, txId string, payloadBytes []byte) (*common.TxResponse, error) {
	return cc.proposalRequestWithTimeout(txType, txId, payloadBytes, -1)
}

func (cc *ChainClient) proposalRequestWithTimeout(txType common.TxType, txId string, payloadBytes []byte, timeout int64) (*common.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	req, err := cc.generateTxRequest(txId, txType, payloadBytes)
	if err != nil {
		return nil, err
	}

	return cc.sendTxRequest(req, timeout)
}

func (cc *ChainClient) sendTxRequest(txRequest *common.TxRequest, timeout int64) (*common.TxResponse, error) {

	var (
		errMsg string
	)

	if timeout < 0 {
		timeout = SendTxTimeout
		if strings.HasPrefix(txRequest.Header.TxType.String(), "QUERY") {
			timeout = GetTxTimeout
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	ignoreAddrs := make(map[string]struct{})
	for {
		client, err := cc.pool.getClientWithIgnoreAddrs(ignoreAddrs)
		if err != nil {
			return nil, err
		}

		if len(ignoreAddrs) > 0 {
			cc.logger.Debugf("[SDK] begin try to connect node [%s]", client.ID)
		}

		resp, err := client.rpcNode.SendRequest(ctx, txRequest)
		if err != nil {
			resp := &common.TxResponse{
				Message: err.Error(),
				ContractResult: &common.ContractResult{
					Code:    common.ContractResultCode_FAIL,
					Result:  []byte(txRequest.Header.TxId),
					Message: common.ContractResultCode_FAIL.String(),
				},
			}

			statusErr, ok := status.FromError(err)
			if ok && (statusErr.Code() == codes.DeadlineExceeded ||
				// desc = "transport: Error while dialing dial tcp 127.0.0.1:12301: connect: connection refused"
				statusErr.Code() == codes.Unavailable) {

				resp.Code = common.TxStatusCode_TIMEOUT
				errMsg = fmt.Sprintf("call [%s] meet network error, try to connect another node if has, %s",
					client.ID, err.Error())

				cc.logger.Errorf(sdkErrStringFormat, errMsg)
				ignoreAddrs[client.ID] = struct{}{}
				continue
			}

			cc.logger.Errorf("statusErr.Code() : %s", statusErr.Code())

			resp.Code = common.TxStatusCode_INTERNAL_ERROR
			errMsg = fmt.Sprintf("client.call failed, %+v", err)
			cc.logger.Errorf(sdkErrStringFormat, errMsg)
			return resp, fmt.Errorf(errMsg)
		}

		cc.logger.Debugf("[SDK] proposalRequest resp: %+v", resp)
		return resp, nil
	}
}

func (cc *ChainClient) GetEVMAddressFromCertPath(certFilePath string) (string, error) {
	certBytes, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		return "", fmt.Errorf("read cert file [%s] failed, %s", certFilePath, err)
	}

	return cc.GetEVMAddressFromCertBytes(certBytes)
}

func (cc *ChainClient) GetEVMAddressFromCertBytes(certBytes []byte) (string, error) {
	block, _ := pem.Decode(certBytes)
	cert, err := bcx509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("ParseCertificate cert failed, %s", err)
	}

	ski := hex.EncodeToString(cert.SubjectKeyId)
	addrInt, err := evmutils.MakeAddressFromHex(ski)
	if err != nil {
		return "", fmt.Errorf("make address from cert SKI failed, %s", err)
	}

	//return fmt.Sprintf("0x%x", addrInt.AsStringKey()), nil
	//address := evmutils.BigToAddress(addrInt)
	//address := evmutils.EVMIntToHashBytes(addrInt)
	//return hex.EncodeToString([]byte(address)), nil
	//return fmt.Sprintf("%s", address), nil

	return addrInt.String(), nil
}
