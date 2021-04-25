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

func TestChainClient_SaveCert(t *testing.T) {

	type args struct {
		userCert    string
		enclaveCert string
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
				userCert:    "",
				enclaveCert: "",
			},
			want:    []byte{},
			wantErr: true,
		},
		{
			name: "test2",
			args: args{
				userCert:    "",
				enclaveCert: "",
			},
			want:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.SaveCert(tt.args.userCert, tt.args.enclaveCert)
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
				contractName: "",
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
		{
			name: "test1",
			args: args{
				contractName: "",
				privateKey:   "",
				userCert:     "",
				dirSign:      "",
			},
			want:    []byte{},
			wantErr: true,
		},
		{
			name: "test2",
			args: args{
				contractName: "",
				privateKey:   "",
				userCert:     "",
				dirSign:      "",
			},
			want:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.GetData(tt.args.contractName, tt.args.privateKey, tt.args.userCert, tt.args.dirSign)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
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
		computeResult string
		contractName  string
		gas           string
		reportSign    string
		userCert      string
		rwSet         *common.TxRWSet
		events        *common.StrSlice
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
				computeResult: "",
				contractName:  "",
				gas:           "",
				reportSign:    "",
				userCert:      "",
				rwSet:         nil,
				events:        nil,
			},
			want:    []byte{},
			wantErr: true,
		},

		{
			name: "test2",
			args: args{
				computeResult: "",
				contractName:  "",
				gas:           "",
				reportSign:    "",
				userCert:      "",
				rwSet:         nil,
				events:        nil,
			},
			want:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.SaveData(tt.args.computeResult, tt.args.contractName, tt.args.gas, tt.args.reportSign,
				tt.args.userCert, tt.args.rwSet, tt.args.events)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainClient_SaveDir(t *testing.T) {

	type args struct {
		userCert   string
		orderId    string
		dirHash    string
		dirSign    string
		privateDir *common.StrSlice
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
				userCert:   "",
				orderId:    "",
				dirHash:    "",
				dirSign:    "",
				privateDir: nil,
			},
			want:    []byte{},
			wantErr: true,
		},
		{
			name: "test2",
			args: args{
				userCert:   "",
				orderId:    "",
				dirHash:    "",
				dirSign:    "",
				privateDir: nil,
			},
			want:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cc, err := createClient()
			require.Nil(t, err)
			got, err := cc.SaveDir(tt.args.userCert, tt.args.orderId, tt.args.dirHash, tt.args.dirSign,
				tt.args.privateDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}
