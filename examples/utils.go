/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package examples

import (
	"errors"
	"fmt"

	"chainmaker.org/chainmaker/common/log"
	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
)

const (
	chainId        = "chain1"
	OrgId1         = "wx-org1.chainmaker.org"
	OrgId2         = "wx-org2.chainmaker.org"
	OrgId3         = "wx-org3.chainmaker.org"
	OrgId4         = "wx-org4.chainmaker.org"
	OrgId5         = "wx-org5.chainmaker.org"
	orgId6         = "wx-org6.chainmaker.org"
	orgId7         = "wx-org7.chainmaker.org"
	certPathPrefix = "../../testdata"
	Version        = "1.0.0"
	UpgradeVersion = "2.0.0"
)

var (
	UserCrtPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.crt"
)

func CheckProposalRequestResp(resp *common.TxResponse, needContractResult bool) error {
	if resp.Code != common.TxStatusCode_SUCCESS {
		return errors.New(resp.Message)
	}

	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil && resp.ContractResult.Code != common.ContractResultCode_OK {
		return errors.New(resp.ContractResult.Message)
	}

	return nil
}

// CreateChainClientWithSDKConf create a chain client with sdk config file path
func CreateChainClientWithSDKConf(sdkConfPath string) (*sdk.ChainClient, error) {
	logger, _ := log.InitSugarLogger(&log.LogConfig{
		Module:       "[SDK]",
		LogPath:      "./sdk.log",
		LogLevel:     log.LEVEL_ERROR,
		MaxAge:       30,
		JsonFormat:   false,
		ShowLine:     true,
		LogInConsole: true,
	})

	var (
		cc  *sdk.ChainClient
		err error
	)

	cc, err = sdk.NewChainClient(
		sdk.WithConfPath(sdkConfPath),
		sdk.WithChainClientLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	// Enable certificate compression
	err = cc.EnableCertHash()
	if err != nil {
		return nil, err
	}
	return cc, nil
}
