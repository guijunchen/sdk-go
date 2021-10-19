/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package examples

import (
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"

	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/common/v2/evmutils"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"
)

const (
	OrgId1 = "wx-org1.chainmaker.org"
	OrgId4 = "wx-org4.chainmaker.org"
	OrgId5 = "wx-org5.chainmaker.org"

	UserNameOrg1Client1 = "org1client1"
	UserNameOrg2Client1 = "org2client1"

	UserNameOrg1Admin1 = "org1admin1"
	UserNameOrg2Admin1 = "org2admin1"
	UserNameOrg3Admin1 = "org3admin1"
	UserNameOrg4Admin1 = "org4admin1"
	UserNameOrg5Admin1 = "org5admin1"

	Version        = "1.0.0"
	UpgradeVersion = "2.0.0"
)

var users = map[string]*User{
	"org1client1": {
		"../../testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key",
		"../../testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt",
		"../../testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key",
		"../../testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt",
	},
	"org2client1": {
		"../../testdata/crypto-config/wx-org2.chainmaker.org/user/client1/client1.tls.key",
		"../../testdata/crypto-config/wx-org2.chainmaker.org/user/client1/client1.tls.crt",
		"../../testdata/crypto-config/wx-org2.chainmaker.org/user/client1/client1.sign.key",
		"../../testdata/crypto-config/wx-org2.chainmaker.org/user/client1/client1.sign.crt",
	},
	"org1admin1": {
		"../../testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key",
		"../../testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt",
		"../../testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key",
		"../../testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org2admin1": {
		"../../testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.key",
		"../../testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.crt",
		"../../testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key",
		"../../testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org3admin1": {
		"../../testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.key",
		"../../testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.crt",
		"../../testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key",
		"../../testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org4admin1": {
		"../../testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.tls.key",
		"../../testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.tls.crt",
		"../../testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.key",
		"../../testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.crt",
	},
	"org5admin1": {
		"../../testdata/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.tls.key",
		"../../testdata/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.tls.crt",
		"../../testdata/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.sign.key",
		"../../testdata/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.sign.crt",
	},
}

type User struct {
	TlsKeyPath, TlsCrtPath   string
	SignKeyPath, SignCrtPath string
}

func GetUser(username string) (*User, error) {
	u, ok := users[username]
	if !ok {
		return nil, errors.New("user not found")
	}

	return u, nil
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

func CalcContractName(contractName string) string {
	return hex.EncodeToString(evmutils.Keccak256([]byte(contractName)))[24:]
}

func GetEndorsers(payload *common.Payload, usernames ...string) ([]*common.EndorsementEntry, error) {
	var endorsers []*common.EndorsementEntry

	for _, name := range usernames {
		u, ok := users[name]
		if !ok {
			return nil, errors.New("user not found")
		}

		var err error
		var entry *common.EndorsementEntry
		p11Handle := sdk.GetP11Handle()
		if p11Handle != nil {
			entry, err = sdkutils.MakeEndorserWithPathAndP11Handle(u.SignKeyPath, u.SignCrtPath, p11Handle, payload)
			if err != nil {
				return nil, err
			}
		} else {
			entry, err = sdkutils.MakeEndorserWithPath(u.SignKeyPath, u.SignCrtPath, payload)
			if err != nil {
				return nil, err
			}
		}

		endorsers = append(endorsers, entry)
	}

	return endorsers, nil
}

func MakeAddrAndSkiFromCrtFilePath(crtFilePath string) (string, string, string, error) {
	crtBytes, err := ioutil.ReadFile(crtFilePath)
	if err != nil {
		return "", "", "", err
	}

	blockCrt, _ := pem.Decode(crtBytes)
	crt, err := bcx509.ParseCertificate(blockCrt.Bytes)
	if err != nil {
		return "", "", "", err
	}

	ski := hex.EncodeToString(crt.SubjectKeyId)
	addrInt, err := evmutils.MakeAddressFromHex(ski)
	if err != nil {
		return "", "", "", err
	}

	fmt.Sprintf("0x%s", addrInt.AsStringKey())

	return addrInt.String(), fmt.Sprintf("0x%x", addrInt.AsStringKey()), ski, nil
}
