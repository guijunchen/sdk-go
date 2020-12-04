/**
 * @Author: jasonruan
 * @Date:   2020-12-01 14:49:44
 */
package chainmaker_sdk_go

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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

	nodeAddr = "127.0.0.1:12301"
	connCnt  = 5

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

func createClient() (*ChainClient, error) {
	client, err := New(
		// 必填字段
		AddNodeAddrWithConnCnt(nodeAddr, connCnt),
		WithLogger(getDefaultLogger()),
		WithUserKeyFilePath(userKeyPath),
		WithUserCrtFilePath(userCrtPath),
		WithOrgId(orgId1),
		WithChainId(chainId),
		// 选填字段
		WithUseTLS(true),
		WithCAPaths(caPaths),
		WithTLSHostName(tlsHostName),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createAdmin(orgId string) (*ChainClient, error) {
	admin, err := New(
		// 必填字段
		AddNodeAddrWithConnCnt(nodeAddr, connCnt),
		WithLogger(getDefaultLogger()),
		WithUserKeyFilePath(fmt.Sprintf(adminKeyPath, orgId)),
		WithUserCrtFilePath(fmt.Sprintf(adminCrtPath, orgId)),
		WithOrgId(orgId),
		WithChainId(chainId),
		// 选填字段
		WithUseTLS(true),
		WithCAPaths(caPaths),
		WithTLSHostName(tlsHostName),
	)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func TestMy(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	txId := "c0044d5082394e498168d401e3af464396f8e73e47e847e391906f19355b01ce"
	tx := testSystemContractGetTxByTxId(t, client, txId)
	fmt.Printf("tx: %+v\n", tx)
}
