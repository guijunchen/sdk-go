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

func TestContractCreate(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	file, err := ioutil.ReadFile(multiSignedPayloadFile)
	require.Nil(t, err)

	resp, err := client.ContractCreate("", file)
	require.Nil(t, err)

	fmt.Printf("resp: %+v\n", resp)
}

func TestContractInvoke(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	//paramsMap := make(map[string]string)
	//paramsMap["aaa"] = "A1"
	//paramsMap["bbb"] = "B1"

	resp, err := client.ContractInvoke(contractName, "increase", "", nil)
	require.Nil(t, err)

	fmt.Printf("resp: %+v\n", resp)
}

func TestContractQuery(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	resp, err := client.ContractQuery(contractName, "query", "", nil)
	require.Nil(t, err)

	fmt.Printf("resp: %+v\n", resp)
}