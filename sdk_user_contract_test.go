/**
 * @Author: jasonruan
 * @Date:   2020-12-02 18:41:47
 **/
package chainmaker_sdk_go

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

func TestUserContractCounterGo(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	testUserContractCounterGoCreate(t, client)
	time.Sleep(5 * time.Second)

	testUserContractCounterGoInvoke(t, client)
	time.Sleep(5 * time.Second)

	testUserContractCounterGoQuery(t, client)

	testUserContractCounterGoUpgrade(t, client)
}

// [用户合约]
func testUserContractCounterGoCreate(t *testing.T, client *ChainClient) {
	file, err := ioutil.ReadFile(multiSignedPayloadFile)
	require.Nil(t, err)

	resp, err := client.CreateContract("", file)
	require.Nil(t, err)

	fmt.Printf("CREATE counter-go contract resp: %+v\n", resp)
}

func testUserContractCounterGoUpgrade(t *testing.T, client *ChainClient) {
	file, err := ioutil.ReadFile(upgradeMultiSignedPayloadFile)
	require.Nil(t, err)

	resp, err := client.UpgradeContract("", file)
	require.Nil(t, err)

	fmt.Printf("UPGRADE counter-go contract resp: %+v\n", resp)
}

func testUserContractCounterGoInvoke(t *testing.T, client *ChainClient) {
	resp, err := client.InvokeContract(contractName, "increase", "", nil)
	require.Nil(t, err)
	fmt.Printf("INVOKE counter-go contract resp: %+v\n", resp)
}

func testUserContractCounterGoQuery(t *testing.T, client *ChainClient) {
	resp, err := client.QueryContract(contractName, "query", nil)
	require.Nil(t, err)
	fmt.Printf("QUERY counter-go contract resp: %+v\n", resp)
}

