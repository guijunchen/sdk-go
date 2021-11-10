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

	c, err := client.SubscribeTx(ctx, 10, 30, "", nil)
	//c, err := client.SubscribeTx(ctx, 0, 4, "", nil)
	//c, err := client.SubscribeTx(ctx, 0, -1, "", []string{"1b70bb886c784a0587590da3a0af8fd336aa1a806be4431db31ceeba4a912f93"})
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
	c, err := client.SubscribeContractEvent(ctx, 0, -1, "claim005", "topic_vx")
	//c, err := client.SubscribeContractEvent(ctx, 0, 0, "claim005", "")
	//c, err := client.SubscribeContractEvent(ctx, "64f50d594c2a739c7088f9fc6785e1934030e17b52f1a894baec61b98633a59f", "9c01b4c21d1907ab27aa23343493b3c9872777e3")

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
