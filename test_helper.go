/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/test/bufconn"

	apipb "chainmaker.org/chainmaker/pb-go/api"
	cmnpb "chainmaker.org/chainmaker/pb-go/common"
	confpb "chainmaker.org/chainmaker/pb-go/config"
)

const (
	sdkConfigForUtPath = "./testdata/sdk_config_for_ut.yml"
)

var _ ConnectionPool = (*mockConnectionPool)(nil)

type mockConnectionPool struct {
	connections                    []*networkClient
	logger                         Logger
	userKeyBytes                   []byte
	userCrtBytes                   []byte
	rpcClientMaxReceiveMessageSize int
}

func newMockChainClient(opts ...ChainClientOption) (*ChainClient, error) {
	conf, err := generateConfig(opts...)
	if err != nil {
		return nil, err
	}

	pool, err := newMockConnPool(conf)
	if err != nil {
		return nil, err
	}

	return &ChainClient{
		pool:            pool,
		logger:          conf.logger,
		chainId:         conf.chainId,
		orgId:           conf.orgId,
		userCrtBytes:    conf.userSignCrtBytes,
		userCrt:         conf.userCrt,
		privateKey:      conf.privateKey,
		archiveConfig:   conf.archiveConfig,
		rpcClientConfig: conf.rpcClientConfig,
	}, nil
}

func newMockConnPool(config *ChainClientConfig) (*mockConnectionPool, error) {
	pool := &mockConnectionPool{
		logger:                         config.logger,
		userKeyBytes:                   config.userKeyBytes,
		userCrtBytes:                   config.userCrtBytes,
		rpcClientMaxReceiveMessageSize: config.rpcClientConfig.rpcClientMaxReceiveMessageSize,
	}

	for idx, node := range config.nodeList {
		for i := 0; i < node.connCnt; i++ {
			cli := &networkClient{
				nodeAddr:    node.addr,
				useTLS:      node.useTLS,
				caPaths:     node.caPaths,
				caCerts:     node.caCerts,
				tlsHostName: node.tlsHostName,
				ID:          fmt.Sprintf("%v-%v-%v", idx, node.addr, node.tlsHostName),
			}
			pool.connections = append(pool.connections, cli)
		}
	}

	// 打散，用作负载均衡
	pool.connections = shuffle(pool.connections)

	return pool, nil
}

// TODO add tls support
func (pool *mockConnectionPool) initGRPCConnect(nodeAddr string, useTLS bool, caPaths, caCerts []string, tlsHostName string) (*grpc.ClientConn, error) {
	maxCallRecvMsgSize := pool.rpcClientMaxReceiveMessageSize * 1024 * 1024
	return grpc.Dial("", grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxCallRecvMsgSize)), grpc.WithContextDialer(dialer()))
}

// 获取空闲的可用客户端连接对象
func (pool *mockConnectionPool) getClient() (*networkClient, error) {
	return pool.getClientWithIgnoreAddrs(nil)
}

func (pool *mockConnectionPool) getClientWithIgnoreAddrs(ignoreAddrs map[string]struct{}) (*networkClient, error) {
	var nc *networkClient

	err := retry.Retry(func(uint) error {
		for _, cli := range pool.connections {

			if ignoreAddrs != nil {
				if _, ok := ignoreAddrs[cli.ID]; ok {
					continue
				}
			}

			if cli.conn == nil || cli.conn.GetState() == connectivity.Shutdown {

				conn, err := pool.initGRPCConnect(cli.nodeAddr, cli.useTLS, cli.caPaths, cli.caCerts, cli.tlsHostName)
				if err != nil {
					pool.logger.Errorf("init grpc connection [nodeAddr:%s] failed, %s", cli.ID, err.Error())
					continue
				}

				cli.conn = conn
				cli.rpcNode = apipb.NewRpcNodeClient(conn)
				nc = cli
				return nil
			}

			s := cli.conn.GetState()
			if s == connectivity.Idle || s == connectivity.Ready {
				nc = cli
				return nil
			}
		}

		return fmt.Errorf("all client connections are busy")

	}, strategy.Wait(retryInterval*time.Millisecond), strategy.Limit(retryLimit))

	if err != nil {
		return nil, err
	}

	return nc, nil
}

func (pool *mockConnectionPool) getLogger() Logger {
	return pool.logger
}

// Close 关闭连接池
func (pool *mockConnectionPool) Close() error {
	for _, c := range pool.connections {
		if c.conn == nil {
			continue
		}

		if err := c.conn.Close(); err != nil {
			pool.logger.Errorf("stop %s connection failed, %s",
				c.nodeAddr, err.Error())

			continue
		}
	}

	return nil
}

type mockRpcNodeServer struct {
	apipb.UnimplementedRpcNodeServer
}

func (s *mockRpcNodeServer) SendRequest(ctx context.Context, req *cmnpb.TxRequest) (*cmnpb.TxResponse, error) {
	switch req.Header.TxType {
	case cmnpb.TxType_ARCHIVE_FULL_BLOCK:
		return &cmnpb.TxResponse{Code: cmnpb.TxStatusCode_SUCCESS}, nil
	case cmnpb.TxType_RESTORE_FULL_BLOCK:
		return &cmnpb.TxResponse{Code: cmnpb.TxStatusCode_SUCCESS}, nil
	}
	return &cmnpb.TxResponse{}, nil
}

func (s *mockRpcNodeServer) Subscribe(req *cmnpb.TxRequest, server apipb.RpcNode_SubscribeServer) error {
	//var (
	//	errCode cmnerr.ErrCode
	//	errMsg  string
	//)

	//tx := &cmnpb.Transaction{
	//	Header:           req.Header,
	//	RequestPayload:   req.Payload,
	//	RequestSignature: req.Signature,
	//	Result:           nil,
	//}

	//errCode, errMsg = s.validate(tx)
	//if errCode != cmnerr.ERR_CODE_OK {
	//	return status.Error(codes.Unauthenticated, errMsg)
	//}

	switch req.Header.TxType {
	case cmnpb.TxType_SUBSCRIBE_BLOCK_INFO:
		//return s.dealBlockSubscription(tx, server)
	case cmnpb.TxType_SUBSCRIBE_TX_INFO:
		//return s.dealTxSubscription(tx, server)
	case cmnpb.TxType_SUBSCRIBE_CONTRACT_EVENT_INFO:
		//return s.dealContractEventSubscription(tx, server)
	}

	return nil
}

func (s *mockRpcNodeServer) GetChainMakerVersion(ctx context.Context,
	req *confpb.ChainMakerVersionRequest) (*confpb.ChainMakerVersionResponse, error) {
	return &confpb.ChainMakerVersionResponse{
		Code:    0,
		Message: "OK",
		Version: "1.0.0",
	}, nil
}

func (s *mockRpcNodeServer) CheckNewBlockChainConfig(ctx context.Context,
	req *confpb.CheckNewBlockChainConfigRequest) (*confpb.CheckNewBlockChainConfigResponse, error) {
	return &confpb.CheckNewBlockChainConfigResponse{
		Code: 0,
	}, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	apipb.RegisterRpcNodeServer(server, &mockRpcNodeServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
