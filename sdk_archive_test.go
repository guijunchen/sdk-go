/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chainmaker_sdk_go

import (
	"testing"

	"github.com/stretchr/testify/require"

	"chainmaker.org/chainmaker/pb-go/common"
)

func TestSendArchiveBlockRequest(t *testing.T) {
	tests := []struct {
		name         string
		height       uint64
		serverTxResp *common.TxResponse
		serverErr    error
	}{
		{
			"valid request",
			100,
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
			nil,
		},
		{
			"block already archived",
			100,
			&common.TxResponse{Code: common.TxStatusCode_ARCHIVED_BLOCK},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(tt.serverTxResp, tt.serverErr, WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			var (
				payload       *common.Payload
				signedPayload *common.Payload
				resp          *common.TxResponse
			)

			payload, err = cli.CreateArchiveBlockPayload(tt.height)
			require.Nil(t, err)

			signedPayload, err = cli.SignArchivePayload(payload)
			require.Nil(t, err)

			resp, err = cli.SendArchiveBlockRequest(signedPayload, -1)
			require.Nil(t, err)

			if resp.Code != tt.wantedResp.Code {
				t.Error("error: expected", tt.wantedResp, "received", resp)
			}
		})
	}
}

func TestSendRestoreBlockRequest(t *testing.T) {
	tests := []struct {
		name                       string
		fullblock                  []byte
		serverTxResp, wantedTxResp *common.TxResponse
		serverErr                  error
		wantedErr                  error
	}{
		{
			"valid request",
			[]byte("fullblock"),
			&common.TxResponse{Code: common.TxStatusCode_SUCCESS},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := newMockChainClient(WithConfPath(sdkConfigPathForUT))
			require.Nil(t, err)
			defer cli.Stop()

			var (
				payload       *common.Payload
				signedPayload *common.Payload
				resp          *common.TxResponse
			)

			payload, err = cli.CreateRestoreBlockPayload(tt.fullblock)
			require.Nil(t, err)

			signedPayload, err = cli.SignArchivePayload(payload)
			require.Nil(t, err)

			resp, err = cli.SendRestoreBlockRequest(signedPayload, -1)
			require.Nil(t, err)

			if resp.Code != tt.wantedResp.Code {
				t.Error("error: expected", tt.wantedResp, "received", resp)
			}
		})
	}
}
