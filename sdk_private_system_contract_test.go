/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"crypto/sha256"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

const (
	computeName = "compute_name"
	computeCode = "compute_code"
	ComputeRes  = "private_compute_result"
	enclaveId   = "enclave_id"
	enclaveCert = "enclave_certificate"
	quoteId     = "quote_id"
	quote       = "quote_content"
	orderId     = "order_id"
)

var priDir *common.StrSlice = &common.StrSlice{
	StrArr: []string{"dir_key1", "dir_key2", "dir_key3"},
}

func TestChainClient_SaveCert(t *testing.T) {

	type args struct {
		enclaveCert    string
		enclaveId      string
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
				enclaveCert:    enclaveCert,
				enclaveId:      enclaveId,
				txId:           "",
				withSyncResult: false,
				timeout:        1,
			},
			want:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.SaveCert(tt.args.enclaveCert, tt.args.enclaveId, tt.args.txId, tt.args.withSyncResult, tt.args.timeout)
			if err != nil {
				t.Errorf("SaveCert() error = %v, response %v", err, got)
				return
			}

		})
	}
}

func TestChainClient_SaveQuote(t *testing.T) {

	type args struct {
		enclaveId      string
		quoteId        string
		quote          string
		sign           string
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
				enclaveId:      enclaveId,
				quoteId:        quoteId,
				quote:          quote,
				sign:           "",
				txId:           "",
				withSyncResult: false,
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
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.SaveQuote(tt.args.enclaveId, tt.args.quoteId, tt.args.quote, tt.args.sign, tt.args.txId, tt.args.withSyncResult, tt.args.timeout)
			if got.ContractResult.Code != common.ContractResultCode_OK || err != nil {
				t.Errorf("SaveQuote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestChainClient_SaveContract(t *testing.T) {
	type args struct {
		codeBytes      []byte
		codeHash       string
		contractName   string
		txId           string
		version        string
		withSyncResult bool
		timeout        int64
	}

	codeHash := sha256.Sum256([]byte(computeCode))
	tests := []struct {
		name    string
		args    args
		want    *common.ContractResult
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				contractName:   computeName,
				codeBytes:      []byte(computeCode),
				codeHash:       string(codeHash[:]),
				version:        version,
				withSyncResult: false,
				timeout:        1,
			},
			want: &common.ContractResult{
				Code:    0,
				Result:  nil,
				Message: "OK",
				GasUsed: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.SaveContract(tt.args.codeBytes, tt.args.codeHash, tt.args.contractName, tt.args.version, tt.args.txId,
				tt.args.withSyncResult, tt.args.timeout)
			if err != nil {
				t.Errorf("SaveContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Message != "OK" {
				t.Errorf("SaveContract() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestChainClient_SaveData(t *testing.T) {

	type args struct {
		result         *common.ContractResult
		contractName   string
		txId           string
		rwSet          *common.TxRWSet
		events         *common.StrSlice
		withSyncResult bool
		timeout        int64
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

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				result: &common.ContractResult{
					Code:    0,
					Result:  nil,
					Message: "",
					GasUsed: 0,
				},
				contractName:   computeName,
				rwSet:          rwSet,
				events:         nil,
				withSyncResult: false,
				timeout:        1,
			},
			want:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			res, err := cc.SaveData(tt.args.contractName, tt.args.result, tt.args.txId, tt.args.rwSet, tt.args.events, tt.args.withSyncResult, tt.args.timeout)

			if res.ContractResult.Code != common.ContractResultCode_OK || err != nil || tt.wantErr != true { //todo check nil
				t.Errorf("SaveData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestChainClient_SaveDir(t *testing.T) {
	type args struct {
		orderId        string
		txId           string
		privateDir     *common.StrSlice
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
				orderId:        orderId,
				privateDir:     priDir,
				txId:           "",
				withSyncResult: false,
				timeout:        1,
			},
			want:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.SaveDir(tt.args.orderId, tt.args.txId, tt.args.privateDir, tt.args.withSyncResult, tt.args.timeout)
			if got.ContractResult.Code != common.ContractResultCode_OK || err != nil { //todo check nil
				t.Errorf("SaveDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestChainClient_GetContract(t *testing.T) {
	type args struct {
		contractName string
		codeHash     string
	}

	codeHash := sha256.Sum256([]byte(computeCode))

	tests := []struct {
		name    string
		args    args
		want    *common.PrivateGetContract
		wantErr interface{}
	}{
		{
			name: "test2",
			args: args{
				contractName: computeName,
				codeHash:     string(codeHash[:]),
			},
			want: &common.PrivateGetContract{
				ContractCode: []byte(computeCode),
				GasLimit:     10000000000,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.GetContract(tt.args.contractName, tt.args.codeHash) //todo check nil
			code := string(got.ContractCode)
			if err != nil {
				t.Errorf("GetContract() error = %v, wantErr %v, code %s", err, tt.wantErr, code)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContract() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetData(t *testing.T) {

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
			want:    []byte(enclaveCert),
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
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.GetData(tt.args.contractName, tt.args.key)
			if err != nil || tt.wantErr != true {
				t.Errorf("SaveData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetCert(t *testing.T) { //
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
				enclaveId: "111",
			},
			want:    []byte(""),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.GetCert(tt.args.enclaveId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetDir(t *testing.T) {

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
			want:    []byte(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.GetDir(tt.args.orderId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetQuote(t *testing.T) {

	type args struct {
		quoteId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{quoteId: "quoteId"},
			want:    []byte(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.GetQuote(tt.args.quoteId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQuote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQuote() got = %v, want %v", got, tt.want)
			}
		})
	}
}
