/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package examples

import (
	"errors"
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"

	"chainmaker.org/chainmaker/common/log"
	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
)

const (
	chainId        = "chain1"
	OrgId1         = "wx-org1.chainmaker.org"
	OrgId2         = "wx-org2.chainmaker.org"
	OrgId3         = "wx-org3.chainmaker.org"
	OrgId4         = "wx-org4.chainmaker.org"
	OrgId5         = "wx-org5.chainmaker.org"
	orgId6         = "wx-org6.chainmaker.org"
	contractName   = "counter-go-1"
	certPathPrefix = "../../testdata"
	tlsHostName    = "chainmaker.org"
	Version        = "1.0.0"
	UpgradeVersion = "2.0.0"

	nodeAddr1 = "127.0.0.1:12301"
	connCnt1  = 5

	nodeAddr2 = "127.0.0.1:12302"
	connCnt2  = 5

	multiSignedPayloadFile        = "../testdata/counter-go-demo/collect-signed-all.pb"
	upgradeMultiSignedPayloadFile = "../testdata/counter-go-demo/upgrade-collect-signed-all.pb"

	byteCodePath        = "../testdata/counter-go-demo/counter-rust-0.7.2.wasm"
	upgradeByteCodePath = "../testdata/counter-go-demo/counter-go-upgrade.wasm"

	certPathFormat = "/crypto-config/%s/ca"
)

var (
	caPaths = []string{
		certPathPrefix + fmt.Sprintf(certPathFormat, OrgId1),
		certPathPrefix + fmt.Sprintf(certPathFormat, OrgId2),
		certPathPrefix + fmt.Sprintf(certPathFormat, OrgId3),
		certPathPrefix + fmt.Sprintf(certPathFormat, OrgId4),
	}

	caCerts = []string{"-----BEGIN CERTIFICATE-----\nMIICsDCCAlWgAwIBAgIDAuGKMAoGCCqBHM9VAYN1MIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTIxMDMyNTA2NDI1MVoXDTMx\nMDMyMzA2NDI1MVowgYoxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMRIwEAYDVQQLEwlyb290LWNlcnQxIjAgBgNVBAMTGWNhLnd4LW9yZzEuY2hhaW5t\nYWtlci5vcmcwWTATBgcqhkjOPQIBBggqgRzPVQGCLQNCAARIG6tdLNtG+eqwTK36\nS/AjzXh9Q0Zwrf7eqyCEQ4Ul7xfgKjCBNVboivH10ieYuh0MAoZj1Ke7z+P6ZUTy\naiuDo4GnMIGkMA4GA1UdDwEB/wQEAwIBpjAPBgNVHSUECDAGBgRVHSUAMA8GA1Ud\nEwEB/wQFMAMBAf8wKQYDVR0OBCIEIJDsy2L0fAK2V4YxOjVEjYj3YKSbX4F24eh0\nZQHoqCr1MEUGA1UdEQQ+MDyCDmNoYWlubWFrZXIub3Jngglsb2NhbGhvc3SCGWNh\nLnd4LW9yZzEuY2hhaW5tYWtlci5vcmeHBH8AAAEwCgYIKoEcz1UBg3UDSQAwRgIh\nAM1oJOU6l4tJVqrCJv5UnMaKLxu4V1dDwu0YsS5Tb1s9AiEA1D8NA3GGy9BEFryq\n5TS0uiqE3QEuDRvs1TrP9H53Sjk=\n-----END CERTIFICATE-----"}

	userKeyPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.key"
	UserCrtPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.crt"

	userSignKeyPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.sign.key"
	userSignCrtPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.sign.crt"

	adminKeyPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.key"
	adminCrtPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.crt"
)

// 中间证书测试
//var (
//	chainOrgId  = OrgId4
//	caPaths     = []string{certPathPrefix + fmt.Sprintf("/crypto-config-middle-cert/%s/ca", chainOrgId)}
//	userKeyPath = certPathPrefix + fmt.Sprintf("/crypto-config-middle-cert/%s/user/client1/client1.tls.key", chainOrgId)
//	UserCrtPath = certPathPrefix + fmt.Sprintf("/crypto-config-middle-cert/%s/user/client1/client1.tls.crt", chainOrgId)
//
//	adminKeyPath = certPathPrefix + "/crypto-config-middle-cert/%s/user/admin1/admin1.tls.key"
//	adminCrtPath = certPathPrefix + "/crypto-config-middle-cert/%s/user/admin1/admin1.tls.crt"
//)

var (
	node1 *sdk.NodeConfig
	node2 *sdk.NodeConfig
)

// CreateNode Create a NodeConfig
func CreateNode(nodeAddr string, connCnt int) *sdk.NodeConfig {
	node := sdk.NewNodeConfig(
		// 节点地址，格式：127.0.0.1:12301
		sdk.WithNodeAddr(nodeAddr),
		// 节点连接数
		sdk.WithNodeConnCnt(connCnt),
		// 节点是否启用TLS认证
		sdk.WithNodeUseTLS(true),
		// 根证书路径，支持多个
		sdk.WithNodeCAPaths(caPaths),
		// TLS Hostname
		sdk.WithNodeTLSHostName(tlsHostName),
	)

	return node
}

func CreateClient() (*sdk.ChainClient, error) {
	return CreateClientWithOrgId(OrgId1)
}

// CreateClientWithOrgId 创建ChainClient
func CreateClientWithOrgId(orgId string) (*sdk.ChainClient, error) {
	if node1 == nil {
		// 创建节点1
		node1 = CreateNode(nodeAddr1, connCnt1)
	}

	if node2 == nil {
		// 创建节点2
		//node2 = CreateNode(nodeAddr2, connCnt2)
	}

	chainClient, err := sdk.NewChainClient(
		// 设置归属组织
		sdk.WithChainClientOrgId(orgId),
		// 设置链ID
		sdk.WithChainClientChainId(chainId),
		// 设置logger句柄，若不设置，将采用默认日志文件输出日志
		sdk.WithChainClientLogger(getDefaultLogger()),
		// 设置客户端用户私钥路径
		sdk.WithUserKeyFilePath(fmt.Sprintf(userKeyPath, orgId)),
		// 设置客户端用户证书
		sdk.WithUserCrtFilePath(fmt.Sprintf(UserCrtPath, orgId)),
		// 添加节点1
		sdk.AddChainClientNodeConfig(node1),
		// 添加节点2
		//AddChainClientNodeConfig(node2),
	)

	if err != nil {
		return nil, err
	}

	//启用证书压缩（开启证书压缩可以减小交易包大小，提升处理性能）
	err = chainClient.EnableCertHash()
	if err != nil {
		return nil, err
	}

	return chainClient, nil
}

// CreateClientWithConfig 创建ChainClient（使用配置文件）
func CreateClientWithConfig() (*sdk.ChainClient, error) {

	chainClient, err := sdk.NewChainClient(
		sdk.WithConfPath("../../testdata/sdk_config.yml"),
	)

	if err != nil {
		return nil, err
	}

	//启用证书压缩（开启证书压缩可以减小交易包大小，提升处理性能）
	err = chainClient.EnableCertHash()
	if err != nil {
		return nil, err
	}

	return chainClient, nil
}

// CreateClientWithCertBytes 创建ChainClient（指定证书内容）
func CreateClientWithCertBytes() (*sdk.ChainClient, error) {

	userCrtBytes, err := ioutil.ReadFile(fmt.Sprintf(UserCrtPath, OrgId1))
	if err != nil {
		return nil, err
	}

	userKeyBytes, err := ioutil.ReadFile(fmt.Sprintf(userKeyPath, OrgId1))
	if err != nil {
		return nil, err
	}

	userSignCrtBytes, err := ioutil.ReadFile(fmt.Sprintf(userSignCrtPath, OrgId1))
	if err != nil {
		return nil, err
	}

	userSignKeyBytes, err := ioutil.ReadFile(fmt.Sprintf(userSignKeyPath, OrgId1))
	if err != nil {
		return nil, err
	}

	chainClient, err := sdk.NewChainClient(
		sdk.WithConfPath("../../testdata/sdk_config.yml"),
		sdk.WithUserCrtBytes(userCrtBytes),
		sdk.WithUserKeyBytes(userKeyBytes),
		sdk.WithUserSignKeyBytes(userSignKeyBytes),
		sdk.WithUserSignCrtBytes(userSignCrtBytes),
	)

	if err != nil {
		return nil, err
	}

	//启用证书压缩（开启证书压缩可以减小交易包大小，提升处理性能）
	err = chainClient.EnableCertHash()
	if err != nil {
		return nil, err
	}

	return chainClient, nil
}

func createNodeWithCaCert(nodeAddr string, connCnt int) *sdk.NodeConfig {
	node := sdk.NewNodeConfig(
		// 节点地址，格式：127.0.0.1:12301
		sdk.WithNodeAddr(nodeAddr),
		// 节点连接数
		sdk.WithNodeConnCnt(connCnt),
		// 节点是否启用TLS认证
		sdk.WithNodeUseTLS(true),
		// 根证书内容，支持多个
		sdk.WithNodeCACerts(caCerts),
		// TLS Hostname
		sdk.WithNodeTLSHostName(tlsHostName),
	)

	return node
}

func createClientWithCaCerts() (*sdk.ChainClient, error) {
	if node1 == nil {
		// 创建节点1
		node1 = createNodeWithCaCert(nodeAddr1, connCnt1)
	}

	chainClient, err := sdk.NewChainClient(
		// 设置归属组织
		sdk.WithChainClientOrgId(OrgId1),
		// 设置链ID
		sdk.WithChainClientChainId(chainId),
		// 设置logger句柄，若不设置，将采用默认日志文件输出日志
		sdk.WithChainClientLogger(getDefaultLogger()),
		// 设置客户端用户私钥路径
		sdk.WithUserKeyFilePath(fmt.Sprintf(userKeyPath, OrgId1)),
		// 设置客户端用户证书
		sdk.WithUserCrtFilePath(fmt.Sprintf(UserCrtPath, OrgId1)),
		// 添加节点1
		sdk.AddChainClientNodeConfig(node1),
	)

	if err != nil {
		return nil, err
	}

	//启用证书压缩（开启证书压缩可以减小交易包大小，提升处理性能）
	err = chainClient.EnableCertHash()
	if err != nil {
		return nil, err
	}

	return chainClient, nil
}

func CreateAdmin(orgId string) (*sdk.ChainClient, error) {
	if node1 == nil {
		node1 = CreateNode(nodeAddr1, connCnt1)
	}

	if node2 == nil {
		node2 = CreateNode(nodeAddr2, connCnt2)
	}

	adminClient, err := sdk.NewChainClient(
		sdk.WithChainClientOrgId(orgId),
		sdk.WithChainClientChainId(chainId),
		sdk.WithChainClientLogger(getDefaultLogger()),
		sdk.WithUserKeyFilePath(fmt.Sprintf(adminKeyPath, orgId)),
		sdk.WithUserCrtFilePath(fmt.Sprintf(adminCrtPath, orgId)),
		sdk.AddChainClientNodeConfig(node1),
		sdk.AddChainClientNodeConfig(node2),
	)
	if err != nil {
		return nil, err
	}

	//启用证书压缩（开启证书压缩可以减小交易包大小，提升处理性能）
	err = adminClient.EnableCertHash()
	if err != nil {
		return nil, err
	}

	return adminClient, nil
}

func getDefaultLogger() *zap.SugaredLogger {
	config := log.LogConfig{
		Module:       "[SDK]",
		LogPath:      "./sdk.log",
		LogLevel:     log.LEVEL_INFO,
		MaxAge:       30,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: true,
	}

	logger, _ := log.InitSugarLogger(&config)
	return logger
}

func CheckProposalRequestResp(resp *common.TxResponse, needContractResult bool) error {
	if resp.Code != common.TxStatusCode_SUCCESS {
		return errors.New(resp.Message)
	}

	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil && resp.ContractResult.Code != common.ContractResultCode_OK {
		return errors.New(resp.ContractResult.Message)
	}

	return nil
}

// CreateChainClientWithSDKConf create a chain client with sdk config file path
func CreateChainClientWithSDKConf(sdkConfPath string) (*sdk.ChainClient, error) {
	logger, _ := log.InitSugarLogger(&log.LogConfig{
		Module:       "[SDK]",
		LogPath:      "./sdk.log",
		LogLevel:     log.LEVEL_ERROR,
		MaxAge:       30,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: true,
	})

	var (
		cc  *sdk.ChainClient
		err error
	)

	cc, err = sdk.NewChainClient(
		sdk.WithConfPath(sdkConfPath),
		sdk.WithChainClientLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	// Enable certificate compression
	err = cc.EnableCertHash()
	if err != nil {
		return nil, err
	}
	return cc, nil
}
