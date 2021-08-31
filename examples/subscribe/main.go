/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"fmt"
	"log"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
)

const (
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	go testSubscribeBlock(client, false)
	go testSubscribeBlock(client, true)
	go testSubscribeTx(client)
	go testSubscribeContractEvent(client)
	select {}
}

func testSubscribeBlock(client *sdk.ChainClient, onlyHeader bool) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := client.SubscribeBlock(ctx, 0, 10, true, onlyHeader)
	//c, err := client.SubscribeBlock(ctx, 10, -1, true, onlyHeader)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case block, ok := <-c:
			if !ok {
				fmt.Println("chan is close!")
				return
			}

			if block == nil {
				log.Fatalln("require not nil")
			}

			if onlyHeader {
				blockHeader, ok := block.(*common.BlockHeader)
				if !ok {
					log.Fatalln("require true")
				}

				fmt.Printf("recv blockHeader [%d] => %+v\n", blockHeader.BlockHeight, blockHeader)
			} else {
				blockInfo, ok := block.(*common.BlockInfo)
				if !ok {
					log.Fatalln("require true")
				}

				fmt.Printf("recv blockInfo [%d] => %+v\n", blockInfo.Block.Header.BlockHeight, blockInfo)
			}

			//if err := client.Stop(); err != nil {
			//	return
			//}
			//return
		case <-ctx.Done():
			return
		}
	}
}

func testSubscribeTx(client *sdk.ChainClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := client.SubscribeTx(ctx, -1, -1, "", nil)
	//c, err := client.SubscribeTx(ctx, 0, 4, "", nil)
	//c, err := client.SubscribeTx(ctx, 0, -1, "", []string{"fc4eac7ac478453aa486ab84e4f814df56b39ccfc5d9418d96a197781327bf31"})
	//c, err := client.SubscribeTx(ctx, 0, -1, syscontract.SystemContract_CERT_MANAGE.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case txI, ok := <-c:
			if !ok {
				fmt.Println("chan is close!")
				return
			}

			if txI == nil {
				log.Fatalln("require not nil")
			}

			tx, ok := txI.(*common.Transaction)
			if !ok {
				log.Fatalln("require true")
			}

			fmt.Printf("recv tx [%s] => %+v\n", tx.Payload.TxId, tx)

			//if err := client.Stop(); err != nil {
			//	return
			//}
			//return
		case <-ctx.Done():
			return
		}
	}
}

func testSubscribeContractEvent(client *sdk.ChainClient) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//订阅指定合约的合约事件
	c, err := client.SubscribeContractEvent(ctx, "topic_vx", "claim001")

	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case event, ok := <-c:
			if !ok {
				fmt.Println("chan is close!")
				return
			}
			if event == nil {
				log.Fatalln("require not nil")
			}
			contractEventInfo, ok := event.(*common.ContractEventInfo)
			if !ok {
				log.Fatalln("require true")
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
