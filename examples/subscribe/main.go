/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"fmt"
	"time"

	"chainmaker.org/chainmaker/common/random/uuid"
	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
)

const (
	sendTxCount       = 5
	claimContractName = "claim001"

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	go testSubscribeBlock()
	go testSubscribeContractEvent()
	go testSubscribeTx()
	select {}
}

func testSubscribeBlock() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := client.SubscribeBlock(ctx, 0, 10, true)
	//c, err := client.SubscribeBlock(ctx, 5, 16, false)
	//c, err := client.SubscribeBlock(ctx, 0, -1, false)
	//c, err := client.SubscribeBlock(ctx, 10, -1, false)
	if err != nil {
		panic(err)
	}

	go func() {
		for i := 0; i < sendTxCount; i++ {
			_, err := testUserContractClaimInvoke(client, "save", false)
			if err != nil {
				panic(err)
			}
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

			if block == nil {
				panic("require not nil")
			}

			blockInfo, ok := block.(*common.BlockInfo)
			if !ok {
				panic("require true")
			}

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

func testSubscribeContractEvent() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//订阅指定合约的合约事件
	c, err := client.SubscribeContractEvent(ctx, "topic_vx", "claim001")

	if err != nil {
		panic(err)
	}

	go func() {
		for i := 0; i < sendTxCount; i++ {
			_, err := testUserContractClaimInvoke(client, "save", false)
			if err != nil {
				panic(err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	for {
		select {
		case event, ok := <-c:
			if !ok {
				fmt.Println("chan is close!")
				return
			}
			if event == nil {
				panic("require not nil")
			}
			contractEventInfo, ok := event.(*common.ContractEventInfo)
			if !ok {
				panic("require true")
			}
			fmt.Printf("recv contract event [%d] => %+v\n", contractEventInfo.BlockHeight, contractEventInfo)

			//if err := client.Stop(); err != nil {
			//	return
			//}
			//return
		case <-ctx.Done():
			return
		}
	}
}

func testSubscribeTx() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := client.SubscribeTx(ctx, -1, -1, -1, nil)
	//c, err := client.SubscribeTx(ctx, 0, 18, -1, nil)
	//c, err := client.SubscribeTx(ctx, 50, -1, -1, nil)
	//c, err := client.SubscribeTx(ctx, 0, 0, -1, []string{"04e98331c02d423c91e5b0bb9b9f8519112d6cee26d94620a3c9773a5ce19147"})
	//c, err := client.SubscribeTx(ctx, -1, -1, common.TxType_INVOKE_USER_CONTRACT, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		for i := 0; i < sendTxCount; i++ {
			_, err := testUserContractClaimInvoke(client, "save", false)
			if err != nil {
				panic(err)
			}
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

			if txI == nil {
				panic("require not nil")
			}

			tx, ok := txI.(*common.Transaction)
			if !ok {
				panic("require true")
			}

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

func testUserContractClaimInvoke(client *sdk.ChainClient, method string, withSyncResult bool) (string, error) {
	//curTime := fmt.Sprintf("%d", CurrentTimeMillisSeconds())
	curTime := time.Now().Format("2006-01-02 15:04:05")

	fileHash := uuid.GetUUID()
	params := map[string]string{
		"time":      curTime,
		"file_hash": fileHash,
		"file_name": fmt.Sprintf("file_%s", curTime),
	}

	err := invokeUserContract(client, claimContractName, method, "", params, withSyncResult)
	//err := invokeUserContractStepByStep(client, claimContractName, method, "", params, withSyncResult)
	if err != nil {
		return "", err
	}

	return fileHash, nil
}

func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string, params map[string]string,
	withSyncResult bool) error {

	resp, err := client.InvokeContract(contractName, method, txId, params, -1, withSyncResult)
	if err != nil {
		return err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if !withSyncResult {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
	} else {
		fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
	}

	return nil
}
