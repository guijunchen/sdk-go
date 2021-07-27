/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package examples

import (
	"encoding/hex"
	"errors"
	"fmt"

	"chainmaker.org/chainmaker/common/evmutils"
	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
)

const (
	OrgId1 = "wx-org1.chainmaker.org"
	OrgId2 = "wx-org2.chainmaker.org"
	OrgId4 = "wx-org4.chainmaker.org"
	OrgId5 = "wx-org5.chainmaker.org"

	UserNameOrg1Admin1 = "org1admin1"
	UserNameOrg2Admin1 = "org2admin1"
	UserNameOrg3Admin1 = "org3admin1"
	UserNameOrg4Admin1 = "org4admin1"
	UserNameOrg5Admin1 = "org5admin1"

	certPathPrefix = "../../testdata"
	Version        = "1.0.0"
	UpgradeVersion = "2.0.0"
)

type user struct {
	SignKeyPath, SignCrtPath string
}

var (
	UserCrtPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.crt"
)

var users = map[string]*user{
	"org1admin1": {
		"../../testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key",
		"../../testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org2admin1": {
		"../../testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key",
		"../../testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org3admin1": {
		"../../testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key",
		"../../testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org4admin1": {
		"../../testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.key",
		"../../testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org5admin1": {
		"../../testdata/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.sign.key",
		"../../testdata/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.sign.crt",
	},
}

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

// CreateChainClientWithSDKConfDisableCertHash create a chain client with sdk config file path, disable cert hash.
func CreateChainClientWithSDKConfDisableCertHash(sdkConfPath string) (*sdk.ChainClient, error) {
	cc, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfPath),
	)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func CalcContractName(contractName string) string {
	return hex.EncodeToString(evmutils.Keccak256([]byte(contractName)))[24:]
}

func GetEndorsers(payload *common.Payload, usernames ...string) ([]*common.EndorsementEntry, error) {
	var endorsementEntrys []*common.EndorsementEntry

	for _, name := range usernames {
		u, ok := users[name]
		if !ok {
			return nil, errors.New("user not found")
		}

		entry, err := sdk.SignPayloadWithPath(u.SignKeyPath, u.SignCrtPath, payload)
		if err != nil {
			return nil, err
		}

		endorsementEntrys = append(endorsementEntrys, entry)
	}

	return endorsementEntrys, nil
}
