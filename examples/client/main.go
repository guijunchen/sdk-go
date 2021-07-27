/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"log"

	"chainmaker.org/chainmaker/sdk-go/examples"
)

const (
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	testChainClientGetEVMAddressFromCertPath()
}

func testChainClientGetEVMAddressFromCertPath() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		log.Fatalln(err)
	}

	userOrg1Client1, err := examples.GetUser(examples.UserNameOrg1Client1)
	if err != nil {
		log.Fatalln(err)
	}

	addrInt, err := client.GetEVMAddressFromCertPath(userOrg1Client1.SignCrtPath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("client1 addrInt: %s\n", addrInt)

	userOrg2Client1, err := examples.GetUser(examples.UserNameOrg2Client1)
	if err != nil {
		log.Fatalln(err)
	}

	addrInt, err = client.GetEVMAddressFromCertPath(userOrg2Client1.SignCrtPath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("client2 addrInt: %s\n", addrInt)
}
