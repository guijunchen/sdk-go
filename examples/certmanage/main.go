/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
)

const (
	sdkConfigOrg1Admin1Path  = "./sdk_config_org1_admin1.yml"
	sdkConfigOrg1Client1Path = "./sdk_config_org1_client1.yml"

	sdkConfigOrg2Client1Path = "./sdk_config_org2_client1.yml"
)

func main() {
	testCertHash()
	testCertManage()
}

func testCertHash() {
	client, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		panic(err)
	}

	certHash := testCertAdd(client)
	time.Sleep(3 * time.Second)

	certInfos := testQueryCert(client, []string{certHash})
	if len(certInfos.CertInfos) != 1 {
		panic("require len(certInfos.CertInfos) == 1")
	}

	testDeleteCert(client, []string{certHash})
	time.Sleep(3 * time.Second)

	certInfos = testQueryCert(client, []string{certHash})
	if len(certInfos.CertInfos) != 1 {
		panic("require len(certInfos.CertInfos) == 1")
	}
	if certInfos.CertInfos[0].Cert != nil {
		panic("require certInfos.CertInfos[0].Cert == nil")
	}
}

func testCertManage() {
	var (
		err             error
		client1, admin1 *sdk.ChainClient

		// org2 client证书
		certs = []string{
			"-----BEGIN CERTIFICATE-----\nMIICiDCCAi6gAwIBAgIDCuSTMAoGCCqGSM49BAMCMIGKMQswCQYDVQQGEwJDTjEQ\nMA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt\nb3JnMi5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD\nExljYS53eC1vcmcyLmNoYWlubWFrZXIub3JnMB4XDTIwMTExNjA2NDYwNFoXDTI1\nMTExNTA2NDYwNFowgZAxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw\nDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcyLmNoYWlubWFrZXIub3Jn\nMQ8wDQYDVQQLEwZjbGllbnQxKzApBgNVBAMTImNsaWVudDEudGxzLnd4LW9yZzIu\nY2hhaW5tYWtlci5vcmcwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQjsmPDqPjx\nikMpRPkmWH8RFgUXwpzwaoMF9OQY6sAty2U8Q6TPlafMbm/xBls//UPZpi5uhwTv\neunkar0HqfvRo3sweTAOBgNVHQ8BAf8EBAMCAaYwDwYDVR0lBAgwBgYEVR0lADAp\nBgNVHQ4EIgQgjqe9Y2WHp+WC/GfKlvwummg3xvKPi9hbDja0QVFKa/EwKwYDVR0j\nBCQwIoAgmZcrtYWpTzN56LDZdqiHah3fG5w0kPaLoEBtyC8GfaEwCgYIKoZIzj0E\nAwIDSAAwRQIgbz8Du0bvtlWVJfBFzUamyfY2OodQDGBbKnr/eFXNeIECIQDnnJs5\nAX2NCT42Be3et+jhwxshehNsYm3WOOdTq/y+yg==\n-----END CERTIFICATE-----\n",
		}

		// org2 client证书的CRL
		certCrl = "-----BEGIN CRL-----\nMIIBXTCCAQMCAQEwCgYIKoZIzj0EAwIwgYoxCzAJBgNVBAYTAkNOMRAwDgYDVQQI\nEwdCZWlqaW5nMRAwDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcyLmNo\nYWlubWFrZXIub3JnMRIwEAYDVQQLEwlyb290LWNlcnQxIjAgBgNVBAMTGWNhLnd4\nLW9yZzIuY2hhaW5tYWtlci5vcmcXDTIxMDEyMTA2NDYwM1oXDTIxMDEyMTEwNDYw\nM1owFjAUAgMK5JMXDTI0MDMyMzE1MDMwNVqgLzAtMCsGA1UdIwQkMCKAIJmXK7WF\nqU8zeeiw2Xaoh2od3xucNJD2i6BAbcgvBn2hMAoGCCqGSM49BAMCA0gAMEUCIEgb\nQsHoMkKAKAurOUUfAJpb++DYyxXS3zhvSWPxIUPWAiEAyLSd4TgB9PbSgHyGzS5D\nU1knUTu/4HKTol6GuzmV0Kg=\n-----END CRL-----"
	)

	client1, err = examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		panic(err)
	}

	admin1, err = examples.CreateChainClientWithSDKConf(sdkConfigOrg1Admin1Path)
	if err != nil {
		panic(err)
	}

	fmt.Println("====================== 证书冻结 ======================")
	testCertManageFrozen(client1, admin1, certs)
	fmt.Println("====================== 已冻结，不可用 ======================")
	// authentication failed, checking certificate frozen list returns error: certificate is frozen
	testCertIsAvailable(false)

	fmt.Println("====================== 证书解冻 ======================")
	testCertManageUnfrozen(client1, admin1, certs)
	fmt.Println("====================== 已解冻，可用 ======================")
	testCertIsAvailable(true)

	fmt.Println("====================== 证书吊销 ======================")
	testCertManageRevoke(client1, admin1, certCrl)
	fmt.Println("====================== 已吊销，不可用 ======================")
	// authentication failed, checking CRL returns error: certificate is revoked
	testCertIsAvailable(false)
}

func testCertIsAvailable(isAvailable bool) {
	_, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg2Client1Path)
	if isAvailable {
		if err != nil {
			panic(err)
		}
	} else {
		if err == nil {
			panic("require err != nil")
		}
	}
}

func testCertManageFrozen(client1, admin1 *sdk.ChainClient, certs []string) {
	var (
		err                error
		payload            []byte
		signedPayloadBytes []byte
		resp               *common.TxResponse
		result             string
	)

	payload, err = client1.CreateCertManageFrozenPayload(certs)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes, err = admin1.SignCertManagePayload(payload)
	if err != nil {
		panic(err)
	}

	resp, err = client1.SendCertManageRequest(signedPayloadBytes, -1, true)
	if err != nil {
		panic(err)
	}

	result = string(resp.ContractResult.Result)

	fmt.Printf("resp: %+v, result:%s\n", resp, result)
}

func testCertManageUnfrozen(client1, admin1 *sdk.ChainClient, certs []string) {
	var (
		err                error
		payload            []byte
		signedPayloadBytes []byte
		resp               *common.TxResponse
		result             string
	)

	payload, err = client1.CreateCertManageUnfrozenPayload(certs)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes, err = admin1.SignCertManagePayload(payload)
	if err != nil {
		panic(err)
	}

	resp, err = client1.SendCertManageRequest(signedPayloadBytes, -1, true)
	if err != nil {
		panic(err)
	}

	result = string(resp.ContractResult.Result)

	fmt.Printf("unfrozen resp: %+v, result:%s\n", resp, result)
}

func testCertManageRevoke(client1, admin1 *sdk.ChainClient, certCrl string) {
	var (
		err                error
		payload            []byte
		signedPayloadBytes []byte
		resp               *common.TxResponse
		result             string
	)

	payload, err = client1.CreateCertManageRevocationPayload(certCrl)
	if err != nil {
		panic(err)
	}

	signedPayloadBytes, err = admin1.SignCertManagePayload(payload)
	if err != nil {
		panic(err)
	}

	resp, err = client1.SendCertManageRequest(signedPayloadBytes, -1, true)
	if err != nil {
		panic(err)
	}

	result = string(resp.ContractResult.Result)

	fmt.Printf("revoke resp: %+v, result:%s\n", resp, result)
}

func testCertAdd(client *sdk.ChainClient) string {
	resp, err := client.AddCert()
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(resp.ContractResult.Result)
}

func testQueryCert(client *sdk.ChainClient, certHashes []string) *common.CertInfos {
	certInfos, err := client.QueryCert(certHashes)
	if err != nil {
		panic(err)
	}
	return certInfos
}

func testDeleteCert(client *sdk.ChainClient, certHashes []string) {
	_, err := client.DeleteCert(certHashes)
	if err != nil {
		panic(err)
	}
}
