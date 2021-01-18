/**
 * @Author: jasonruan
 * @Date:   2021-01-16 15:22:24
 **/
package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"encoding/binary"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	assetContractName = "asset001"
	assetVersion = "1.0.0"
	assetByteCodePath = "./testdata/asset-wasm-demo/asset-rust-0.7.2.wasm"
)

func TestUserContractAsset(t *testing.T) {
	client, err := createClientWithConfig()
	require.Nil(t, err)

	client2, err := createClientWithOrgId(orgId2)
	require.Nil(t, err)

	admin1, err := createAdmin(orgId1)
	require.Nil(t, err)
	admin2, err := createAdmin(orgId2)
	require.Nil(t, err)
	admin3, err := createAdmin(orgId3)
	require.Nil(t, err)
	admin4, err := createAdmin(orgId4)
	require.Nil(t, err)

	fmt.Println("====================== 1)安装钱包合约 ======================")
	pairs := []*pb.KeyValuePair{
		{
			Key:   "issue_limit",
			Value: "100000000",
		},
		{
			Key:   "total_supply",
			Value: "100000000",
		},
	}
	testUserContractAssetCreate(t, client, admin1, admin2, admin3, admin4, pairs, true, true)

	fmt.Println("====================== 2)注册另一个用户 ======================")
	testUserContractAssetInvokeRegister(t, client2, "register", true)

	fmt.Println("====================== 3)查询钱包地址 ======================")
	addr1 := testUserContractAssetQuery(t, client, "query_address", nil)
	fmt.Printf("client1 address: %s\n", addr1)
	addr2 := testUserContractAssetQuery(t, client2, "query_address", nil)
	fmt.Printf("client2 address: %s\n", addr2)

	fmt.Println("====================== 4)给用户分别发币100000 ======================")
	amount := "100000"
	testUserContractAssetInvoke(t, client, "issue_amount", amount, addr1, true)
	testUserContractAssetInvoke(t, client, "issue_amount", amount, addr2, true)

	fmt.Println("====================== 5)分别查看余额 ======================")
	params := map[string]string {
		"owner": addr1,
	}
	balance1 := testUserContractAssetQuery(t, client, "balance_of", params)
	val1, err := BytesToInt([]byte(balance1), binary.LittleEndian)
	require.Nil(t, err)
	fmt.Printf("client1 balance: %d\n", val1)

	params = map[string]string {
		"owner": addr2,
	}
	balance2 := testUserContractAssetQuery(t, client, "balance_of", params)
	val2, err := BytesToInt([]byte(balance2), binary.LittleEndian)
	require.Nil(t, err)
	fmt.Printf("client2 balance: %d\n", val2)

	fmt.Println("====================== 6)A给B转账100 ======================")
	amount = "100"
	testUserContractAssetInvoke(t, client, "transfer", amount, addr2, true)

	fmt.Println("====================== 7)再次分别查看余额 ======================")
	params = map[string]string {
		"owner": addr1,
	}
	balance1 = testUserContractAssetQuery(t, client, "balance_of", params)
	val1, err = BytesToInt([]byte(balance1), binary.LittleEndian)
	require.Nil(t, err)
	fmt.Printf("client1 balance: %d\n", val1)

	params = map[string]string {
		"owner": addr2,
	}
	balance2 = testUserContractAssetQuery(t, client, "balance_of", params)
	val2, err = BytesToInt([]byte(balance2), binary.LittleEndian)
	require.Nil(t, err)
	fmt.Printf("client2 balance: %d\n", val2)
}

func testUserContractAssetCreate(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient, kvs []*pb.KeyValuePair, withSyncResult bool, isIgnoreSameContract bool) {


	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		assetContractName, assetVersion, assetByteCodePath, pb.RuntimeType_WASMER, kvs, withSyncResult)
	if !isIgnoreSameContract {
		require.Nil(t, err)
	}

	fmt.Printf("CREATE asset contract resp: %+v\n", resp)
}

func testUserContractAssetInvokeRegister(t *testing.T, client *ChainClient, method string, withSyncResult bool) {
	err := invokeUserContract(client, assetContractName, method, "", nil, withSyncResult)
	require.Nil(t, err)
}

func testUserContractAssetQuery(t *testing.T, client *ChainClient, method string, params map[string]string) string {
	resp, err := client.QueryContract(assetContractName, method, params, -1)
	require.Nil(t, err)
	fmt.Printf("QUERY asset contract [%s] resp: %+v\n", method, resp)

	err = checkProposalRequestResp(resp, true)
	require.Nil(t, err)
	return string(resp.ContractResult.Result)
}

func testUserContractAssetInvoke(t *testing.T, client *ChainClient, method string, amount, addr string, withSyncResult bool) {
	params := map[string]string {
		"amount": amount,
		"to": addr,
	}
	err := invokeUserContract(client, assetContractName, method, "", params, withSyncResult)
	require.Nil(t, err)
}

