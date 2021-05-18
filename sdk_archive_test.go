/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestArchive(t *testing.T) {

	client1, err := createClientWithOrgId(orgId1)
	require.Nil(t, err)

	admin1, err := createAdmin(orgId1)
	require.Nil(t, err)

	fmt.Println("====================== 数据归档 ======================")
	var targetBlockHeight int64 = 4
	testArchiveBlock(t, client1, admin1, targetBlockHeight)
}

func testArchiveBlock(t *testing.T, client1, admin1 *ChainClient, targetBlockHeight int64) {
	var (
		err                error
		payload            []byte
		signedPayloadBytes []byte
		resp               *common.TxResponse
		result             string
	)

	payload, err = client1.CreateArchiveBlockPayload(targetBlockHeight)
	require.Nil(t, err)

	signedPayloadBytes, err = admin1.SignArchivePayload(payload)
	require.Nil(t, err)

	resp, err = client1.SendArchiveBlockRequest(signedPayloadBytes, -1, true)
	require.Nil(t, err)

	result = string(resp.ContractResult.Result)

	fmt.Printf("resp: %+v, result:%+s\n", resp, result)
}

