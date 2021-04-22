/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"
)

const (
	storageContractName = "storage001"
	storageVersion      = "1.0.0"
	storageByteCodePath = "./testdata/storage-evm-demo/storage.bin"
	storageABIPath      = "./testdata/storage-evm-demo/storage.abi"
)

func TestUserContractStorageEVM(t *testing.T) {
	fmt.Println("====================== create client ======================")
	client, err := createClientWithConfig()
	require.Nil(t, err)

	fmt.Println("====================== create admin1 ======================")
	admin1, err := createAdmin(orgId1)
	require.Nil(t, err)
	fmt.Println("====================== create admin2 ======================")
	admin2, err := createAdmin(orgId2)
	require.Nil(t, err)
	fmt.Println("====================== create admin3 ======================")
	admin3, err := createAdmin(orgId3)
	require.Nil(t, err)
	fmt.Println("====================== create admin4 ======================")
	admin4, err := createAdmin(orgId4)
	require.Nil(t, err)


	fmt.Println("====================== 创建Storage合约 ======================")
	testUserContractStorageEVMCreate(t, client, admin1, admin2, admin3, admin4, true, true)

	fmt.Println("====================== 设置数值 ======================")
	testUserContractStorageEVMSet(t, client, 123, true)

	fmt.Println("====================== 查看数值 ======================")
	testUserContractStorageEVMGet(t, client, true)
}

func testUserContractStorageEVMCreate(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient, withSyncResult bool, isIgnoreSameContract bool) {

	codeBytes, err := ioutil.ReadFile(storageByteCodePath)
	require.Nil(t, err)

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		storageContractName, storageVersion, string(codeBytes), common.RuntimeType_EVM, nil, withSyncResult)
	if !isIgnoreSameContract {
		require.Nil(t, err)
	}

	fmt.Printf("CREATE EVM storage contract resp: %+v\n", resp)
}

func testUserContractStorageEVMSet(t *testing.T, client *ChainClient, data int64, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(storageABIPath)
	require.Nil(t, err)

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	require.Nil(t, err)

	dataByte, err := myAbi.Pack("set", big.NewInt(data))
	require.Nil(t, err)

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	pairs := map[string]string{
		"data": dataString,
	}

	err = invokeUserContract(client, storageContractName, method, "", pairs, withSyncResult)
	require.Nil(t, err)
}
func testUserContractStorageEVMGet(t *testing.T, client *ChainClient, withSyncResult bool) {

	abiJson, err := ioutil.ReadFile(storageABIPath)
	require.Nil(t, err)

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	require.Nil(t, err)

	dataByte, err := myAbi.Pack("get")
	require.Nil(t, err)

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	pairs := map[string]string{
		"data": dataString,
	}

	result, err := invokeUserContractWithResult(client, storageContractName, method, "", pairs, withSyncResult)
	require.Nil(t, err)

	val, err := myAbi.Unpack("get", result)
	require.Nil(t, err)
	fmt.Printf("val: %d\n", val)
}