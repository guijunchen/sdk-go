package chainmaker_sdk_go

import (
	"encoding/hex"
	"testing"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/stretchr/testify/require"
)

// [系统合约]
func TestSystemContract(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	blockInfo := testSystemContractGetBlockByHeight(t, client, -1)
	testSystemContractGetTxByTxId(t, client, blockInfo.Block.Txs[0].Header.TxId)
	testSystemContractGetBlockByHash(t, client, hex.EncodeToString(blockInfo.Block.Header.BlockHash))
	testSystemContractGetBlockByTxId(t, client, blockInfo.Block.Txs[0].Header.TxId)
	testSystemContractGetLastConfigBlock(t, client)
	testSystemContractGetChainInfo(t, client)
	testSystemContractGetContractInfo(t, client)

	client, err = New(
		// 必填字段
		AddNodeAddrWithConnCnt(nodeAddr, connCnt),
		WithLogger(getDefaultLogger()),
		WithUserKeyFilePath(userKeyPath),
		WithUserCrtFilePath(userCrtPath),
		WithOrgId(orgId1),
		WithChainId(SYSTEM_CHAIN),
		// 选填字段
		WithUseTLS(true),
		WithCAPaths(caPaths),
		WithTLSHostName(tlsHostName),
	)
	require.Nil(t, err)
	testSystemContractGetNodeChainList(t, client)
}

func testSystemContractGetTxByTxId(t *testing.T, client *ChainClient, txId string) *pb.TransactionInfo {
	transactionInfo, err := client.GetTxByTxId(txId)
	require.Nil(t, err)
	return transactionInfo
}

func testSystemContractGetBlockByHeight(t *testing.T, client *ChainClient, blockHeight int64) *pb.BlockInfo {
	blockInfo, err := client.GetBlockByHeight(blockHeight, true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGetBlockByHash(t *testing.T, client *ChainClient, blockHash string) *pb.BlockInfo {
	blockInfo, err := client.GetBlockByHash(blockHash, true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGetBlockByTxId(t *testing.T, client *ChainClient, txId string) *pb.BlockInfo {
	blockInfo, err := client.GetBlockByTxId(txId, true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGetLastConfigBlock(t *testing.T, client *ChainClient) *pb.BlockInfo {
	blockInfo, err := client.GetLastConfigBlock(true)
	require.Nil(t, err)
	return blockInfo
}

func testSystemContractGetChainInfo(t *testing.T, client *ChainClient) *pb.ChainInfo {
	chainInfo, err := client.GetChainInfo()
	require.Nil(t, err)
	return chainInfo
}

func testSystemContractGetContractInfo(t *testing.T, client *ChainClient) *pb.ContractInfo {
	contractInfo, err := client.GetContractInfo()
	require.Nil(t, err)
	return contractInfo
}

func testSystemContractGetNodeChainList(t *testing.T, client *ChainClient) *pb.ChainList {
	chainList, err := client.GetNodeChainList()
	require.Nil(t, err)
	return chainList
}
