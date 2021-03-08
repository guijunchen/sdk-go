/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-go/pb/common"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	sendTxCount = 5
)

func TestSubscribeBlock(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := client.SubscribeBlock(ctx, 0, 10, true)
	//c, err := client.SubscribeBlock(ctx, 5, 16, false)
	//c, err := client.SubscribeBlock(ctx, 0, -1, false)
	//c, err := client.SubscribeBlock(ctx, 10, -1, false)
	require.Nil(t, err)

	go func() {
		for i := 0; i < sendTxCount; i++ {
			_, err := testUserContractClaimInvoke(client, "save", false)
			require.Nil(t, err)
			time.Sleep(2 * time.Second)
		}
	}()

	for {
		select {
		case block, ok := <-c:
			if !ok {
				fmt.Println("chan is close!")
				return
			}

			require.NotNil(t, block)

			blockInfo, ok := block.(*common.BlockInfo)
			require.Equal(t, true, ok)

			fmt.Printf("recv block [%d] => %+v\n", blockInfo.Block.Header.BlockHeight, blockInfo)

			//if err := client.Stop(); err != nil {
			//	return
			//}
			//return
		case <-ctx.Done():
			return
		}
	}
}

func TestSubscribeTx(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := client.SubscribeTx(ctx, -1, -1, -1, nil)
	//c, err := client.SubscribeTx(ctx, 0, 18, -1, nil)
	//c, err := client.SubscribeTx(ctx, 50, -1, -1, nil)
	//c, err := client.SubscribeTx(ctx, 0, 0, -1, []string{"04e98331c02d423c91e5b0bb9b9f8519112d6cee26d94620a3c9773a5ce19147"})
	//c, err := client.SubscribeTx(ctx, -1, -1, common.TxType_INVOKE_USER_CONTRACT, nil)
	require.Nil(t, err)

	go func() {
		for i := 0; i < sendTxCount; i++ {
			_, err := testUserContractClaimInvoke(client, "save", false)
			require.Nil(t, err)
			time.Sleep(2 * time.Second)
		}
	}()

	for {
		select {
		case txI, ok := <-c:
			if !ok {
				fmt.Println("chan is close!")
				return
			}

			require.NotNil(t, txI)

			tx, ok := txI.(*common.Transaction)
			require.Equal(t, true, ok)

			fmt.Printf("recv tx [%s] => %+v\n", tx.Header.TxId, tx)

			//if err := client.Stop(); err != nil {
			//	return
			//}
			//return
		case <-ctx.Done():
			return
		}
	}
}
