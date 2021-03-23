/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/common/random/uuid"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	claimContractName = "claim001"
	claimVersion      = "1.0.0"
	claimByteCodePath = "./testdata/claim-wasm-demo/rust-fact-1.0.0.wasm"
)

func TestUserContractClaim(t *testing.T) {
	client, err := createClientWithConfig()
	require.Nil(t, err)

	admin1, err := createAdmin(orgId1)
	require.Nil(t, err)
	admin2, err := createAdmin(orgId2)
	require.Nil(t, err)
	admin3, err := createAdmin(orgId3)
	require.Nil(t, err)
	admin4, err := createAdmin(orgId4)
	require.Nil(t, err)

	fmt.Println("====================== 创建合约 ======================")
	testUserContractClaimCreate(t, client, admin1, admin2, admin3, admin4, true, true)

	fmt.Println("====================== 调用合约 ======================")
	fileHash, err := testUserContractClaimInvoke(client, "save", true)
	require.Nil(t, err)

	fmt.Println("====================== 执行合约查询接口 ======================")
	//txId := "1cbdbe6106cc4132b464185ea8275d0a53c0261b7b1a470fb0c3f10bd4a57ba6"
	//fileHash = txId[len(txId)/2:]
	params := map[string]string{
		"file_hash": fileHash,
	}
	testUserContractClaimQuery(t, client, "find_by_file_hash", params)
}

func testUserContractClaimCreate(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient, withSyncResult bool, isIgnoreSameContract bool) {

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		claimContractName, claimVersion, claimByteCodePath, common.RuntimeType_WASMER, []*common.KeyValuePair{}, withSyncResult)
	if !isIgnoreSameContract {
		require.Nil(t, err)
	}

	fmt.Printf("CREATE claim contract resp: %+v\n", resp)
}

func testUserContractClaimInvoke(client *ChainClient,
	method string, withSyncResult bool) (string, error) {

	curTime := fmt.Sprintf("%d", CurrentTimeMillisSeconds())
	fileHash := uuid.GetUUID()
	params := map[string]string{
		"time":      curTime,
		"file_hash": fileHash,
		"file_name": fmt.Sprintf("file_%s", curTime),
	}

	err := invokeUserContract(client, claimContractName, method, "", params, withSyncResult)
	if err != nil {
		return "", err
	}

	return fileHash, nil
}

func testUserContractClaimQuery(t *testing.T, client *ChainClient,
	method string, params map[string]string) {
	resp, err := client.QueryContract(claimContractName, method, params, -1)
	require.Nil(t, err)
	fmt.Printf("QUERY claim contract resp: %+v\n", resp)
}
