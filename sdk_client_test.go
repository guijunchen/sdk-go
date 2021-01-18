/**
 * @Author: jasonruan
 * @Date:   2020-12-01 14:49:44
 */
package chainmaker_sdk_go

import (
	"fmt"
	"log"
)

const (
	chainId        = "chain1"
	orgId1         = "wx-org1.chainmaker.org"
	orgId2         = "wx-org2.chainmaker.org"
	orgId3         = "wx-org3.chainmaker.org"
	orgId4         = "wx-org4.chainmaker.org"
	orgId5         = "wx-org5.chainmaker.org"
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
)

var (
	caPaths     = []string{
		certPathPrefix + fmt.Sprintf("/crypto-config/%s/ca", orgId1),
		certPathPrefix + fmt.Sprintf("/crypto-config/%s/ca", orgId2),
		certPathPrefix + fmt.Sprintf("/crypto-config/%s/ca", orgId3),
		certPathPrefix + fmt.Sprintf("/crypto-config/%s/ca", orgId4),
	}

	userKeyPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.key"
	userCrtPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.crt"

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
		log.Fatal(err)
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
		log.Fatal(err)
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
		log.Fatal(err)
	}

	return adminClient, nil
}
