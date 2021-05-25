/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

const (
	chainId        = "chain1"
	orgId1         = "wx-org1.chainmaker.org"
	orgId2         = "wx-org2.chainmaker.org"
	orgId3         = "wx-org3.chainmaker.org"
	orgId4         = "wx-org4.chainmaker.org"
	orgId5         = "wx-org5.chainmaker.org"
	orgId6         = "wx-org6.chainmaker.org"
	contractName   = "counter-go-1"
	certPathPrefix = "./testdata"
	tlsHostName    = "chainmaker.org"
	version        = "1.0.0"
	upgradeVersion = "2.0.0"

	nodeAddr1 = "127.0.0.1:12301"
	connCnt1  = 5

	nodeAddr2 = "127.0.0.1:12301"
	connCnt2  = 5

	multiSignedPayloadFile        = "./testdata/counter-go-demo/collect-signed-all.pb"
	upgradeMultiSignedPayloadFile = "./testdata/counter-go-demo/upgrade-collect-signed-all.pb"

	byteCodePath        = "./testdata/counter-go-demo/counter-rust-0.7.2.wasm"
	upgradeByteCodePath = "./testdata/counter-go-demo/counter-go-upgrade.wasm"

	certPathFormat = "/crypto-config/%s/ca"
)

var (
	caPaths = []string{
		certPathPrefix + fmt.Sprintf(certPathFormat, orgId1),
		certPathPrefix + fmt.Sprintf(certPathFormat, orgId2),
		certPathPrefix + fmt.Sprintf(certPathFormat, orgId3),
		certPathPrefix + fmt.Sprintf(certPathFormat, orgId4),
	}

	caCerts = []string{"-----BEGIN CERTIFICATE-----\nMIICsDCCAlWgAwIBAgIDAuGKMAoGCCqBHM9VAYN1MIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTIxMDMyNTA2NDI1MVoXDTMx\nMDMyMzA2NDI1MVowgYoxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn\nMRIwEAYDVQQLEwlyb290LWNlcnQxIjAgBgNVBAMTGWNhLnd4LW9yZzEuY2hhaW5t\nYWtlci5vcmcwWTATBgcqhkjOPQIBBggqgRzPVQGCLQNCAARIG6tdLNtG+eqwTK36\nS/AjzXh9Q0Zwrf7eqyCEQ4Ul7xfgKjCBNVboivH10ieYuh0MAoZj1Ke7z+P6ZUTy\naiuDo4GnMIGkMA4GA1UdDwEB/wQEAwIBpjAPBgNVHSUECDAGBgRVHSUAMA8GA1Ud\nEwEB/wQFMAMBAf8wKQYDVR0OBCIEIJDsy2L0fAK2V4YxOjVEjYj3YKSbX4F24eh0\nZQHoqCr1MEUGA1UdEQQ+MDyCDmNoYWlubWFrZXIub3Jngglsb2NhbGhvc3SCGWNh\nLnd4LW9yZzEuY2hhaW5tYWtlci5vcmeHBH8AAAEwCgYIKoEcz1UBg3UDSQAwRgIh\nAM1oJOU6l4tJVqrCJv5UnMaKLxu4V1dDwu0YsS5Tb1s9AiEA1D8NA3GGy9BEFryq\n5TS0uiqE3QEuDRvs1TrP9H53Sjk=\n-----END CERTIFICATE-----",}

	userKeyPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.key"
	userCrtPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.crt"

	userSignKeyPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.sign.key"
	userSignCrtPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.sign.crt"

	adminKeyPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.key"
	adminCrtPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.crt"
)

// 中间证书测试
//var (
//	chainOrgId  = orgId4
//	caPaths     = []string{certPathPrefix + fmt.Sprintf("/crypto-config-middle-cert/%s/ca", chainOrgId)}
//	userKeyPath = certPathPrefix + fmt.Sprintf("/crypto-config-middle-cert/%s/user/client1/client1.tls.key", chainOrgId)
//	userCrtPath = certPathPrefix + fmt.Sprintf("/crypto-config-middle-cert/%s/user/client1/client1.tls.crt", chainOrgId)
//
//	adminKeyPath = certPathPrefix + "/crypto-config-middle-cert/%s/user/admin1/admin1.tls.key"
//	adminCrtPath = certPathPrefix + "/crypto-config-middle-cert/%s/user/admin1/admin1.tls.crt"
//)

var (
	node1 *NodeConfig
	node2 *NodeConfig
)

// 创建节点
func createNode(nodeAddr string, connCnt int) *NodeConfig {
	node := NewNodeConfig(
		// 节点地址，格式：127.0.0.1:12301
		WithNodeAddr(nodeAddr),
		// 节点连接数
		WithNodeConnCnt(connCnt),
		// 节点是否启用TLS认证
		WithNodeUseTLS(true),
		// 根证书路径，支持多个
		WithNodeCAPaths(caPaths),
		// TLS Hostname
		WithNodeTLSHostName(tlsHostName),
	)

	return node
}

func createClient() (*ChainClient, error) {
	return createClientWithOrgId(orgId1)
}

// 创建ChainClient
func createClientWithOrgId(orgId string) (*ChainClient, error) {
	if node1 == nil {
		// 创建节点1
		node1 = createNode(nodeAddr1, connCnt1)
	}

	if node2 == nil {
		// 创建节点2
		node2 = createNode(nodeAddr2, connCnt2)
	}

	chainClient, err := NewChainClient(
		// 设置归属组织
		WithChainClientOrgId(orgId),
		// 设置链ID
		WithChainClientChainId(chainId),
		// 设置logger句柄，若不设置，将采用默认日志文件输出日志
		WithChainClientLogger(getDefaultLogger()),
		// 设置客户端用户私钥路径
		WithUserKeyFilePath(fmt.Sprintf(userKeyPath, orgId)),
		// 设置客户端用户证书
		WithUserCrtFilePath(fmt.Sprintf(userCrtPath, orgId)),
		// 添加节点1
		AddChainClientNodeConfig(node1),
		// 添加节点2
		AddChainClientNodeConfig(node2),
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

// 创建ChainClient（使用配置文件）
func createClientWithConfig() (*ChainClient, error) {

	chainClient, err := NewChainClient(
		WithConfPath("./testdata/sdk_config.yml"),
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

// 创建ChainClient（指定证书内容）
func createClientWithCertBytes() (*ChainClient, error) {

	userCrtBytes, err := ioutil.ReadFile(fmt.Sprintf(userCrtPath, orgId1))
	if err != nil {
		return nil, err
	}

	userKeyBytes, err := ioutil.ReadFile(fmt.Sprintf(userKeyPath, orgId1))
	if err != nil {
		return nil, err
	}

	userSignCrtBytes, err := ioutil.ReadFile(fmt.Sprintf(userSignCrtPath, orgId1))
	if err != nil {
		return nil, err
	}

	userSignKeyBytes, err := ioutil.ReadFile(fmt.Sprintf(userSignKeyPath, orgId1))
	if err != nil {
		return nil, err
	}

	chainClient, err := NewChainClient(
		WithConfPath("./testdata/sdk_config.yml"),
		WithUserCrtBytes(userCrtBytes),
		WithUserKeyBytes(userKeyBytes),
		WithUserSignKeyBytes(userSignKeyBytes),
		WithUserSignCrtBytes(userSignCrtBytes),
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

func createNodeWithCaCert(nodeAddr string, connCnt int) *NodeConfig {
	node := NewNodeConfig(
		// 节点地址，格式：127.0.0.1:12301
		WithNodeAddr(nodeAddr),
		// 节点连接数
		WithNodeConnCnt(connCnt),
		// 节点是否启用TLS认证
		WithNodeUseTLS(true),
		// 根证书内容，支持多个
		WithNodeCACerts(caCerts),
		// TLS Hostname
		WithNodeTLSHostName(tlsHostName),
	)

	return node
}

func createClientWithCaCerts() (*ChainClient, error) {
	if node1 == nil {
		// 创建节点1
		node1 = createNodeWithCaCert(nodeAddr1, connCnt1)
	}

	chainClient, err := NewChainClient(
		// 设置归属组织
		WithChainClientOrgId(orgId1),
		// 设置链ID
		WithChainClientChainId(chainId),
		// 设置logger句柄，若不设置，将采用默认日志文件输出日志
		WithChainClientLogger(getDefaultLogger()),
		// 设置客户端用户私钥路径
		WithUserKeyFilePath(fmt.Sprintf(userKeyPath, orgId1)),
		// 设置客户端用户证书
		WithUserCrtFilePath(fmt.Sprintf(userCrtPath, orgId1)),
		// 添加节点1
		AddChainClientNodeConfig(node1),
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

func createAdmin(orgId string) (*ChainClient, error) {
	if node1 == nil {
		node1 = createNode(nodeAddr1, connCnt1)
	}

	if node2 == nil {
		node2 = createNode(nodeAddr2, connCnt2)
	}

	adminClient, err := NewChainClient(
		WithChainClientOrgId(orgId),
		WithChainClientChainId(chainId),
		WithChainClientLogger(getDefaultLogger()),
		WithUserKeyFilePath(fmt.Sprintf(adminKeyPath, orgId)),
		WithUserCrtFilePath(fmt.Sprintf(adminCrtPath, orgId)),
		AddChainClientNodeConfig(node1),
		AddChainClientNodeConfig(node2),
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

func TestChainClient_GetEVMAddressFromCertPath(t *testing.T) {
	client, err := createClientWithConfig()
	require.Nil(t, err)

	certFilePath := fmt.Sprintf(userCrtPath, orgId1)
	addr, err := client.GetEVMAddressFromCertPath(certFilePath)
	require.Nil(t, err)
	fmt.Printf("client1 address: %s\n", addr)

	certFilePath = fmt.Sprintf(userCrtPath, orgId2)
	addr, err = client.GetEVMAddressFromCertPath(certFilePath)
	require.Nil(t, err)
	fmt.Printf("client2 address: %s\n", addr)
}