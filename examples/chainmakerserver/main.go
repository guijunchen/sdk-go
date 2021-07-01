/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"

	"chainmaker.org/chainmaker/sdk-go/examples"
)

func main() {
	testChainClientGetChainMakerServerVersion()
}

func testChainClientGetChainMakerServerVersion() {
	client, err := examples.CreateClient()
	if err != nil {
		panic(err)
	}
	version, err := client.GetChainMakerServerVersion()
	if err != nil {
		panic(err)
	}
	fmt.Println("get chainmaker server version:", version)
}
