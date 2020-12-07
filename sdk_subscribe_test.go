/**
 * @Author: jasonruan
 * @Date:   2020-12-03 11:31:40
 **/
package chainmaker_sdk_go

import (
	"context"
	"fmt"
	"testing"
	"time"

	"chainmaker.org/chainmaker-go/chainmaker-sdk-go/pb"
	"github.com/stretchr/testify/require"
)

const (
	sendTxCount = 10
)

func TestSubscribeBlock(t *testing.T) {
	client, err := createClient()
	require.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := client.SubscribeBlock(ctx, 45, 55, false)
	require.Nil(t, err)

	go func() {
		for i := 0; i < sendTxCount; i++ {
			testUserContractCounterGoInvoke(t, client, "increase", nil, false)
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

			blockInfo, ok := block.(*pb.BlockInfo)
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

	c, err := client.SubscribeTx(ctx, 45, 55, -1, nil)
	//c, err := client.SubscribeTx(ctx, 45, 55, -1, []string{"b7bd37a15fbc49998612bd85b0c918796e3c12eae7384945bf7a82bc523b796d"})
	//c, err := client.SubscribeTx(ctx, -1, -1, pb.TxType_CREATE_USER_CONTRACT, nil)
	require.Nil(t, err)

	go func() {
		for i := 0; i < sendTxCount; i++ {
			testUserContractCounterGoInvoke(t, client, "increase", nil, false)
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

			tx, ok := txI.(*pb.Transaction)
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
