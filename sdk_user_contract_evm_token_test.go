/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/common/evmutils"
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
	tokenContractName = "token001"
	tokenVersion      = "1.0.0"
	tokenByteCodePath = "./testdata/token-evm-demo/token.bin"
	tokenABIPath      = "./testdata/token-evm-demo/token.abi"

	//client1Addr       = "0xbe8cd0f69425253c3bab99344d61eb34f57bdc21"
	//client2Addr       = "0xa55f1e0cb68b0cc589906078237094bdb9715bfd"

	client1Addr         = "1087848554046178479107522336262214072175637027873"
	client2Addr         = "944104665674401770091203869615921096651560803325"

	amount            = 500000
)

func TestUserContractTokenEVM(t *testing.T) {
	fmt.Println("====================== create client ======================")
	client, err := createClientWithCertBytes()
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


	fmt.Println("====================== 创建Token合约,给client1地址分配初始代币 ======================")
	testUserContractTokenEVMCreate(t, client, admin1, admin2, admin3, admin4, true, true)

	//fmt.Println("====================== client1给client2地址转账200 ======================")
	//testUserContractTokenEVMTransfer(t, client, true)

	//fmt.Println("====================== 查看余额 ======================")
	testUserContractTokenEVMBalanceOf(t, client, client1Addr, true)
	//testUserContractTokenEVMBalanceOf(t, client, client2Addr, true)
}

func testUserContractTokenEVMCreate(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient, withSyncResult bool, isIgnoreSameContract bool) {

	// 构造合约参数，RLP编码
	//addr := evmutils.StringToAddress(client1Addr)

	//params := []interface{} {
	//	evmCommon.FromHex(client1Addr),
	//	//big.NewInt(amount),
	//}
	//
	//b, err := rlp.EncodeToBytes(params)
	//pairs := []*common.KeyValuePair{
	//	{
	//		Key:   "data",
	//		Value: hex.EncodeToString(b),
	//	},
	//}

	abiJson, err := ioutil.ReadFile(tokenABIPath)
	require.Nil(t, err)

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	require.Nil(t, err)

	//addr := evmutils.StringToAddress(client1Addr)
	addr := evmutils.BigToAddress(evmutils.FromDecimalString(client1Addr))

	dataByte, err := myAbi.Pack("", addr)
	require.Nil(t, err)

	data := hex.EncodeToString(dataByte)
	pairs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: data,
		},
	}

	codeBytes, err := ioutil.ReadFile(tokenByteCodePath)
	require.Nil(t, err)

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		tokenContractName, tokenVersion, string(codeBytes), common.RuntimeType_EVM, pairs, withSyncResult)
	if !isIgnoreSameContract {
		require.Nil(t, err)
	}

	fmt.Printf("CREATE EVM token contract resp: %+v\n", resp)
}

func testUserContractTokenEVMTransfer(t *testing.T, client *ChainClient, withSyncResult bool) {

	abiJson, err := ioutil.ReadFile(tokenABIPath)
	require.Nil(t, err)

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	require.Nil(t, err)

	//addr := evmutils.StringToAddress(client2Addr)
	addr := evmutils.BigToAddress(evmutils.FromDecimalString(client2Addr))

	dataByte, err := myAbi.Pack("transfer", addr, big.NewInt(amount))
	require.Nil(t, err)

	data := hex.EncodeToString(dataByte)
	method := data[0:8]

	pairs := map[string]string{
		"data": data,
	}

	err = invokeUserContract(client, tokenContractName, method, "", pairs, withSyncResult)
	require.Nil(t, err)
}

func testUserContractTokenEVMBalanceOf(t *testing.T, client *ChainClient, address string, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(tokenABIPath)
	require.Nil(t, err)

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	require.Nil(t, err)

	//addr := evmutils.StringToAddress(address)
	addr := evmutils.BigToAddress(evmutils.FromDecimalString(client1Addr))

	dataByte, err := myAbi.Pack("balanceOf", addr)
	require.Nil(t, err)

	data := hex.EncodeToString(dataByte)
	method := data[0:8]

	pairs := map[string]string{
		"data": data,
	}

	err = invokeUserContract(client, tokenContractName, method, "", pairs, withSyncResult)
	require.Nil(t, err)
}