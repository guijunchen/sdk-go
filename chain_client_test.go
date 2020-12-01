/**
 * @Author: jasonruan
 * @Date:   2020-12-01 14:49:44
 */
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	chainId         = "chain1"
	orgId           = "wx-org1.chainmaker.org"
	contractName    = "contract1"
	runtimeType     = pb.RuntimeType_WASMER_RUST
	certPathPrefix  = "/home/jason/Work/ChainMaker/chainmaker-go/config"
	tlsHostName     = "chainmaker.org"

	nodeAddr        = "127.0.0.1:12301"
	connCnt         = 5
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

func TestChainClient_ContractInvoke(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	paramsMap := make(map[string]string)
	paramsMap["aaa"] = "A1"
	paramsMap["bbb"] = "B1"

	client.ContractInvoke(contractName, "save", "", paramsMap)
}

