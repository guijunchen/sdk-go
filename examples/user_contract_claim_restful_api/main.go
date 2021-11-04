/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/httputils"
	"chainmaker.org/chainmaker/common/v2/random/uuid"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
	"crypto/md5"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"strconv"
	"time"
)

const (
	claimContractName     = "claim_restful_001"
	claimVersion          = "2.0.0"
	claimByteCodePath     = "../../testdata/claim-wasm-demo/rust-fact-2.0.0.wasm"

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

func main() {
	testUserContractClaim()
}

func testUserContractClaim() {
	fmt.Println("====================== create client ======================")
	client, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfigOrg1Client1Path),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("====================== 创建合约 ======================")
	usernames := []string{examples.UserNameOrg1Admin1, examples.UserNameOrg2Admin1,
		examples.UserNameOrg3Admin1, examples.UserNameOrg4Admin1}
	testUserContractClaimCreate(client, usernames...)

	time.Sleep(3 * time.Second)

	fmt.Println("====================== 调用合约 ======================")
	fileHash, err := testUserContractClaimInvoke(client, "save", true)
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(3 * time.Second)

	fmt.Println("====================== 执行合约查询接口 ======================")
	kvs := []*common.KeyValuePair{
		{
			Key:   "file_hash",
			Value: []byte(fileHash),
		},
	}
	testUserContractClaimQuery(client, "find_by_file_hash", kvs)
}

func testUserContractClaimInvoke(client *sdk.ChainClient,
	method string, withSyncResult bool) (string, error) {

	curTime := strconv.FormatInt(time.Now().Unix(), 10)

	fileHash := uuid.GetUUID()
	kvs := []*common.KeyValuePair{
		{
			Key:   "time",
			Value: []byte(curTime),
		},
		{
			Key:   "file_hash",
			Value: []byte(fileHash),
		},
		{
			Key:   "file_name",
			Value: []byte(fmt.Sprintf("file_%s", curTime)),
		},
	}

	payload := client.CreatePayload("", common.TxType_INVOKE_CONTRACT, claimContractName, method, kvs, 0)

	req, err := client.GenerateTxRequest(payload, nil)
	if err != nil {
		log.Fatalln(err)
	}

	url := "http://localhost:12301/v1/sendrequest"

	resp, err := httputils.POST(nil, url, req)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("resp: %s\n\n", string(resp))

	return fileHash, nil
}

func testUserContractClaimQuery(client *sdk.ChainClient, method string, kvs []*common.KeyValuePair) {

	payload := client.CreatePayload("", common.TxType_QUERY_CONTRACT, claimContractName, method, kvs, 0)

	req, err := client.GenerateTxRequest(payload, nil)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("req: [len:%d] %+v\n\n", req.Size(), req)
	reqBytes, err := proto.Marshal(req)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("reqByte: [len:%d] %s\n\n", len(reqBytes), fmt.Sprintf("%x", md5.Sum(reqBytes)))

	url := "http://localhost:12301/v1/sendrequest"

	resp, err := httputils.POST(nil, url, req)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("resp: %s\n\n", string(resp))
}

func testUserContractClaimCreate(client *sdk.ChainClient, usernames ...string) {
	payload, err := client.CreateContractCreatePayload(claimContractName, claimVersion, claimByteCodePath,
		common.RuntimeType_WASMER, []*common.KeyValuePair{})
	if err != nil {
		log.Fatalln(err)
	}

	endorsers, err := examples.GetEndorsersWithAuthType(crypto.HashAlgoMap[client.GetHashType()],
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := client.GenerateTxRequest(payload, endorsers)
	if err != nil {
		log.Fatalln(err)
	}

	url := "http://localhost:12301/v1/sendrequest"

	resp, err := httputils.POST(nil, url, req)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("CREATE claim contract resp: %+v\n", string(resp))
}