/**
 * @Author: jasonruan
 * @Date:   2020-12-01 14:49:44
 */
package chainmaker_sdk_go

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

const (
	chainId                 = "chain1"
	orgId                   = "wx-org1.chainmaker.org"
	contractName            = "counter-go-1"
	certPathPrefix          = "./testdata"
	tlsHostName             = "chainmaker.org"

	nodeAddr                = "127.0.0.1:12301"
	connCnt                 = 5

	multiSignedPayloadFile  = "./testdata/counter-go-demo/collect-signed-all.pb"
)

var (
	caPaths         = []string{certPathPrefix + fmt.Sprintf("/crypto-config/%s/ca", orgId)}
	userKeyPath     = certPathPrefix + fmt.Sprintf("/crypto-config/%s/user/client1/client1.tls.key", orgId)
	userCrtPath     = certPathPrefix + fmt.Sprintf("/crypto-config/%s/user/client1/client1.tls.crt", orgId)
)

func createClient() (*ChainClient, error) {
	client, err := New(
		// 必填字段
		AddNodeAddrWithConnCnt(nodeAddr, connCnt),
		WithLogger(getDefaultLogger()),
		WithUserKeyFilePath(userKeyPath),
		WithUserCrtFilePath(userCrtPath),
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

	return client, nil
}

func TestContractCounterGo(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	testContractCounterGoCreate(t, client)
	time.Sleep(5 * time.Second)

	testContractCounterGoInvoke(t, client)
	time.Sleep(5 * time.Second)

	testContractCounterGoQuery(t, client)
}

func testContractCounterGoCreate(t *testing.T, client *ChainClient) {
	file, err := ioutil.ReadFile(multiSignedPayloadFile)
	require.Nil(t, err)

	resp, err := client.ContractCreate("", file)
	require.Nil(t, err)

	fmt.Printf("CREATE counter-go contract resp: %+v\n", resp)
}

func testContractCounterGoInvoke(t *testing.T, client *ChainClient) {
	resp, err := client.ContractInvoke(contractName, "increase", "", nil)
	require.Nil(t, err)
	fmt.Printf("INVOKE counter-go contract resp: %+v\n", resp)
}

func testContractCounterGoQuery(t *testing.T, client *ChainClient) {
	resp, err := client.ContractQuery(contractName, "query", "", nil)
	require.Nil(t, err)
	fmt.Printf("QUERY counter-go contract resp: %+v\n", resp)
}