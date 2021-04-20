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
	balanceContractName = "balance001"
	balanceVersion      = "1.0.0"
	balanceByteCodePath = "./testdata/balance-evm-demo/ledger_balance.bin"
	balanceABIPath      = "./testdata/balance-evm-demo/ledger_balance.abi"
)

func TestUserContractBalanceEVM(t *testing.T) {
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


	fmt.Println("====================== 创建Balance合约 ======================")
	testUserContractBalanceEVMCreate(t, client, admin1, admin2, admin3, admin4, true, true)

	fmt.Println("====================== 设置余额 ======================")
	testUserContractBalanceEVMUpdateBalance(t, client, 123, true)

	fmt.Println("====================== 查看余额 ======================")
	testUserContractBalanceEVMGetBalance(t, client, client1Addr, true)
}

func testUserContractBalanceEVMCreate(t *testing.T, client *ChainClient,
	admin1, admin2, admin3, admin4 *ChainClient, withSyncResult bool, isIgnoreSameContract bool) {

	codeBytes, err := ioutil.ReadFile(balanceByteCodePath)
	require.Nil(t, err)

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		balanceContractName, balanceVersion, string(codeBytes), common.RuntimeType_EVM, nil, withSyncResult)
	if !isIgnoreSameContract {
		require.Nil(t, err)
	}

	fmt.Printf("CREATE EVM balance contract resp: %+v\n", resp)
}

func testUserContractBalanceEVMUpdateBalance(t *testing.T, client *ChainClient, data int64, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(balanceABIPath)
	require.Nil(t, err)

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	require.Nil(t, err)

	dataByte, err := myAbi.Pack("updateBalance", big.NewInt(data))
	require.Nil(t, err)

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	pairs := map[string]string{
		"data": dataString,
	}

	err = invokeUserContract(client, balanceContractName, method, "", pairs, withSyncResult)
	require.Nil(t, err)
}

func testUserContractBalanceEVMGetBalance(t *testing.T, client *ChainClient, address string, withSyncResult bool) {
	abiJson, err := ioutil.ReadFile(balanceABIPath)
	require.Nil(t, err)

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	require.Nil(t, err)

	addr := evmutils.BigToAddress(evmutils.FromDecimalString(address))

	dataByte, err := myAbi.Pack("balances", addr)
	require.Nil(t, err)

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]

	pairs := map[string]string{
		"data": dataString,
	}

	err = invokeUserContract(client, balanceContractName, method, "", pairs, withSyncResult)
	require.Nil(t, err)
}