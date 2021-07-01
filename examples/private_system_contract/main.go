/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"chainmaker.org/chainmaker/pb-go/common"
	sdk "chainmaker.org/chainmaker/sdk-go"
	"chainmaker.org/chainmaker/sdk-go/examples"
)

const (
	computeName  = "compute_name"
	computeCode  = "compute_code"
	computeCode2 = "compute_code2"
	ComputeRes   = "private_compute_result"
	//enclaveId    = "enclave_id"
	quoteId = "quote_id"
	quote   = "quote_content"
	orderId = "order_id"

	sdkConfigOrg1Client1Path = "../sdk_configs/sdk_config_org1_client1.yml"
)

var (
	proof     []byte
	enclaveId string
	caCert    []byte
	report    string
)

func main() {
	testChainClientSaveData()
	testChainClientSaveDir()
	testChainClientGetContract()
	testChainClientGetData()
	testChainClientGetDir()
	testChainClientSaveCACert()
	testChainClientGetCACert()
	testChainClientSaveEnclaveReport()
	testChainClientSaveRemoteAttestationProof()
	testChainClientGetEnclaveEncryptPubKey()
	testChainClientGetEnclaveVerificationPubKey()
	testChainClientGetEnclaveReport()
	testChainClientGetEnclaveChallenge()
	testChainClientGetEnclaveSignature()
	testChainClientGetEnclaveProof()
}

func readFileData(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return data
}

func initCaCert() {
	caCert = readFileData("../../testdata/remote_attestation/enclave_cacert.crt")
}

func initProof() {
	var err error
	proofHex := readFileData("../../testdata/remote_attestation/proof.hex")
	proof, err = hex.DecodeString(string(proofHex))
	if err != nil {
		panic(err)
	}
}

func initEnclaveId() {
	enclaveId = "global_enclave_id"
}

func initReport() {
	reportBytes := readFileData("../../testdata/remote_attestation/report.dat")
	report = hex.EncodeToString(reportBytes)
}

var priDir = &common.StrSlice{
	StrArr: []string{"dir_key1", "dir_key2", "dir_key3"},
}

func testChainClientSaveData() {
	codeHash := sha256.Sum256([]byte(computeCode))
	txid := sdk.GetRandTxId()
	result := &common.ContractResult{
		Code:    0,
		Result:  nil,
		Message: "",
		GasUsed: 0,
	}
	rwSet := &common.TxRWSet{
		TxReads: []*common.TxRead{
			{Key: []byte("key2"), Value: []byte("value2"), ContractName: computeName},
		},
		TxWrites: []*common.TxWrite{
			{Key: []byte("key1"), Value: []byte("value_1"), ContractName: computeName},
			{Key: []byte("key3"), Value: []byte("value_3"), ContractName: computeName},
			{Key: []byte("key4"), Value: []byte("value_4"), ContractName: computeName},
			{Key: []byte("key5"), Value: []byte("value_5"), ContractName: computeName},
		},
	}
	// todo add reportHash,sign
	var reportHash []byte
	var reportSign []byte

	cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		panic(err)
	}
	res, err := cc.SaveData(computeName, examples.Version, false, codeHash[:], reportHash, result, []byte(""), txid, rwSet,
		reportSign, nil, []byte(""), false, 1)
	if err != nil {
		panic(err)
	}
	if res.ContractResult.Code != common.ContractResultCode_OK {
		panic("res.ContractResult.Code != common.ContractResultCode_OK")
	}
}

func testChainClientSaveDir() {
	txid := sdk.GetRandTxId()

	cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
	if err != nil {
		panic(err)
	}
	got, err := cc.SaveDir(orderId, txid, priDir, false, 1)
	if err != nil {
		panic(err)
	}
	if got.ContractResult.Code != common.ContractResultCode_OK {
		panic("got.ContractResult.Code != common.ContractResultCode_OK")
	}
}

func testChainClientGetContract() {
	type args struct {
		contractName string
		codeHash     string
	}

	codeHash := sha256.Sum256([]byte(computeCode))
	codeHash2 := sha256.Sum256([]byte(computeCode2))

	tests := []struct {
		name    string
		args    args
		want    *common.PrivateGetContract
		wantErr interface{}
	}{
		{
			name: "test1",
			args: args{
				contractName: computeName,
				codeHash:     string(codeHash[:]),
			},
			want: &common.PrivateGetContract{
				ContractCode: []byte(computeCode),
				Version:      examples.Version,
				GasLimit:     10000000000,
			},
			wantErr: nil,
		},
		{
			name: "test2",
			args: args{
				contractName: computeName,
				codeHash:     string(codeHash2[:]),
			},
			want: &common.PrivateGetContract{
				ContractCode: []byte(computeCode2),
				Version:      examples.UpgradeVersion,
				GasLimit:     10000000000,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.GetContract(tt.args.contractName, tt.args.codeHash)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			panic("!reflect.DeepEqual(got, tt.want)")
		}
	}
}

func testChainClientGetData() {

	type args struct {
		contractName string
		key          string
	}

	dirByte, _ := priDir.Marshal()
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				contractName: computeName,
				key:          "key1",
			},
			want:    []byte("value_1"),
			wantErr: true,
		},
		{
			name: "test2",
			args: args{
				contractName: "",
				key:          orderId,
			},
			want:    dirByte,
			wantErr: true,
		},
		{
			name: "test3",
			args: args{
				contractName: "",
				key:          enclaveId,
			},
			want:    caCert,
			wantErr: true,
		},
		{
			name: "test4",
			args: args{
				contractName: "",
				key:          quoteId,
			},
			want:    []byte(quote),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.GetData(tt.args.contractName, tt.args.key)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			panic("!reflect.DeepEqual(got, tt.want)")
		}
	}
}

func testChainClientGetDir() {

	type args struct {
		orderId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				orderId: "orderId",
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.GetDir(tt.args.orderId)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			panic("!reflect.DeepEqual(got, tt.want)")
		}
	}
}

func testChainClientSaveCACert() {

	initCaCert()

	type args struct {
		caCert         string
		txId           string
		withSyncResult bool
		timeout        int64
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				caCert:         string(caCert),
				txId:           "",
				withSyncResult: true,
				timeout:        1,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.SaveEnclaveCACert(tt.args.caCert, tt.args.txId, tt.args.withSyncResult, tt.args.timeout)
		if err != nil {
			panic(err)
		}
		fmt.Printf("testChainClientSaveCACert got %+v\n", got)
	}
}

func testChainClientGetCACert() {

	initCaCert()

	type args struct {
		txId           string
		withSyncResult bool
		timeout        int64
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				txId:           "",
				withSyncResult: true,
				timeout:        1,
			},
			want:    caCert,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.GetEnclaveCACert()
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			panic("!reflect.DeepEqual(got, tt.want)")
		}
	}
}

func testChainClientSaveEnclaveReport() {

	initEnclaveId()
	initReport()

	type args struct {
		enclaveId      string
		report         string
		txId           string
		withSyncResult bool
		timeout        int64
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId:      enclaveId,
				report:         report,
				txId:           "",
				withSyncResult: true,
				timeout:        1,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.SaveEnclaveReport(tt.args.enclaveId, tt.args.report, tt.args.txId, tt.args.withSyncResult, tt.args.timeout)
		if err != nil {
			panic(err)
		}
		fmt.Printf("testChainClientSaveEnclaveReport got %+v\n", got)
	}
}

func testChainClientSaveRemoteAttestationProof() {

	initProof()

	type args struct {
		proof          string
		txId           string
		withSyncResult bool
		timeout        int64
	}

	tests := []struct {
		name    string
		args    args
		want    *common.TxResponse
		wantErr bool
	}{
		{
			name: "TEST1",
			args: args{
				proof:          string(proof),
				txId:           "",
				withSyncResult: true,
				timeout:        -1,
			},
			want: &common.TxResponse{
				Code:           0,
				Message:        "OK",
				ContractResult: nil,
			},
		},
	}

	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.SaveRemoteAttestationProof(tt.args.proof, tt.args.txId, tt.args.withSyncResult, tt.args.timeout)
		if err != nil {
			panic(err)
		}
		fmt.Printf("enclaveId = 0x%x \n", got.ContractResult.Result)
	}
}

func testChainClientGetEnclaveEncryptPubKey() {

	initEnclaveId()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.GetEnclaveEncryptPubKey(tt.args.enclaveId)
		if err != nil {
			panic(err)
		}
		fmt.Printf("encrypt pub key => %s \n", got)
	}
}

func testChainClientGetEnclaveVerificationPubKey() {

	initEnclaveId()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.GetEnclaveVerificationPubKey(tt.args.enclaveId)
		if err != nil {
			panic(err)
		}
		fmt.Printf("verification pub key => %s \n", got)
	}
}

func testChainClientGetEnclaveReport() {

	initEnclaveId()
	initReport()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    report,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.GetEnclaveReport(tt.args.enclaveId)
		if err != nil {
			panic(err)
		}
		fmt.Printf("testChainClientGetEnclaveReport got %+v\n", got)
	}
}

func testChainClientGetEnclaveChallenge() {

	initEnclaveId()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.GetEnclaveChallenge(tt.args.enclaveId)
		if err != nil {
			panic(err)
		}
		fmt.Printf("challenge => %s \n", got)
	}
}

func testChainClientGetEnclaveSignature() {

	initEnclaveId()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.GetEnclaveSignature(tt.args.enclaveId)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			panic("!reflect.DeepEqual(got, tt.want)")
		}
	}
}

func testChainClientGetEnclaveProof() {

	initEnclaveId()
	initProof()

	type args struct {
		enclaveId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				enclaveId: enclaveId,
			},
			want:    proof,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cc, err := examples.CreateChainClientWithSDKConf(sdkConfigOrg1Client1Path)
		if err != nil {
			panic(err)
		}
		got, err := cc.GetEnclaveSignature(tt.args.enclaveId)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			panic("!reflect.DeepEqual(got, tt.want)")
		}
	}
}
