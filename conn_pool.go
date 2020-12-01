/**
 * @Author: jasonruan
 * @Date:   2020-11-30 15:10:17
 **/
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"chainmaker.org/chainmaker-go/common/ca"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"math/rand"
	"time"
)

const (
	retryInterval = 500 // 获取可用客户端连接对象重试时间间隔，单位：ms
	retryLimit = 5      // 获取可用客户端连接对象最大重试次数
)

// 客户端连接结构定义
type networkClient struct {
	rpcNode     pb.RpcNodeClient
	conn        *grpc.ClientConn
	nodeAddr    string
}

// 客户端连接池结构定义
type ConnectionPool struct {
	connections         []*networkClient
	logger              Logger
	useTLS              bool
	caPaths             []string
	userKeyFilePath     string
	userCrtFilePath     string
	tlsHostName         string
}

// 创建连接池
func NewConnPool(config *Config) (*ConnectionPool, error) {
	pool := &ConnectionPool{
		logger: config.logger,
		useTLS: config.useTLS,
		caPaths: config.caPaths,
		userKeyFilePath: config.userKeyFilePath,
		userCrtFilePath: config.userCrtFilePath,
		tlsHostName: config.tlsHostName,
	}

	for nodeAddr, cnt := range config.addrsWithConnCnt {
		cli := &networkClient {
			nodeAddr: nodeAddr,
		}

		for i:=0; i<cnt; i++ {
			pool.connections = append(pool.connections, cli)
		}
	}

	// 打散，用作负载均衡
	pool.connections = shuffle(pool.connections)

	return pool, nil
}

// 初始化GPRC客户端连接
func (pool *ConnectionPool) initGRPCConnect(nodeAddr string) (*grpc.ClientConn, error) {
	if pool.useTLS {
		tlsClient := ca.CAClient{
			ServerName: pool.tlsHostName,
			CaPaths:    pool.caPaths,
			CertFile:   pool.userCrtFilePath,
			KeyFile:    pool.userKeyFilePath,
		}
		c, err := tlsClient.GetCredentialsByCA()
		if err != nil {
			return nil, err
		}

		return grpc.Dial(nodeAddr, grpc.WithTransportCredentials(*c))
	} else {
		return grpc.Dial(nodeAddr, grpc.WithInsecure())
	}
}

// 获取空闲的可用客户端连接对象
func (pool *ConnectionPool) getClient() (*networkClient, error) {
	var nc *networkClient

	if err := retry.Retry(func(attempt uint) error {
		for _, cli := range pool.connections {
			if cli.conn == nil || cli.conn.GetState() == connectivity.Shutdown {

				conn, err := pool.initGRPCConnect(cli.nodeAddr)
				if err != nil {
					pool.logger.Errorf("init grpc connection [nodeAddr:%s] failed, %s", cli.nodeAddr, err.Error())
					continue
				}

				cli.conn = conn
				cli.rpcNode = pb.NewRpcNodeClient(conn)
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

	}, strategy.Wait(retryInterval * time.Millisecond), strategy.Limit(retryLimit)); err != nil {
		return nil, err
	}

	return nc, nil
}

// 关闭连接池
func (pool *ConnectionPool) Close() error {
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

// 数组打散
func shuffle(vals []*networkClient) []*networkClient {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]*networkClient, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}

	return ret
}