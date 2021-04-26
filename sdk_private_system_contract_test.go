/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

const (
	computeContract = "private_computation"
)

func TestChainClient_SaveCert(t *testing.T) {

	type args struct {
		userCert       string
		enclaveCert    string
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
				userCert:       "",
				enclaveCert:    "",
				txId:           "",
				withSyncResult: true,
				timeout:        1000,
			},
			want:    []byte{},
			wantErr: true,
		},
		{
			name: "test2",
			args: args{
				userCert:       "",
				enclaveCert:    "",
				txId:           "",
				withSyncResult: true,
				timeout:        1000,
			},
			want:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.SaveCert(tt.args.userCert, tt.args.enclaveCert, tt.args.txId, tt.args.withSyncResult, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveCert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_GetContract(t *testing.T) {

	type args struct {
		userCert     string
		contractName string
		codeHash     string
		hashSign     string
	}

	tests := []struct {
		name    string
		args    args
		want    *common.PrivateGetContract
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				userCert:     "1",
				contractName: "2",
				codeHash:     "1",
				hashSign:     "1",
			},

			want:    nil,
			wantErr: true,
		},

		{
			name: "test2",
			args: args{
				userCert:     "",
				contractName: "my_asset",
				codeHash:     "",
				hashSign:     "",
			},

			want: &common.PrivateGetContract{
				ContractCode: nil,
				GasLimit:     0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.GetContract(tt.args.userCert, tt.args.contractName, tt.args.codeHash, tt.args.hashSign)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContract() error = %v, wantErr %v", err, tt.wantErr)
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
		privateKey   string
		userCert     string
		dirSign      string
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		//{
		//	name: "test1",
		//	args: args{
		//		contractName: computeContract,
		//		privateKey:   "key1",
		//		userCert:     "",
		//		dirSign:      "",
		//	},
		//	want:    []byte{},
		//	wantErr: true,
		//},
		//{
		//	name: "test2",
		//	args: args{
		//		contractName: computeContract,
		//		privateKey:   "key3",
		//		userCert:     "",
		//		dirSign:      "",
		//	},
		//	want:    []byte{},
		//	wantErr: true,
		//},
		{
			name: "test3",
			args: args{
				contractName: computeContract,
				privateKey:   "key5",
				userCert:     "",
				dirSign:      "",
			},
			want:    []byte("value_5"),
			wantErr: true,
		},
		//{
		//	name: "test4",
		//	args: args{
		//		contractName: computeContract,
		//		privateKey:   "dir_key1",
		//		userCert:     "",
		//		dirSign:      "",
		//	},
		//	want:    []byte{},
		//	wantErr: true,
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.GetData(tt.args.contractName, tt.args.privateKey, tt.args.userCert, tt.args.dirSign)
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

func TestChainClient_SaveData(t *testing.T) {

	type args struct {
		computeResult  string
		contractName   string
		gas            string
		reportSign     string
		userCert       string
		txId           string
		rwSet          *common.TxRWSet
		events         *common.StrSlice
		withSyncResult bool
		timeout        int64
	}

	rwSet := &common.TxRWSet{
		TxReads: []*common.TxRead{
			{Key: []byte("key1"), Value: []byte("value1"), ContractName: computeContract},
			{Key: []byte("key2"), Value: []byte("value2"), ContractName: computeContract},
			{Key: []byte("key3"), Value: []byte("value3"), ContractName: computeContract},
		},
		TxWrites: []*common.TxWrite{
			{Key: []byte("key3"), Value: []byte("value_3"), ContractName: computeContract},
			{Key: []byte("key4"), Value: []byte("value_4"), ContractName: computeContract},
			{Key: []byte("key5"), Value: []byte("value_5"), ContractName: computeContract},
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
				computeResult:  "true",
				contractName:   computeContract,
				gas:            "",
				reportSign:     "",
				userCert:       "",
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
			res, err := cc.SaveData(tt.args.computeResult, tt.args.contractName, tt.args.gas, tt.args.reportSign,
				tt.args.userCert, tt.args.txId, tt.args.rwSet, tt.args.events, tt.args.withSyncResult, tt.args.timeout)
			if string(res) != "OK" || err != nil || tt.wantErr != true {
				t.Errorf("SaveData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestChainClient_SaveDir(t *testing.T) {

	type args struct {
		userCert       string
		orderId        string
		dirHash        string
		dirSign        string
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
				userCert: "",
				orderId:  computeContract,
				dirHash:  "",
				dirSign:  "",
				privateDir: &common.StrSlice{
					StrArr: []string{"dir_key1", "dir_key2", "dir_key3"},
				},
				txId:           "",
				withSyncResult: true,
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
			got, err := cc.SaveDir(tt.args.userCert, tt.args.orderId, tt.args.dirHash, tt.args.dirSign, tt.args.txId,
				tt.args.privateDir, tt.args.withSyncResult, tt.args.timeout)
			if string(got) != "OK" || err != nil || tt.wantErr != true {
				t.Errorf("SaveDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestChainClient_SaveContract(t *testing.T) {

	type args struct {
		contractCode   []byte
		codeHash       string
		contractName   string
		txId           string
		version        string
		withSyncResult bool
		timeout        int64
	}
	tests := []struct {
		name    string
		args    args
		want    *common.ContractResult
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				contractCode:   nil,
				codeHash:       "",
				version:        version,
				withSyncResult: false,
				timeout:        1,
			},
			want: &common.ContractResult{
				Code:    0,
				Result:  nil,
				Message: "",
				GasUsed: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.SaveContract(tt.args.contractCode, tt.args.codeHash, tt.args.contractName, tt.args.txId,
				tt.args.version, tt.args.withSyncResult, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveContract() got = %v, want %v", got, tt.want)
			}
		})
	}
}
