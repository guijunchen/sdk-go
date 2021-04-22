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

	client1Addr       = "0xbe8cd0f69425253c3bab99344d61eb34f57bdc21"
	client2Addr       = "0xa55f1e0cb68b0cc589906078237094bdb9715bfd"

	client1AddrInt    = "1087848554046178479107522336262214072175637027873"
	client2AddrInt    = "944104665674401770091203869615921096651560803325"

	client1AddrSki    = "7081212378e72d4ecf406c30384f82a74c2a0c8d9e91ccfa94c245023942240f"
	client2AddrSki    = "320cd73d87b5b238a2d09cce54bc6796288e4ce8760fb0561ed39db68044d3a6"

	amount            = 200
)

func TestUserContractTokenEVM(t *testing.T) {
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

	fmt.Println("====================== 创建Token合约,给client1地址分配初始代币 ======================")
	testUserContractTokenEVMCreate(t, client, admin1, admin2, admin3, admin4, true, true)

	fmt.Println("====================== 查看余额 ======================")
	testUserContractTokenEVMBalanceOf(t, client, client1AddrInt, true)
	testUserContractTokenEVMBalanceOf(t, client, client2AddrInt, true)

	fmt.Println("====================== client1给client2地址转账 ======================")
	testUserContractTokenEVMTransfer(t, client, amount, true)

	fmt.Println("====================== 查看余额 ======================")
	testUserContractTokenEVMBalanceOf(t, client, client1AddrInt, true)
	testUserContractTokenEVMBalanceOf(t, client, client2AddrInt, true)
}

func testUserContractTokenEVMCreate(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient, withSyncResult bool, isIgnoreSameContract bool) {

	abiJson, err := ioutil.ReadFile(tokenABIPath)
	require.Nil(t, err)

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	require.Nil(t, err)

	// 方式1: 16进制
	//addr := evmutils.BigToAddress(evmutils.FromHexString(client1Addr[2:]))
	// 方式2: Int
	//addr := evmutils.BigToAddress(evmutils.FromDecimalString(client1AddrInt))
	// 方式3: ski
	addrInt, err := evmutils.MakeAddressFromHex(client1AddrSki)
	require.Nil(t, err)
	addr := evmutils.BigToAddress(addrInt)

	dataByte, err := myAbi.Pack("", addr)
	require.Nil(t, err)

	data := hex.EncodeToString(dataByte)
	pairs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: data,
		},
	}

	byteCode, err := ioutil.ReadFile(tokenByteCodePath)
	require.Nil(t, err)

	//bc := string(byteCode)
	//bc = strings.TrimSpace(bc)

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		//tokenContractName, tokenVersion, bc + data, common.RuntimeType_EVM, pairs, withSyncResult)
		tokenContractName, tokenVersion, string(byteCode), common.RuntimeType_EVM, pairs, withSyncResult)
	if !isIgnoreSameContract {
		require.Nil(t, err)
	}

	fmt.Printf("CREATE EVM token contract resp: %+v\n", resp)
}

func testUserContractTokenEVMTransfer(t *testing.T, client *ChainClient, amount int64, withSyncResult bool) {

	abiJson, err := ioutil.ReadFile(tokenABIPath)
	require.Nil(t, err)

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	require.Nil(t, err)

	//addr := evmutils.StringToAddress(client2Addr)
	addr := evmutils.BigToAddress(evmutils.FromDecimalString(client2AddrInt))

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
	addr := evmutils.BigToAddress(evmutils.FromDecimalString(address))

	methodName := "balanceOf"
	dataByte, err := myAbi.Pack(methodName, addr)
	require.Nil(t, err)

	data := hex.EncodeToString(dataByte)
	method := data[0:8]

	pairs := map[string]string{
		"data": data,
	}

	result, err := invokeUserContractWithResult(client, tokenContractName, method, "", pairs, withSyncResult)
	require.Nil(t, err)

	balance, err := myAbi.Unpack(methodName, result)
	require.Nil(t, err)
	fmt.Printf("addr [%s] => %d\n", address, balance)
}