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
	"fmt"
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
}

func New(opts ...Option) (*ChainClient, error) {
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

func (cc ChainClient) GetChainConfigBeforeBlockHeight(blockHeight int) (*pb.ChainConfig, error) {
	panic("implement me")
}

func (cc ChainClient) SendTransaction(tx *pb.Transaction) (*pb.TxResponse, error) {
	panic("implement me")
}


func (cc ChainClient) Stop() error {
	return cc.pool.Close()
}

func (cc ChainClient) proposalRequest(txType pb.TxType, txId string, payloadBytes []byte) (*pb.TxResponse, error) {
	if txId == "" {
		txId = GetRandTxId()
	}

	timeout := SendTxTimeout
	if strings.HasPrefix(txType.String(), "QUERY") {
		timeout = GetTxTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := cc.pool.getClient()
	if err != nil {
		return nil, err
	}

	// 构造Sender
	sender := &pb.SerializedMember{
		OrgId:      cc.orgId,
		MemberInfo: cc.userCrtPEM,
		IsFullCert: true,
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

	resp, err := client.rpcNode.SendRequest(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				return nil, fmt.Errorf("client.call failed, deadline: %+v", err)
			}
		}

		return nil, fmt.Errorf("client.call failed, %+v", err)
	}

	cc.logger.Debugf("[SDK] proposalRequest resp: %+v", resp)

	return resp, nil
}

