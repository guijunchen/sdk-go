/**
 * @Author: jasonruan
 * @Date:   2020-12-01 14:49:44
 */
package chainmaker_sdk_go

import (
	"fmt"
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

	byteCodePath        = "./testdata/counter-go-demo/counter-go.wasm"
	upgradeByteCodePath = "./testdata/counter-go-demo/counter-go-upgrade.wasm"
)

var (
	caPaths     = []string{certPathPrefix + fmt.Sprintf("/crypto-config/%s/ca", orgId1)}
	userKeyPath = certPathPrefix + fmt.Sprintf("/crypto-config/%s/user/client1/client1.tls.key", orgId1)
	userCrtPath = certPathPrefix + fmt.Sprintf("/crypto-config/%s/user/client1/client1.tls.crt", orgId1)

	adminKeyPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.key"
	adminCrtPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.crt"
)

var (
	node1 *NodeConfig
	node2 *NodeConfig
)

func createNode(nodeAddr string, connCnt int) *NodeConfig {
	node := NewNodeConfig(
		WithNodeAddr(nodeAddr),
		WithNodeConnCnt(connCnt),
		WithNodeUseTLS(true),
		WithNodeCAPaths(caPaths),
		WithNodeTLSHostName(tlsHostName),
	)

	return node
}

func createClient() (*ChainClient, error) {
	if node1 == nil {
		node1 = createNode(nodeAddr1, connCnt1)
	}

	if node2 == nil {
		node2 = createNode(nodeAddr2, connCnt2)
	}

	chainClient, err := NewChainClient(
		WithChainClientOrgId(orgId1),
		WithChainClientChainId(chainId),
		WithChainClientLogger(getDefaultLogger()),
		WithUserKeyFilePath(userKeyPath),
		WithUserCrtFilePath(userCrtPath),
		AddChainClientNodeConfig(node1),
		AddChainClientNodeConfig(node2),
		)

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
		WithChainClientOrgId(orgId1),
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

	return adminClient, nil
}
