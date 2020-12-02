/**
 * @Author: jasonruan
 * @Date:   2020-12-01 14:49:44
 */
package chainmaker_sdk_go

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/stretchr/testify/require"
)

const (
	chainId        = "chain1"
	orgId          = "wx-org1.chainmaker.org"
	contractName   = "counter-go-1"
	certPathPrefix = "./testdata"
	tlsHostName    = "chainmaker.org"

	nodeAddr = "127.0.0.1:12301"
	connCnt  = 5

	multiSignedPayloadFile = "./testdata/counter-go-demo/collect-signed-all.pb"
)

var (
	caPaths     = []string{certPathPrefix + fmt.Sprintf("/crypto-config/%s/ca", orgId)}
	userKeyPath = certPathPrefix + fmt.Sprintf("/crypto-config/%s/user/client1/client1.tls.key", orgId)
	userCrtPath = certPathPrefix + fmt.Sprintf("/crypto-config/%s/user/client1/client1.tls.crt", orgId)
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

func TestUserContractCounterGo(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	testUserContractCounterGoCreate(t, client)
	time.Sleep(5 * time.Second)

	testUserContractCounterGoInvoke(t, client)
	time.Sleep(5 * time.Second)

	testUserContractCounterGoQuery(t, client)
}

func testUserContractCounterGoCreate(t *testing.T, client *ChainClient) {
	file, err := ioutil.ReadFile(multiSignedPayloadFile)
	require.Nil(t, err)

	resp, err := client.ContractCreate("", file)
	require.Nil(t, err)

	fmt.Printf("CREATE counter-go contract resp: %+v\n", resp)
}

func testUserContractCounterGoInvoke(t *testing.T, client *ChainClient) {
	resp, err := client.ContractInvoke(contractName, "increase", "", nil)
	require.Nil(t, err)
	fmt.Printf("INVOKE counter-go contract resp: %+v\n", resp)
}

func testUserContractCounterGoQuery(t *testing.T, client *ChainClient) {
	resp, err := client.ContractQuery(contractName, "query", "", nil)
	require.Nil(t, err)
	fmt.Printf("QUERY counter-go contract resp: %+v\n", resp)
}

func TestSystemContractGo(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	blockInfo := testSystemContractGoGetBlockByHeight(t, client, -1)
	testSystemContractGoGetTxByTxId(t, client, blockInfo.Block.Txs[0].Header.TxId)
	testSystemContractGoGetBlockByHash(t, client, hex.EncodeToString(blockInfo.Block.Header.BlockHash))
	testSystemContractGoGetBlockByTxId(t, client, blockInfo.Block.Txs[0].Header.TxId)
	testSystemContractGoGetLastConfigBlock(t, client)
	testSystemContractGoGetChainInfo(t, client)
	testSystemContractGoGetContractInfo(t, client)

	client, err = New(
		// 必填字段
		AddNodeAddrWithConnCnt(nodeAddr, connCnt),
		WithLogger(getDefaultLogger()),
		WithUserKeyFilePath(userKeyPath),
		WithUserCrtFilePath(userCrtPath),
		WithOrgId(orgId),
		WithChainId(SYSTEM_CHAIN),
		// 选填字段
		WithUseTLS(true),
		WithCAPaths(caPaths),
		WithTLSHostName(tlsHostName),
	)
	require.Nil(t, err)
	testSystemContractGoGetNodeChainList(t, client)
}

func testSystemContractGoGetTxByTxId(t *testing.T, client *ChainClient, txId string) *pb.TransactionInfo {
	transactionInfo, err := client.GetTxByTxId(txId)
	require.Nil(t, err)
	return transactionInfo
}

func testSystemContractGoGetBlockByHeight(t *testing.T, client *ChainClient, blockHeight int64) *pb.BlockInfo {
	blockInfo, err := client.GetBlockByHeight(blockHeight, true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGoGetBlockByHash(t *testing.T, client *ChainClient, blockHash string) *pb.BlockInfo {
	blockInfo, err := client.GetBlockByHash(blockHash, true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGoGetBlockByTxId(t *testing.T, client *ChainClient, txId string) *pb.BlockInfo {
	blockInfo, err := client.GetBlockByTxId(txId, true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGoGetLastConfigBlock(t *testing.T, client *ChainClient) *pb.BlockInfo {
	blockInfo, err := client.GetLastConfigBlock(true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGoGetChainInfo(t *testing.T, client *ChainClient) *pb.ChainInfo {
	chainInfo, err := client.GetChainInfo()
	require.Nil(t, err)
	return chainInfo
}

func testSystemContractGoGetContractInfo(t *testing.T, client *ChainClient) *pb.ContractInfo {
	contractInfo, err := client.GetContractInfo()
	require.Nil(t, err)
	return contractInfo
}

func testSystemContractGoGetNodeChainList(t *testing.T, client *ChainClient) *pb.ChainList {
	chainList, err := client.GetNodeChainList()
	require.Nil(t, err)
	return chainList
}
