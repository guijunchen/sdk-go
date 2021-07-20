/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package examples

import (
	"chainmaker.org/chainmaker/common/evmutils"
	"encoding/hex"
	"errors"
	"fmt"

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

	if resp.ContractResult != nil && resp.ContractResult.Code != 0 {
		return errors.New(resp.ContractResult.Message)
	}

	return nil
}

// CreateChainClientWithSDKConf create a chain client with sdk config file path
func CreateChainClientWithSDKConf(sdkConfPath string) (*sdk.ChainClient, error) {
	cc, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfPath),
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

func GetEndorsers(payload *common.Payload, admins ...*sdk.ChainClient) ([]*common.EndorsementEntry, error) {
	var endorsers []*common.EndorsementEntry

	for _, admin := range admins {
		signedPayload, err := admin.SignContractManagePayload(payload)
		if err != nil {
			return nil, err
		}

		endorsers = append(endorsers, signedPayload)
	}

	return endorsers, nil
}

func CalcContractName(contractName string) string {
	return hex.EncodeToString(evmutils.Keccak256([]byte(contractName)))[24:]
}