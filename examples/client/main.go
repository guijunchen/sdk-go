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

const (
	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	testChainClientGetEVMAddressFromCertPath()
}

func testChainClientGetEVMAddressFromCertPath() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		panic(err)
	}

	certFilePath := fmt.Sprintf(examples.UserCrtPath, examples.OrgId1)
	addr, err := client.GetEVMAddressFromCertPath(certFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("client1 address: %s\n", addr)

	certFilePath = fmt.Sprintf(examples.UserCrtPath, examples.OrgId2)
	addr, err = client.GetEVMAddressFromCertPath(certFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("client2 address: %s\n", addr)
}
