/**
 * @Author: jasonruan
 * @Date:   2020-11-30 14:44:30
 */
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"chainmaker.org/chainmaker-go/common/crypto"
	bcx509 "chainmaker.org/chainmaker-go/common/crypto/x509"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

var _ SDKInterface = (*ChainClient)(nil)

type ChainClient struct {
	logger      Logger
	pool        *ConnectionPool
	chainId     string
	orgId       string
	userCrtPEM  []byte
	userCrt     *bcx509.Certificate
	privateKey  crypto.PrivateKey
	// 用户压缩证书
	enabledCrtHash bool
	userCrtHash []byte
}

func NewNodeConfig(opts ...NodeOption) *NodeConfig {
	config := &NodeConfig{}
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
		pool:           pool,
		logger:         config.logger,
		chainId:        config.chainId,
		orgId:          config.orgId,
		userCrtPEM:     config.userCrtPEM,
		userCrt:        config.userCrt,
		privateKey:     config.privateKey,
	}, nil
}

func (cc ChainClient) Stop() error {
	return cc.pool.Close()
}

func (cc *ChainClient) EnableCertHash() error {
	if len(cc.userCrtHash) > 0 {
		ok, err := cc.getCheckCertHash(cc.userCrtHash)
		if err != nil {
			errMsg := fmt.Sprintf("enable cert hash, get and check cert hash failed, %s", err.Error())
			cc.logger.Errorf("[SDK] %s", errMsg)
			return errors.New(errMsg)
		}

		if !ok {
			// 如果链上证书Hash为空，说明还没有在链上添加证书，执行后续方法添加之
			cc.logger.Warnf("[SDK] %s", "havenot get user cert on chain [%s], begin add cert", cc.chainId)
		}
	}

	resp, err := cc.AddCert()
	if err != nil {
		errMsg := fmt.Sprintf("enable cert hash AddCert failed, %s", err.Error())
		cc.logger.Errorf("[SDK] %s", errMsg)
		return errors.New(errMsg)
	}

	if err = checkProposalRequestResp(resp, true); err != nil {
		errMsg := fmt.Sprintf("enable cert hash AddCert got invalid resp, %s", err.Error())
		cc.logger.Errorf("[SDK] %s", errMsg)
		return errors.New(errMsg)
	}

	cc.userCrtHash = resp.ContractResult.Result
	cc.enabledCrtHash = true

	err = cc.checkUserCertOnChain(cc.userCrtHash)
	if err != nil {
		errMsg := fmt.Sprintf("check user cert on chain failed, %s", err.Error())
		cc.logger.Errorf("[SDK] %s", errMsg)
		return errors.New(errMsg)
	}

	return nil
}

// 检查证书是否成功上链
func (cc ChainClient) checkUserCertOnChain(userCrtHash []byte) error {
	if err := retry.Retry(func(attempt uint) error {
		ok, err := cc.getCheckCertHash(cc.userCrtHash)
		if err != nil {
			errMsg := fmt.Sprintf("check user cert on chain, get and check cert hash failed, %s", err.Error())
			cc.logger.Errorf("[SDK] %s", errMsg)
			return errors.New(errMsg)
		}

		if !ok {
			errMsg := fmt.Sprintf("user cert havenot on chain yet, and try again")
			cc.logger.Debugf("[SDK] %s", errMsg)
			return errors.New(errMsg)
		}

		return nil
	}, strategy.Limit(10), strategy.Wait(time.Second),
	); err != nil {
		errMsg := fmt.Sprintf("check user upload cert on chain failed, try again later, %s", err.Error())
		cc.logger.Errorf("[SDK] %s", errMsg)
		return errors.New(errMsg)
	}

	return nil
}

func (cc ChainClient) getCheckCertHash(userCrtHash []byte) (bool, error) {
	// 根据已缓存证书Hash，查链上是否存在
	certInfo, err := cc.QueryCert([]string{hex.EncodeToString(userCrtHash)})
	if err != nil {
		errMsg := fmt.Sprintf("QueryCert failed, %s", err.Error())
		cc.logger.Errorf("[SDK] %s", errMsg)
		return false, errors.New(errMsg)
	}

	if len(certInfo.CertInfos) == 0 {
		return false, nil
	}

	// 返回链上证书列表长度不为1，即报错
	if len(certInfo.CertInfos) > 1 {
		errMsg := fmt.Sprintf("CertInfos != 1")
		cc.logger.Errorf("[SDK] %s", errMsg)
		return false, errors.New(errMsg)
	}

	// 如果链上证书Hash不为空
	if certInfo.CertInfos[0].Hash != "" {
		// 如果和缓存的证书Hash不一致则报错
		if hex.EncodeToString(cc.userCrtHash) != certInfo.CertInfos[0].Hash {
			errMsg := fmt.Sprintf("not equal certHash, [expected:%s]/[actual:%s]",
				cc.userCrtHash, certInfo.CertInfos[0].Hash)
			cc.logger.Errorf("[SDK] %s", errMsg)
			return false, errors.New(errMsg)
		}

		// 如果和缓存的证书Hash一致，则说明已经上传好了证书，具备提交压缩证书交易的能力
		return true, nil
	}

	return false, nil
}

func (cc ChainClient) DisableCertHash() error {
	cc.enabledCrtHash = false
	return nil
}

func (cc ChainClient) generateTxRequest(txId string, txType pb.TxType, payloadBytes []byte) (*pb.TxRequest, error) {
	var (
		sender *pb.SerializedMember
	)

	// 构造Sender
	if cc.enabledCrtHash && len(cc.userCrtHash) > 0 {
		sender = &pb.SerializedMember{
			OrgId:      cc.orgId,
			MemberInfo: cc.userCrtHash,
			IsFullCert: false,
		}
	} else {
		sender = &pb.SerializedMember{
			OrgId:      cc.orgId,
			MemberInfo: cc.userCrtPEM,
			IsFullCert: true,
		}
	}

	// 构造Header
	header := &pb.TxHeader{
		ChainId:        cc.chainId,
		Sender:         sender,
		TxType:         txType,
		TxId:           txId,
		Timestamp:      time.Now().Unix(),
		ExpirationTime: 0,
	}

	req := &pb.TxRequest{
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

func (cc ChainClient) proposalRequest(txType pb.TxType, txId string, payloadBytes []byte) (*pb.TxResponse, error) {
	return cc.proposalRequestWithTimeout(txType, txId, payloadBytes, -1)
}

func (cc ChainClient) proposalRequestWithTimeout(txType pb.TxType, txId string, payloadBytes []byte, timeout int64) (*pb.TxResponse, error) {
	var (
		errMsg string
	)

	if txId == "" {
		txId = GetRandTxId()
	}

	if timeout < 0 {
		timeout = SendTxTimeout
		if strings.HasPrefix(txType.String(), "QUERY") {
			timeout = GetTxTimeout
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout) * time.Second)
	defer cancel()

	client, err := cc.pool.getClient()
	if err != nil {
		return nil, err
	}

	req, err := cc.generateTxRequest(txId, txType, payloadBytes)
	if err != nil {
		return nil, err
	}

	resp, err := client.rpcNode.SendRequest(ctx, req)
	if err != nil {
		resp := &pb.TxResponse{
			Message: err.Error(),
			ContractResult: &pb.ContractResult{
				Code:    pb.ContractResultCode_FAIL,
				Result:  []byte(txId),
				Message: pb.ContractResultCode_FAIL.String(),
			},
		}

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				resp.Code = pb.TxStatusCode_TIMEOUT
				errMsg = fmt.Sprintf("client.call failed, deadline: %+v", err)
				cc.logger.Errorf("[SDK] %s", errMsg)
				return resp, fmt.Errorf(errMsg)
			}
		}

		resp.Code = pb.TxStatusCode_INTERNAL_ERROR
		errMsg = fmt.Sprintf("client.call failed, %+v", err)
		cc.logger.Errorf("[SDK] %s", errMsg)
		return resp, fmt.Errorf(errMsg)
	}

	cc.logger.Debugf("[SDK] proposalRequest resp: %+v", resp)

	return resp, nil
}

